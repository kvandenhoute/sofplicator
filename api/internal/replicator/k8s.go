package replicator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kvandenhoute/sofplicator/api/internal/util"
	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateSecret(name string, username string, password string, namespace string) string {
	log.Info("Creating AZ secret")
	client := util.GetKubeClient()
	secretName := util.GenerateName(name)
	secrets := client.CoreV1().Secrets(namespace)
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		StringData: map[string]string{
			"AZ_ACR_USERNAME": username,
			"AZ_ACR_PASSWORD": password,
		},
	}
	_, err := secrets.Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create K8S secret: %s", err)
	}
	log.Info("Created secret: ", secretName)
	return secretName
}

func CreateConfigmap(name string, images []Artifact, charts []Artifact, namespace string) string {
	log.Info("Creating Configmap")
	client := util.GetKubeClient()
	configMapName := util.GenerateName(name)
	configMaps := client.CoreV1().ConfigMaps(namespace)
	imagesJson, err := json.Marshal(images)
	if err != nil {
		log.Error(fmt.Errorf("Error unmarschalling images"))
	}
	chartsJson, err := json.Marshal(charts)
	if err != nil {
		log.Error(fmt.Errorf("Error unmarschalling charts"))
	}
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Data: map[string]string{
			"images.json": string(imagesJson),
			"charts.json": string(chartsJson),
		},
	}
	_, err = configMaps.Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		log.Error("Failed to create ConfigMap: ", err)
	}
	log.Info("ConfigMap created: ", configMapName)
	return configMapName
}

func CreateJob(name string, image string, namespace string, imagePullSecret string, mountPath string, configMapName string, secretName string, vaultURI string) string {
	client := util.GetKubeClient()
	jobs := client.BatchV1().Jobs(namespace)
	log.Info("Starting job")
	jobName := util.GenerateName(name)
	var backOffLimit int32 = 2
	log.Info(vaultURI)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					ImagePullSecrets: []v1.LocalObjectReference{
						{
							Name: imagePullSecret,
						},
					},
					Containers: []v1.Container{
						{
							Name:            "job-container",
							Image:           image,
							ImagePullPolicy: "Always",
							Env: []v1.EnvVar{
								{
									Name:  "IMAGES_JSON_IN_FLAVOR_REPORTS",
									Value: fmt.Sprintf("%s/images.json", mountPath),
								},
								{
									Name:  "CHARTS_JSON_IN_FLAVOR_REPORTS",
									Value: fmt.Sprintf("%s/charts.json", mountPath),
								},
								{
									Name:  "target_url",
									Value: vaultURI,
								},
								{
									Name: "TARGET_USERNAME",
									ValueFrom: &v1.EnvVarSource{
										SecretKeyRef: &v1.SecretKeySelector{
											Key: "AZ_ACR_USERNAME",
											LocalObjectReference: v1.LocalObjectReference{
												Name: secretName,
											},
										},
									},
								},
								{
									Name: "TARGET_PASSWORD",
									ValueFrom: &v1.EnvVarSource{
										SecretKeyRef: &v1.SecretKeySelector{
											Key: "AZ_ACR_PASSWORD",
											LocalObjectReference: v1.LocalObjectReference{
												Name: secretName,
											},
										},
									},
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "artifacts",
									MountPath: mountPath,
								},
							},
							// Command: strings.Split(*cmd, " "),
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "artifacts",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: configMapName,
									},
								},
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}
	_, err := jobs.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create K8S job: %s", err)
	}
	log.Infof("Job created successfully with name: %s", jobName)

	return jobName
}
