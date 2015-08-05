package main

import (
	"flag"
	"github.com/coreos/go-etcd/etcd"
	"log"
)

var filename = flag.String("config", "./haproxy.cfg.tpl", "Template config file used for HAproxy")
var name = flag.String("name", "back", "Base name used for backends")
var haproxy = flag.String("ha", "/usr/bin/haproxy", "Path to the `haproxy` executable")
var etcdHost = flag.String("etcd", "http://localhost:4001", "etcd server(s)")
var etcdKey = flag.String("key", "services", "etcd root key to look for")

var configFile = ".haproxy.cfg"

func reloadConf(etcdClient *etcd.Client) error {
	backends, _ := GetBackends(etcdClient, *etcdKey, *name)

	err := createConfigFile(backends, *filename, configFile)
	if err != nil {
		log.Println("Cannot generate haproxy configuration: ", err)
		return err
	}
	return reloadHAproxy(*haproxy, configFile)
}

func main() {
	flag.Parse()

	var etcdClient = etcd.NewClient([]string{*etcdHost})
	err := reloadConf(etcdClient)
	if err != nil {
		log.Println("Cannot reload haproxy: ", err)
	}

	changeChan := make(chan *etcd.Response)
	stopChan := make(chan bool)

	go func() {
		for msg := range changeChan {
			reload := (msg.PrevNode == nil) || (msg.PrevNode.Key != msg.Node.Key) || (msg.Action != "set")
			if reload {
				err := reloadConf(etcdClient)
				if err != nil {
					log.Println("Cannot reload haproxy: ", err, msg)
				}
			}
		}
	}()

	log.Println("Start watching changes in etcd")
	if _, err := etcdClient.Watch(*etcdKey, 0, true, changeChan, stopChan); err != nil {
		log.Println("Cannot register watcher for changes in etcd: ", err)
	}

}
