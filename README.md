# Kunlun

Kunlun is to tool for deploying and managing common OSS based workloads on Azure.
Using Kunlun allows users with no familiarity with Azure to deploy Java and LAMP
applications in an optimized way.

[![Build Status](https://xplaceholderci.gugagaga.fun/buildStatus/icon?job=kunlun/master)](https://xplaceholderci.gugagaga.fun/job/kunlun/job/master/)

[![GoDoc](https://godoc.org/github.com/Microsoft/kunlun?status.svg)](https://godoc.org/github.com/Microsoft/kunlun)

[![Go Report Card](https://goreportcard.com/badge/Microsoft/kunlun)](https://goreportcard.com/report/Microsoft/kunlun)

## Building Kunlun from Source

```
go get github.com/Microsoft/kunlun/cmd/kl
```

Now you will have a `kl` command. To validate the install run:

```
kl help
```

If you get a `No command 'kl' found` error then you like neglected to 
add your $GOPATH/bin to the path when installing Go.

```
export PATH=$PATH:$GOTPATH/bin
```

If you hit the error `undefined: strings.Builder`, please upgrade your 
Go to version `1.10` or later.


## Using Kunlun

### Prepare the Environment

* Install [Go](https://golang.org/doc/install) version `1.10` or later

* Install [Terraform](https://www.terraform.io/intro/getting-started/install.html)

  * Note that you need to ensure that the terraform binary is on your path

* Install [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)

* Install [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest)

#### Login to Azure

```
az login
```

If you have more than one subscription you should check you have the right one activated with 
`az account show`. If necessary change it with `az account set --subscription ...`, you can 
view your avilable subscriptions with `az account list`.
    
Once you are sure you are using the correct subscription place it's ID in an environment 
variable, while we are at it we'll grab the Tenant ID too:

```
export KL_AZURE_SUBSCRIPTION_ID=$(az account show --output tsv --query id)
export KL_AZURE_TENANT_ID="$(az account show --output tsv --query tenantId)"
```

#### Service Principal for Kunlun

If you don't already have a service principal for Kunlun we need to create one now. If you
have already created one you simply need to grab its client id, see the last command in this
section.

To create/use a service principal for Kunlun to use to manage your resources we will first
capture some important values in environment variables. First we need a name for the service
principle. The below command generates a name that includes a UUID, you may choose to provide
a more memorable name:

```
export KL_AZURE_APP_NAME=kunlun-$(uuidgen)
```

For convenience for this tutorial ONLY we'll store the client secret in an environment variable.
Obviously you don't want to do this in the real world. 

```
export KL_AZURE_CLIENT_SECRET=password
```

Now we are ready to create the service principal:

```
az ad sp create-for-rbac --name $KL_AZURE_APP_NAME --password $KL_AZURE_CLIENT_SECRET
```

For convenience we will store the Client ID in an environment variable:

```
export KL_AZURE_CLIENT_ID="$(az ad sp show --id http://$KL_AZURE_APP_NAME --output tsv --query appId)"
```

#### A Few More Environment Variables

It is also useful to set a few other convenience environment variables:

```
export KL_IAAS=azure
export KL_AZURE_ENVIRONMENT=public
export KL_AZURE_REGION=southcentralus
```

#### Check Environment Setup

At this point you should have a set of environment variables that will be used to make Kunlun use easier.
To view the current setup use:

```
env | grep '^KL_'
```

This will give you an output something like:

```
KL_IAAS=azure
KL_AZURE_CLIENT_ID=c53dc238-****-****-****-6217f401a917
KL_AZURE_REGION=southcentralus
KL_AZURE_TENANT_ID=49e892d5-****-****-****-98be9fe068e2
KL_AZURE_CLIENT_SECRET=password
KL_AZURE_SUBSCRIPTION_ID=325e7c34-****-****-****-1df746c67705
KL_AZURE_ENVIRONMENT=public
KL_AZURE_APP_NAME=kunlun
```

### Analyze the Application you wish to deploy

Change into your project working directory. For our demo we will create a
new project directory and fetch the kunlun tool

```
mkdir kunlun-test
cd kunlun-test
go get github.com/kunlun/kun-lun
cd go/src/github.com/kun-lun/kunlun/cmd/kl
```

And now we will analyze the application under consideration in order to select the correct
infrastructure.

```
kl analyze
```

This will ask a number of questions, and for the most part you can safely use
the defaults. You will, however, need to provide a resource group name.
This is the name of the Azure Resource Group into which Kunlun will place
all created resources (e.g. virtual machines, networks, storage). The 
resource group name must be unique to your subscription.

When asked for your application code you should provide a public Git 
repository. For your convenience we provide a couple of simple sample
applications:

PHP  : https://github.com/kun-lun/2048php.git
Java : https://github.com/kun-lun/2048java.git

You will get one folder called `artifacts` in your working dir. With 
a `main.yml` file and one `patches` folder. The format of these outputs
is unique to Kunlun and, for the most part you will not work with them.
However, they are provided so that advanced users have a significant
amount of control over their deployed resources. See the 'Advanced
Users' section for more information.

### Plan the infrastructure

Now we need to convert this Kunlun spec into something that can be
used to deploy the infrastructure required. The `plan_infra` command
will do this:

```
kl plan_infra
```

The outputs of this command are Terraform templates, which can be 
found in the `infra` folder. If the generated Terraform output is not
sufficient for your needs you can customize the plan. Please refer to the
'Advanced Users' section below.
 
### Infrastrcuture Configuration

Now it is time to deploy the infrastructure required for your application. 
This can be done with:

```
kl apply_infra
```

This command creates an `outputs.yml` in the `artifacts/patches` folder. This contains contents
such as:
 
```
- type: replace
  path: /vm_groups/name=jumpbox/networks/0/outputs?
  value:
    - public_ip: 40.87.54.187
- type: replace
  path: /vm_groups/name=web-servers/networks/0/outputs?
  value:
    - ip: 10.0.0.4
```

[FIXME: what does this next sentence mean?]
This file will be applied to the original artifact, 
and then our deployment component would digest and produce the deployment script, now, in ansible.
 
### Plan Software Deployment
 
We can now plan our software deployment to this infrastructure:

```
kl plan_deployment
```

After running this command you will see a folder called `deployments`. This folder contains 
the ansible deployment scripts for your application. If this doesn't look like it covers your
requirements you can customize it. See the 'Advanced Users' section below.

### Deploy Software

Run `kl apply_deployment` to do the real deployment.

## Advanced Users

Kunlun is designed to be very flexible. As such it provides a number
of ways advanced users can tune it for their specific purposes.

### Analyze

The `analyze` command outputs a `main.yml` file that contains a description 
of the system that will be deployed. If you think the infrastrcuture selected
for your application does not meet your requirements, you can create a 
patch file under the `patches` folder to change the final configuration:

```
- type: replace
  path: /vm_groups/name=jumpbox/sku
  value: Standard_DS2_v2
- type: replace
  path: /vm_groups/name=web-servers/sku
  value: Standard_DS2_v2
```

### Plan

The `plan_infra` command creates Terraform templates to deploy the actual 
resources. If you want to setup some additional resources, you can add
additional Terraform files in the `infra` folder.

The `plan_deployment` command attempts to build Ansible deployments for your 
software code. If you think our built-in artifacts do not meet your requirements, 
you can create a patch file to add more roles into the artifact and run 
`kl plan_deployment` again. For example, you might want to add a firewall 
component:

```
- type: replace
  path: /vm_groups/name=web-servers/roles/-
  value:
    name: geerlingguy.firewall
```

## SSH into your VMs

Run `kl ssh -group <YOUR VM GROUP NAME> -index <YOUR NODE INDEX>` to ssh into your vm instance in one group.
