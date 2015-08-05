package main

import (
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"strings"
)

type Backend struct {
	Name string
	Ip   string
	Port string
}

func GetBackends(client *etcd.Client, service, backendName string) ([]Backend, error) {

	resp, err := client.Get(service, false, true)
	if err != nil {
		log.Println("Error when reading etcd: ", err)
		return nil, err
	} else {
		backends := make([]Backend, len(resp.Node.Nodes))
		for index, element := range resp.Node.Nodes {

			key := (*element).Key // key format is: /service/IP:PORT
			service := strings.Split(key[strings.LastIndex(key, "/")+1:], ":")

			backends[index] = Backend{Name: fmt.Sprintf("back-%v", index), Ip: service[0], Port: service[1]}
		}
		return backends, nil
	}

}
