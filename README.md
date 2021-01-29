# Maker

### noun 

1. a person or thing that makes.
2. a manufacturer (used in combination)
3. a person who has the hobby of creating tangible physical products, especially do-it-yourself technology and engineering projects or handmade crafts

Maker is (well, will be...this is a new project) a CLI written to create various types of services in various cloud providers such as VM's, perhaps K8s clusters, storage buckets, etc. Its not meant to be a full replacement for each providers own CLI's or clients, but something handy for just creating a few key objects without needing 4+ CLI's or multiple Terraform providers. Handy for spinning up and down infra for labs and devlopment work kinda thing. 

Goals:
* Target many of the key cloud providers (Digital Ocean, Linode, AWS, Azure, GCP)
* Create common compute objects (VM's, K8s Clusters)
* Create Non-VM storage objects (Object storage or similar)
* Create DB Services (SQL, NoSQL specific to each provider)
* Learn more Go, and the [Cobra](https://cobra.dev/) CLI framework


### Basic Usage

Create the required configuration/credential files
```shell
maker auth -p PROVIDER
# This will prompt for specific data or paths to key files depending on the provider
# A file is then created under $HOME/.maker

/home/USERNAME/.maker/
├── aws_credentials
├── do_config
└── gcp_config
```

Create a VM
```shell
maker create vm -p do -n test-vm -s s-1vcpu-1gb -i ubuntu-16-04-x64
# or perhaps
maker create vm -p gcp -s e2-micro -i ubuntu-os-cloud/ubuntu-1604-xenial-v20210119 -n test-gce
```

Create a S3 Bucket
```shell
maker create bucket -p aws -n my-super-special-bucket
```

Get the status of a VM
```shell
maker status vm -p aws -n ec2-test-instance
```

Delete a DigitalOcean Space
```shell
maker delete bucket -p do -n super-special-do-space
```