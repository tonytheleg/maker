# Maker

### noun 

1. a person or thing that makes.
2. a manufacturer (used in combination)
3. a person who has the hobby of creating tangible physical products, especially do-it-yourself technology and engineering projects or handmade crafts

Maker is (well, will be...this is a new project) a CLI written to create various types of services in various cloud providers such as VM's, perhaps K8s clusters, storage buckets, etc. Its not meant to now be a full replacement for each providers own CLI's or clients, but something handy for just creating a few key objects without needing 4+ CLI's or multiple Terraform providers. Handy for spinning up and down infra for labs and devlopment work kinda thing. 

Goals:
* Target many of the key cloud providers (Digital Ocean, Linode, AWS, Azure, GCP)
* Create common compute objects (VM's, K8s Clusters)
* Create Non-VM storage objects (Object storage or similar)
* Create DB Services (SQL, NoSQL specific to each provider)
* Learn more Go, and the [Cobra](https://cobra.dev/) CLI framework