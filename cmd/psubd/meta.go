package main

import (
	"sync"

	"github.com/funkygao/gafka/ctx"
	"github.com/funkygao/gafka/zk"
)

type MetaStore interface {
	Start()
	Stop()
	Refresh()
	Clusters() []string
	ZkAddrs() string
	BrokerList() []string
	AuthPub(pubkey string) bool
	AuthSub(subkey string) bool
}

type zkMetaStore struct {
	brokerList []string
	zkcluster  *zk.ZkCluster
	mu         sync.Mutex
}

func newZkMetaStore(zone string, cluster string) MetaStore {
	zkzone := zk.NewZkZone(zk.DefaultConfig(zone, ctx.ZoneZkAddrs(zone)))
	return &zkMetaStore{
		zkcluster: zkzone.NewCluster(cluster),
	}
}

func (this *zkMetaStore) Start() {
	this.brokerList = this.zkcluster.BrokerList()
}

func (this *zkMetaStore) Stop() {
	this.zkcluster.Close()
}

func (this *zkMetaStore) Refresh() {
	this.brokerList = this.zkcluster.BrokerList()
}

func (this *zkMetaStore) BrokerList() []string {
	return this.brokerList
}

func (this *zkMetaStore) ZkAddrs() string {
	return this.zkcluster.ZkAddrs()
}

func (this *zkMetaStore) Clusters() []string {
	r := make([]string, 0)
	this.zkcluster.ZkZone().WithinClusters(func(name, path string) {
		r = append(r, name)
	})
	return r
}

func (this *zkMetaStore) AuthPub(pubkey string) (ok bool) {
	return true
}

func (this *zkMetaStore) AuthSub(subkey string) (ok bool) {
	return true
}
