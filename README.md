# go-nodetool

A golang library to manage Cassandra / Scylla clusters.

## Dependencies 

This library needs the [Jolokia](https://jolokia.org/) JMX<->HTTP bridge to work.
It uses HTTP requests to the Jolokia server in order to read/write/exec Cassandra MBeans.
To see how you can use Jolokia in Cassandra, check [this guide](https://blog.pythian.com/two-easy-ways-poll-apache-cassandra-metrics-using-jmx-http-bridge/).


## Usage

`go-nodetool` performs administrative functions on a Cassandra / Scylla cluster. The commands available at the moment are:

|  Command  |  Explanation  |
|   :---:     |     :---:       |
|  `Status`   |  information about the cluster |
|  `Version`  |  version of the specific instance  |
| `OperationMode` | operation mode of the specific instance |
| `Decommission`  | decommission the specific instance |


```go
import (
	"net/url"
	"github.com/yanniszark/go-nodetool/nodetool"
)

url, _ = url.Parse("127.0.0.1:8778/jolokia/")
nt := nodetool.NewFromURL(url)

nodeMap, err := nt.Status()
```

## Contributing

Adding new commands is pretty straightforward. If you wish to extend this library with more commands, open an issue to request an enhancement or a PR to submit one. Use the existing commands as a guide on how to implement new ones and include unit tests.

## Acknowledgement

A signification portion of code from this library comes from [Jetstack Navigator](https://github.com/jetstack/navigator/tree/67d871f4d0d4e765b0c46425e139d4517ea26068/pkg/cassandra/nodetool). Modifications were made so that the library could be extended more easily. 
