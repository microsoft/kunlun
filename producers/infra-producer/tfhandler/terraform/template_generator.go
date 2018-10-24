package terraform

import (
	artifacts "github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/common/storage"
	. "github.com/Microsoft/kunlun/producers/infra-producer/tfhandler/terraform/templates"
)

type TemplateGenerator struct{}

func NewTemplateGenerator() TemplateGenerator {
	return TemplateGenerator{}
}

func (t TemplateGenerator) GenerateTemplate(manifest artifacts.Manifest, state storage.State) (string, error) {
	template := ""

	tmpl, _ := NewProviderTemplate()
	template += tmpl
	tmpl, _ = NewResourceGroupTemplate()
	template += tmpl

	for _, nsg := range manifest.NetworkSecurityGroups {
		tmpl, err := NewNSGTemplate(nsg)
		if err != nil {
			return "", err
		}
		template += tmpl
	}

	for _, lb := range manifest.LoadBalancers {
		tmpl, err := NewLoadBalancerTemplate(lb)
		if err != nil {
			return "", err
		}
		template += tmpl
	}

	for _, vnet := range manifest.VNets {
		tmpl, err := NewVirtualNetworkTemplate(vnet)
		if err != nil {
			return "", err
		}
		template += tmpl
	}

	for _, vmg := range manifest.VMGroups {
		var tmpl string
		var err error
		if vmg.Type == "vm" {
			tmpl, err = NewVMTemplate(vmg)
		} else {
			tmpl, err = NewVMSSTemplate(vmg)
		}
		if err != nil {
			return "", err
		}
		template += tmpl
	}

	for _, mysqlDB := range manifest.MysqlDatabases {
		tmpl, err := NewMysqlTemplate(mysqlDB)
		if err != nil {
			return "", err
		}
		template += tmpl
	}
	return template, nil
}
