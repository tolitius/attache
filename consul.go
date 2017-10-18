package attache

import (
	"fmt"
	"log"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type ConsulSpec struct {
	Address    string
	Datacenter string
	Token      string
}

func connectToConsul(consulSpec ConsulSpec) (*consulapi.Client, error) {

	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = consulSpec.Address
	consulConfig.Datacenter = consulSpec.Datacenter
	consulConfig.Token = consulSpec.Token

	conn, err := consulapi.NewClient(consulConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to consul: %+v, due to %v", consulSpec, err)
	}

	return conn, nil
}

func ConsulToMap(consulSpec ConsulSpec, rootPath string) (map[string]string, error) {

	consul, err := connectToConsul(consulSpec)
	if err != nil {
		return nil, err
	}

	kv := consul.KV()

	config := make(map[string]string)

	kvps, _, err := kv.List(rootPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch k/v pairs from consul: %+v, root path: %s. due to %v", consulSpec, rootPath, err)
	}

	for _, kvp := range kvps {
		if val := kvp.Value; val != nil {
			config[kvp.Key] = string(val[:])
		}
	}

	for k, v := range config {
		log.Printf("read consul map entry: {:%s, %s}\n", k, v)
	}

	return config, nil
}

func MapToConsul(consulSpec ConsulSpec, config map[string]string) (time.Duration, error) {

	consul, err := connectToConsul(consulSpec)
	if err != nil {
		return -1, err
	}

	kv := consul.KV()

	var duration int64

	for k, v := range config {
		took, err := kv.Put(&consulapi.KVPair{Key: k, Value: []byte(v)}, nil)
		if err != nil {
			return -1, fmt.Errorf("could not put a key, value: {%s, %s} to consul due to %v", k, v, err)
		}
		duration += int64(took.RequestTime)
	}

	return time.Duration(duration), nil
}
