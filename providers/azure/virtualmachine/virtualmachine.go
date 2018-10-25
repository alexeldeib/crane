package virtualmachine

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/mitchellh/mapstructure"
)

// GetSchema for Virtual Machine resource
func GetSchema(args map[string]interface{}) (*compute.VirtualMachine, error) {
	result := new(compute.VirtualMachine)
	err := mapstructure.Decode(args, &result)
	return result, err
}

// Create a new virtual machine
func Create(args *compute.VirtualMachine) {
	resourceGroupName, subscriptionID := "test-rg", "e02554d5-76c6-489e-8406-a86956e6669f"
	computeClient := compute.NewVirtualMachinesClient(subscriptionID)
	groupsClient := resources.NewGroupsClient(subscriptionID)
	authorizer, err := auth.NewAuthorizerFromFile("https://management.azure.com")
	if err != nil {
		log.Fatalf("Failed to get authorizer.\n")
	}

	computeClient.Authorizer = authorizer
	groupsClient.Authorizer = authorizer

	exists, err := groupsClient.CheckExistence(context.Background(), resourceGroupName)

	if err != nil {
		log.Fatalln(err)
	}

	var group resources.Group

	if exists.StatusCode == 404 {
		group, err = groupsClient.CreateOrUpdate(context.Background(), resourceGroupName, resources.Group{Location: args.Location})
		if err != nil {
			log.Fatalln(err)
		}
	} else if exists.StatusCode == 204 {
		group, err = groupsClient.Get(context.Background(), resourceGroupName)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln("Could not find or create a resource group for this resource.")
	}

	vm, err := computeClient.CreateOrUpdate(context.Background(),
		*group.Name,
		*(args.Name),
		*args)
	if err == nil {
		fmt.Println("Success!")
	} else {
		fmt.Println(err)
	}

	fmt.Println(vm)
}
