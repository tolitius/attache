package consul

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

type ConsulSpec struct {
	Address    string
	Datacenter string
	Token      string
}

func consulToMap(consulSpec ConsulSpec, rootPath string) map[string][]byte {

	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = consulSpec.Address
	consulConfig.Datacenter = consulSpec.Datacenter
	consulConfig.Token = consulSpec.Token

	consul, err := consulapi.NewClient(consulConfig)

	if err != nil {
		log.Fatalf("failed to connect to consul: %+v, due to %v", consulSpec, err)
	}

	kv := consul.KV()

	config := make(map[string][]byte)

	kvps, _, err := kv.List(rootPath, nil)
	if err != nil {
		log.Fatalf("failed to fetch k/v pairs from consul: %+v, root path: %s. due to %v", consulSpec, rootPath, err)
	}

	for _, kvp := range kvps {
		if val := kvp.Value; val != nil {
			config[kvp.Key] = val
		}
	}

	for k, v := range config {
		fmt.Printf("read consul map entry: {:%s, %s}\n", k, v)
	}

	return config
}
