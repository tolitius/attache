# attaché

## What

Interop between Consul API and Go maps:

```go
ConsulToMap(ConsulSpec{"consulHost:8500", "datacenterName", "token or empty string"},
            "rootPath")
```

returns a `map[string][]byte` of all the `k/v` under the `rootPath`.

`ConsulSpec` struct is:

```go
type ConsulSpec struct {
	Address    string
	Datacenter string
	Token      string
}
```

## Show me

Say we have this structure in Consul at `localhost:8500` and datacenter named `dc1`:

```json
{"hubble":
    {"store": "spacecraft://tape",
     "camera":
        {"mode": "color"},
     "mission":
        {"target": "Horsehead Nebula"}}}
```

attaché could read it all into a Go map by:

```go
ConsulToMap(ConsulSpec{"locahost:8500", "dc1", ""},
            "hubble")
```

will produce a Go map:

```go
{"hubble/store": "spacecraft://tape"
 "hubble/camera/mode": "color"
 "hubble/mission/target": "Horsehead Nebula"}
```

where keys are `[]byte`

## License

Copyright © 2017 tolitius

Distributed under the Eclipse Public License either version 1.0 or (at your option) any later version.
