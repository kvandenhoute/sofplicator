# SOFPLICATOR: The Sophisticated Replicator
Is Harbor the glue of your company?
Is harbor failing too much?
Is harbor not scaling nicely?
Enter: The Sofplicator

Sofplicator aimes to provide easy scaling from one registry to another.


# Requirements
The Sofplicator requires:

* A job replicator image
    * Can be set using JOB_REGISTRY, JOB_REPOSITORY and JOB_TAG env variables
* An Azure Vault that contains a username and password (to log in to the ACR)
    * Can be set using ACR_USERNAME_KEY and ACR_PASSWORD_KEY env variables
* Labels can be set on an ACR to include it in global replication (global replication = start a replication to all available targets instead of once specific)
    * Can be set using ACR_TARGET_LABEL_KEY, ACR_TARGET_LABEL_VALUE




## Build
`docker build -f ./build/package/sofplicator . -t sofplicator:0.5.0`

## RUN
```
docker run \
-p 8080:8080 \
sofplicator:0.1.2

```