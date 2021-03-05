# Maker - a 100DaysOfCode Project

### noun 

1. a person or thing that makes.
2. a manufacturer (used in combination)
3. a person who has the hobby of creating tangible physical products, especially do-it-yourself technology and engineering projects or handmade crafts

Maker is a CLI written to create various types of services in various cloud providers such as VM's, K8s clusters, storage buckets and SQL Databases. Its not meant to be a full replacement for each providers own CLI's or clients, but something handy for just creating a few key objects without needing 4+ CLI's or multiple Terraform providers. Handy for spinning up and down infra for labs and devlopment work kinda thing. 

This project is a personal [100 Days of Code](https://www.100daysofcode.com/) project in an effort to break the tutorial death cycle and start having some fun wtih Go! If you're curious to learn more about the various cloud provider's API's this is a good repo to poke through and see what it takes to do some general tasks. I learned a ton!


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