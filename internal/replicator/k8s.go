package replicator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/kvandenhoute/sofplicator/internal/config"
	"github.com/kvandenhoute/sofplicator/internal/util"
	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateSecret(name string, uuid string, username string, password string, namespace string) (string, error) {
	log.Info("Creating AZ secret")
	client := util.KubeClient()
	secretName := util.GenerateName(name, uuid)
	secrets := client.CoreV1().Secrets(namespace)
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Labels: map[string]string{
				"sofplicator-uuid": uuid,
			},
		},
		StringData: map[string]string{
			"AZ_ACR_USERNAME": username,
			"AZ_ACR_PASSWORD": password,
		},
	}
	_, err := secrets.Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	log.Info("Created secret: ", secretName)
	return secretName, nil
}

func CreateConfigmap(name string, uuid string, images []Artifact, charts []Artifact, namespace string) (string, error) {
	log.Info("Creating Configmap")
	client := util.KubeClient()
	configMapName := util.GenerateName(name, uuid)
	configMaps := client.CoreV1().ConfigMaps(namespace)
	imagesJson, err := json.Marshal(images)
	if err != nil {
		log.Error(fmt.Errorf("error unmarshalling images"))
	}
	chartsJson, err := json.Marshal(charts)
	if err != nil {
		log.Error(fmt.Errorf("error unmarshalling charts"))
	}
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
			Labels: map[string]string{
				"sofplicator-uuid": uuid,
			},
		},
		Data: map[string]string{
			"images.json": string(imagesJson),
			"charts.json": string(chartsJson),
		},
	}
	_, err = configMaps.Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	log.Info("ConfigMap created: ", configMapName)
	return configMapName, nil
}

func CreateJob(name string, uuid string, image string, namespace string, imagePullSecret string, mountPath string, configMapName string, secretName string, vaultURI string) (string, error) {
	client := util.KubeClient()
	jobs := client.BatchV1().Jobs(namespace)
	log.Info("Starting job")
	jobName := util.GenerateName(name, uuid)
	var backOffLimit int32 = 2
	log.Info(vaultURI)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
			Labels: map[string]string{
				"sofplicator-uuid": uuid,
			},
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
		return "", nil
	}
	log.Infof("Job created successfully with name: %s", jobName)

	return jobName, nil
}

func CleanupResources(uuid string, namespace string) error {
	err := deleteJob(uuid, namespace)
	if err != nil {
		return err
	}
	err = deleteConfigmap(uuid, namespace)
	if err != nil {
		return err
	}
	err = deleteSecret(uuid, namespace)
	if err != nil {
		return err
	}
	return nil
}

func deleteJob(uuid string, namespace string) error {
	client := util.KubeClient()
	job, err := GetJobOnLabel(uuid, namespace)
	if err != nil {
		return err
	}
	jobs := client.BatchV1().Jobs(namespace)
	log.Debug("Deleting job: %s", job.Name)
	err = jobs.Delete(context.TODO(), job.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func deleteConfigmap(uuid string, namespace string) error {
	client := util.KubeClient()
	configmap, err := getConfigmapOnLabel(uuid, namespace)
	if err != nil {
		return err
	}
	configmaps := client.CoreV1().ConfigMaps(namespace)
	log.Debug("Deleting configmap: %s", configmap.Name)
	err = configmaps.Delete(context.TODO(), configmap.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func deleteSecret(uuid string, namespace string) error {
	client := util.KubeClient()
	secret, err := getSecretOnLabel(uuid, namespace)
	if err != nil {
		return err
	}
	secrets := client.CoreV1().Secrets(namespace)
	log.Debug("Deleting secrets: %s", secret.Name)
	err = secrets.Delete(context.TODO(), secret.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func GetJobOnLabel(uuid string, namespace string) (batchv1.Job, error) {
	client := util.KubeClient()
	jobs := client.BatchV1().Jobs(namespace)
	jobList, err := jobs.List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("sofplicator-uuid=%s", uuid),
	})
	if err != nil {
		return batchv1.Job{}, fmt.Errorf("error listing jobs: %s", err)
	}
	log.Debugf("Get jobs: ", jobList)
	if len(jobList.Items) == 0 {
		return batchv1.Job{}, errors.New("no jobs found with unique label")
	}
	if len(jobList.Items) > 1 {
		return batchv1.Job{}, errors.New("more than one job found with unique label")
	}
	return jobList.Items[0], nil
}

func getSecretOnLabel(uuid string, namespace string) (v1.Secret, error) {
	client := util.KubeClient()
	secrets := client.CoreV1().Secrets(namespace)
	secretList, err := secrets.List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("sofplicator-uuid=%s", uuid),
	})
	if err != nil {
		return v1.Secret{}, fmt.Errorf("error listing secrets: %s", err)
	}
	log.Debugf("Get secrets: ", secretList)
	if len(secretList.Items) == 0 {
		return v1.Secret{}, errors.New("no secrets found with unique label")
	}
	if len(secretList.Items) > 1 {
		return v1.Secret{}, errors.New("more than one secret found with unique label")
	}
	return secretList.Items[0], nil
}

func getCurrentNamespace() string {
	if config.Get().TargetNamespace != "" && len(config.Get().TargetNamespace) == 0 {
		log.Trace("Use environment variable namespce")
		return config.Get().TargetNamespace
	}
	namespace, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	log.Trace("Tried reading namespace")
	log.Trace(string(namespace))
	if err != nil {
		log.Trace("Use default namespace (dev-tooling)")
		return "dev-tooling"
	}
	log.Trace("Use current namespace (" + string(namespace) + ")")
	return string(namespace)
}

func getConfigmapOnLabel(uuid string, namespace string) (v1.ConfigMap, error) {
	client := util.KubeClient()
	configmaps := client.CoreV1().ConfigMaps(namespace)
	configmapList, err := configmaps.List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("sofplicator-uuid=%s", uuid),
	})
	if err != nil {
		return v1.ConfigMap{}, fmt.Errorf("error listing configmaps: %s", err)
	}
	log.Debugf("Get configmaps: ", configmapList)
	if len(configmapList.Items) == 0 {
		return v1.ConfigMap{}, errors.New("no configmap found with unique label")
	}
	if len(configmapList.Items) > 1 {
		return v1.ConfigMap{}, errors.New("more than one secret found with unique label")
	}
	return configmapList.Items[0], nil
}
