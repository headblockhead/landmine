package landmine

import (
	"github.com/headblockhead/landmine/core"
	"go.etcd.io/etcd/client/v3"
	"time"
)

type Landmine struct {
	Config     core.Config
	EtcdClient *clientv3.Client
	//Workspaces []core.Workspace
}

func NewLandmine(config core.Config) *Landmine {
	l := new(Landmine)
	l.Config = config

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.EtcdHostname},
		DialTimeout: 5 * time.Second, // TODO
	})
	if err != nil {
		panic(err) // TODO
	}
	l.EtcdClient = client

	return l
}
