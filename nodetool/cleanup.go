package nodetool

import (
	"github.com/yanniszark/go-nodetool/client"
)

func (t *Nodetool) Cleanup(keyspace string, tables []string) error {

	req := &client.JolokiaExecRequest{
		MBean:     storageServiceMBean,
		Operation: "cleanup",
	}

	if keyspace != "" {
		req.Arguments = append(req.Arguments, keyspace)
		for _, table := range tables {
			if table != "" {
				req.Arguments = append(req.Arguments, table)
			}
		}
	}

	// If operation succeeded we don't need to know anything else
	_, err := t.Client.Do(req)
	return err
}
