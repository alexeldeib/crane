package controller

import (
	"fmt"
	"log"

	"github.com/alexeldeib/build/crane/providers/azure/virtualnetwork"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Resource describe the generic shape of a cloud object.
type Resource struct {
	Type string
	Name string
	Args map[string]interface{}
}

// Execute ...
func Execute() {
	// resources := make(map[string]Resource)
	var resources []Resource
	mapstructure.Decode(viper.Get("resources"), &resources)

	fmt.Printf("%s --- %s\n", resources[0].Type, resources[0].Name)
	vm, err := virtualnetwork.GetSchema(resources[0].Args)
	if err != nil {
		log.Fatal(err)
	}
	virtualnetwork.Create(vm)
}
