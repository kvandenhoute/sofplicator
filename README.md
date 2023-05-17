# SOFPLICATOR: The Sophisticated Replicator
Is Harbor the glue of your company?
Is harbor failing too much?
Is harbor not scaling nicely?
Enter: The Sofplicator

Sofplicator aimes to provide easy replicating from one registry to another.


# Requirements
The Sofplicator requires:

* A job replicator image
    * Can be set using JOB_REGISTRY, JOB_REPOSITORY and JOB_TAG env variables
* An Azure Vault that contains a username and password (to log in to the ACR)
    * Can be set using ACR_USERNAME_KEY and ACR_PASSWORD_KEY env variables

# How to replicate

## Check liveness
**Request**
`GET <SOFPLICATOR_URL>/healthz`
```
**Response**
{
    "message": "I'm alive!"
}
```
## Basic replication from unauthenticated registry, to another unauthenticated registry
**Request**
`POST <SOFPLICATOR_URL>/startReplication`
```
{
    "identifier": "kvdh",
    "source": {
        "url": "docker.io",
        "unauthenticated": true
    },
    "target": {
        "url": "nexus",
        "unauthenticated": true
    }
}
```
*identifier* is a unique identifier which will be used in all of the resources created for this replication
**Response**
```
{
    "jobId": "69a6264f-7406-4ade-86a3-74af6028ad57"
}
```


## Registries
When doing a replication from a source registry to a target registry, there are 3 ways of providing credentials *(replace xxx with source or target)*
1. Calling the replication with `xxxRegistry.username` and `xxxRegistry.password`
2. Calling the replication with `xxxRegistry.UseCredentialsFromAzureVault`. In this case, the ACR_USERNAME_KEY and ACR_PASSWORD_KEY environment variables have to be set, these will be the secrets keys to be used for in the Azure vault. The Sofplicator will use the vault provided via the ACR_VAULT environment variable.
3. Calling the replication with `xxxRegistry.useExistingSecret`, `xxxRegistry.usernameKey`, `xxxRegistry.passwordKey`. In this case, credentials will be read from an existing secret in the cluster. The secret should exist in the same namespace as the sofplicator is running. 



## Build
`podman build -f ./build/package/sofplicator . -t sofplicator:0.5.0`

## RUN
```
podman run \
-p 8080:8080 \
sofplicator:0.5.0

```