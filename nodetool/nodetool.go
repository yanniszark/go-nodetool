package nodetool

import (
	"net/http"
	"net/url"

	"github.com/yanniszark/go-nodetool/client"
)

const (
	storageServiceMBean = "org.apache.cassandra.db:type=StorageService"
)

type Nodetool struct {
	Client client.Interface
}

func New(client client.Interface) *Nodetool {
	return &Nodetool{
		Client: client,
	}
}

func NewFromURL(u *url.URL) *Nodetool {
	return New(client.New(u, &http.Client{}))
}
