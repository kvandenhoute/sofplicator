package util

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerregistry/mgmt/containerregistry"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/mgmt/keyvault"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/subscription/mgmt/subscription"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"

	"github.com/jongio/azidext/go/azidext"
	log "github.com/sirupsen/logrus"
)

type RegistryWithVault struct {
	Registry containerregistry.Registry
	Vault    keyvault.Vault
}

func ListSubscriptions() subscription.ListResultPage {
	// Passing nil configures the credential for Azure Public Cloud. To run in another cloud, such as Azure Government,
	// specify the Azure Resource Manager scope for that cloud in azidext.DefaultAzureCredentialOptions
	log.Debug("Using public cloud credential manager")
	authorizer, err := azidext.NewDefaultAzureCredentialAdapter(nil)
	if err != nil {
		log.Fatalf("Could not use credential manager for azure public cloud: %+v", err)
		panic("Could not use credential manager for azure public cloud")
	}
	log.Debug("Get subscriptions client")
	client := subscription.NewSubscriptionsClient()
	client.Authorizer = authorizer

	log.Debug("List all subscriptions")
	subscriptions, err := client.List(context.Background())
	if err != nil {
		log.Fatalf("Could not use the subscriptions client to get all subscriptions: %+v", err)
		panic("Could not use the subscriptions client to get all subscriptions")
	}

	log.Trace("%d subscriptions found %+v", len(subscriptions.Values()), subscriptions.Values())
	return subscriptions
}

func GetAllACRsWithLabel(subscriptions subscription.ListResultPage, requiredTagKey string, requiredTagValue string) []RegistryWithVault {
	var registriesWithVault []RegistryWithVault
	log.Debug("Using public cloud credential manager")
	authorizer, err := azidext.NewDefaultAzureCredentialAdapter(nil)
	if err != nil {
		log.Fatalf("Could not use credential manager for azure public cloud: %+v", err)
		panic("Could not use credential manager for azure public cloud")
	}

	for _, s := range subscriptions.Values() {
		client := containerregistry.NewRegistriesClient(*s.SubscriptionID)
		client.Authorizer = authorizer

		keyvaultClient := keyvault.NewVaultsClient(*s.SubscriptionID)
		keyvaultClient.Authorizer = authorizer

		allRegistries, err := client.List(context.Background())
		if err != nil {
			log.Fatalf("Could not list registries: %+v", err)
			panic("Could not list registries")
		}
		for _, registry := range allRegistries.Values() {
			for tagKey, tagValue := range registry.Tags {
				if tagKey == requiredTagKey && *tagValue == requiredTagValue {
					var top int32 = 1
					allVaults, err := keyvaultClient.ListBySubscription(context.Background(), &top)
					if err != nil {
						log.Fatalf("Could not list registries: %+v", err)
						panic("Could not list registries")
					}
					registriesWithVault = append(registriesWithVault, RegistryWithVault{Registry: registry, Vault: allVaults.Values()[0]})

				}
			}
		}
	}

	log.Trace("%d registries found %+v", len(subscriptions.Values()), subscriptions.Values())
	return registriesWithVault
}

func GetAzSecret(secretName string, vaultURI string) string {
	log.Debug("Getting azure credentials from vault")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
		panic("Could not get valid azure login credentials")
	}
	client, err := azsecrets.NewClient(vaultURI, cred, nil)
	if err != nil {
		log.Fatalf("Error establishing connection to Vault: %+v", err)
		panic("Could not create a client for the keyvault")
	}
	// Get a secret. An empty string version gets the latest version of the secret.
	version := ""
	resp, err := client.GetSecret(context.TODO(), secretName, version, nil)
	if err != nil {
		log.Fatalf("failed to get the secret: %v", err)
	}
	return *resp.Value
}
