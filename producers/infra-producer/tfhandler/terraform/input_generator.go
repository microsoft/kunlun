package terraform

import (
	"fmt"

	artifacts "github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/common/storage"
	. "github.com/Microsoft/kunlun/producers/infra-producer/tfhandler/terraform/templates"
)

type InputGenerator struct {
}

func NewInputGenerator() InputGenerator {
	return InputGenerator{}
}

func (i InputGenerator) GenerateInput(manifest artifacts.Manifest, state storage.State) (string, error) {
	input := ""
	ipt, _ := NewResourceGroupInput(manifest.ResourceGroupName, manifest.Location, manifest.EnvName)
	input += ipt

	for _, lb := range manifest.LoadBalancers {
		ipt, err := NewLoadBalancerInput(lb)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		input += ipt
	}

	for _, nsg := range manifest.NetworkSecurityGroups {
		ipt, err := NewNSGInput(nsg)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		input += ipt
	}

	for _, vnet := range manifest.VNets {
		ipt, err := NewVirtualNetworkInput(vnet)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		input += ipt
	}

	for _, vmg := range manifest.VMGroups {
		var ipt string
		var err error
		if vmg.Type == "vm" {
			ipt, err = NewVMInput(vmg)
		} else {
			ipt, err = NewVMSSInput(vmg)
		}
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		input += ipt
	}

	for _, mysqlDB := range manifest.MysqlDatabases {
		ipt, _ := NewMysqlInput(mysqlDB)
		input += ipt
	}
	return input, nil
}

func (i InputGenerator) Credentials(state storage.State) map[string]string {
	return map[string]string{
		"azure_environment":     state.Azure.Environment,
		"azure_subscription_id": state.Azure.SubscriptionID,
		"azure_tenant_id":       state.Azure.TenantID,
		"azure_client_id":       state.Azure.ClientID,
		"azure_client_secret":   state.Azure.ClientSecret,
	}
}
