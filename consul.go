/*
Package attache is a younger brother of https://github.com/tolitius/envoy
that makes a bridge between Consul and application data structures a little more beautiful.
*/
package attache

import (
	"fmt"
	"log"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

// ConsulToMap takes a consul config and a path offset
// Connects to consul "key/value".
// Reads all (i.e. "recurse") {k, v} pairs under the path offset
// into a map[string]string preserving path hierarchy in map keys: i.e. {"universe/answers/main": "42"}
func ConsulToMap(consulSpec *consulapi.Config, offset string, keysWithOffset ...bool) (map[string]string, error) {

	consul, err := consulapi.NewClient(consulSpec)
	if err != nil {
		return nil, err
	}

	kv := consul.KV()

	config := make(map[string]string)

	kvps, _, err := kv.List(offset, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch k/v pairs from consul: %+v, path offset: %s. due to %v", consulSpec, offset, err)
	}

	withOffset := true
	if len(keysWithOffset) > 0 {
		withOffset = keysWithOffset[0]
	}

	for _, kvp := range kvps {
		if val := kvp.Value; val != nil {
			k := kvp.Key
			if !withOffset {
				k = strings.Split(kvp.Key, offset)[1]
			}
			config[k] = string(val[:])
		}
	}

	log.Printf("read consul map at offset: /%s\n", offset)

	excludeList := []string{"password", "secret", "auth"}
	for k, v := range config {
		_, found := func(slice []string, key string) (int, bool) {
			for i, item := range slice {
				search := strings.Contains(key, item)
				if search != false {
					return i, true
				}
			}
			return -1, false
		}(excludeList, k)

		if !found {
			fmt.Printf("read consul map entry: {:%s, %s }\n", k, v)
		} else {
			fmt.Printf("read consul map entry: {:%s, %s }\n", k, "*******")
		}
	}

	return config, nil
}

// MapToConsul takes a consul config and a map[string]string
// Connects to consul "key/value".
// Walks over a given map and "PUT"s its etries to consul
// respecting path hierarchy encoded in keys: i.e. {"universe/answer/main": 42}.
// Returns a total time.Duration of all the "PUT" operations
func MapToConsul(consulSpec *consulapi.Config, config map[string]string) (time.Duration, error) {

	consul, err := consulapi.NewClient(consulSpec)
	if err != nil {
		return 0, err
	}

	kv := consul.KV()

	var duration time.Duration

	for k, v := range config {
		took, err := kv.Put(&consulapi.KVPair{Key: k, Value: []byte(v)}, nil)
		if err != nil {
			return 0, fmt.Errorf("could not put a key, value: {%s, %s} to consul %+v due to %v", k, v, consulSpec, err)
		}
		duration += took.RequestTime
	}

	return duration, nil
}
