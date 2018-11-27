# kunlun

[![Build Status](https://xplaceholderci.gugagaga.fun/buildStatus/icon?job=kunlun/master)](https://xplaceholderci.gugagaga.fun/job/kunlun/job/master/)

[![GoDoc](https://godoc.org/github.com/Microsoft/kunlun?status.svg)](https://godoc.org/github.com/Microsoft/kunlun)

[![Go Report Card](https://goreportcard.com/badge/Microsoft/kunlun)](https://goreportcard.com/report/Microsoft/kunlun)

## Prepare the Environment

* Install [Go](https://golang.org/doc/install) version `1.10` or later

* Install [Terraform](https://www.terraform.io/intro/getting-started/install.html)

  * Note that you need to ensure that the terraform binary is on your path

* Install [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)

* Login to Azure

    ```
    az login
    export KL_AZURE_SUBSCRIPTION_ID=<YOUR SUBSCRIPTION ID>
    az account set --subscription $KL_AZURE_SUBSCRIPTION_ID
    ```

* Create a Service Principle for your application

    ```
    export KL_AZURE_CLIENT_SECRET=password
    export KL_AZURE_CLIENT_ID="$(az ad sp create-for-rbac --name kunlun --password $KL_AZURE_CLIENT_SECRET --output tsv --query appId)"
    export KL_AZURE_TENANT_ID=$(az ad sp show --id $KL_AZURE_CLIENT_ID --output tsv --query appOwnerTenantId)
    ```

* Set some environment variables

    ```
    export KL_IAAS=azure
    export KL_AZURE_ENVIRONMENT=public
    export KL_AZURE_REGION=southcentralus
    ```

## Building from Source

```
go get github.com/Microsoft/kunlun/cmd/kl
```

Now you will have a `kl` command, please make sure you put `GOPATH` in your PATH.

If you hit the error `undefined: strings.Builder`, please upgrade your Go to version `1.10` or later.

## Analyze the Application you wish to deploy

Run `kl analyze`.

Please type in the git repository address for your application, e.g. https://github.com/kun-lun/2048php.git, as the project path.

You will get one folder called `artifacts` in your working dir. With a `main.yml` file and one `patches` folder.

The `main.yml` contains a description of the system that will be deployed.
 
If you think the file does not meet your requirements, you can create one patch file under the `patches` folder like this:

```
- type: replace
  path: /vm_groups/name=jumpbox/sku
  value: Standard_DS2_v2
- type: replace
  path: /vm_groups/name=web-servers/sku
  value: Standard_DS2_v2
```

## Plan the infrastructure

Run `kl plan_infra` to generate the terraform scripts that will deploy your infrastrcuture.

You will got an `infra` folder in your working dir.

If you want to setup some additional resources, you can also put the resources terraform files under the same folder.
 
## Infrastrcuture Configuration

Run `kl apply_infra`.

`outputs.yml` will be generated under the `patches` folder, with the content like this:
 
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
 
## Plan Deployment
 
Run `kl plan_deployment`

You will get a folder called `deployments` which contains the deployment scripts.

And if you think our built-in artifacts does not meet your requirements, 
you can create one patch file to add more roles into the artifact and run 
`kl plan_deployment` again. For example, you might want to add a firewall component:

```
- type: replace
  path: /vm_groups/name=web-servers/roles/-
  value:
    name: geerlingguy.firewall
```

## Deploy

Run `kl apply_deployment` to do the real deployment.
