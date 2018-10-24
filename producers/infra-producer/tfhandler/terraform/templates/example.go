package templates

import (
	"fmt"
	"io/ioutil"

	artifacts "github.com/Microsoft/kunlun/artifacts"
)

func MainExample() {
	m, err := artifacts.NewManifestFromYAML(example_artifact)
	if err != nil {
		fmt.Println(err)
	}

	template := ""

	tmpl, _ := NewResourceGroupTemplate()
	template += tmpl

	for _, nsg := range m.NetworkSecurityGroups {
		tmpl, err := NewNSGTemplate(nsg)
		if err != nil {
			fmt.Println(err)
			return
		}
		template += tmpl
	}

	for _, lb := range m.LoadBalancers {
		tmpl, err := NewLoadBalancerTemplate(lb)
		if err != nil {
			fmt.Println(err)
			return
		}
		template += tmpl
	}

	for _, vnet := range m.VNets {
		tmpl, err := NewVirtualNetworkTemplate(vnet)
		if err != nil {
			fmt.Println(err)
			return
		}
		template += tmpl
	}

	for _, vmg := range m.VMGroups {
		var tmpl string
		if vmg.Type == "vm" {
			tmpl, _ = NewVMTemplate(vmg)
		} else {
			tmpl, _ = NewVMSSTemplate(vmg)
		}
		template += tmpl
	}

	for _, mysqlDB := range m.MysqlDatabases {
		tmpl, _ := NewMysqlTemplate(mysqlDB)
		template += tmpl
	}

	input := ""
	ipt, _ := NewResourceGroupInput(m.ResourceGroupName, m.Location, m.EnvName)
	input += ipt

	for _, lb := range m.LoadBalancers {
		ipt, err := NewLoadBalancerInput(lb)
		if err != nil {
			fmt.Println(err)
			return
		}
		input += ipt
	}

	for _, nsg := range m.NetworkSecurityGroups {
		ipt, err := NewNSGInput(nsg)
		if err != nil {
			fmt.Println(err)
			return
		}
		input += ipt
	}

	for _, vnet := range m.VNets {
		ipt, err := NewVirtualNetworkInput(vnet)
		if err != nil {
			fmt.Println(err)
			return
		}
		input += ipt
	}

	for _, vmg := range m.VMGroups {
		fmt.Println(vmg.NetworkInfos[0].SubnetName)
		var ipt string
		var err error
		if vmg.Type == "vm" {
			ipt, err = NewVMInput(vmg)
		} else {
			ipt, err = NewVMSSInput(vmg)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		input += ipt
	}

	for _, mysqlDB := range m.MysqlDatabases {
		ipt, _ := NewMysqlInput(mysqlDB)
		input += ipt
	}

	ioutil.WriteFile("template.tf", []byte(template), 0755)
	ioutil.WriteFile("input.tfvars", []byte(input), 0755)
}
