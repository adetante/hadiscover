package main

import (
    "github.com/coreos/go-etcd/etcd"
    "fmt"
    "flag"
    "os"
)

var filename = flag.String("config","./haproxy.cfg.tpl","Template config file used for HAproxy")
var name     = flag.String("name","back","Base name used for backends")
var haproxy  = flag.String("ha","/usr/bin/haproxy","Path to the `haproxy` executable")
var pidFile  = flag.String("pid","/var/run/haproxy.pid","Path to the haproxy pid file")
var etcdHost = flag.String("etcd","http://localhost:4001","etcd server(s)")
var etcdKey  = flag.String("key","services","etcd root key to look for")

var configFile = ".haproxy.cfg"

func reloadConf(etcdClient *etcd.Client)(error){
    backends,_ := GetBackends(etcdClient,*etcdKey,*name)

    err := createConfigFile(backends, *filename, configFile)
    if(err != nil){
        fmt.Fprintln(os.Stderr,"Cannot generate haproxy configuration: ",err)
        return err
    }
    return reloadHAproxy(*haproxy, configFile)
}

func main(){
    defer func(){
        os.Remove(configFile)
    }()
    flag.Parse()

    var etcdClient = etcd.NewClient([]string{*etcdHost})
    err := reloadConf(etcdClient)
    if(err != nil){
        fmt.Fprintln(os.Stderr,"Cannot reload haproxy: ",err)
    }

    changeChan  := make(chan *etcd.Response)
    stopChan    := make(chan bool)

    go func(){
        for msg := range changeChan {
            err := reloadConf(etcdClient)
            if(err != nil){
                fmt.Fprintln(os.Stderr,"Cannot reload haproxy: ",err,msg)
            }
        }
    }()

    fmt.Println("Watching...")
    if _, err := etcdClient.Watch(*etcdKey,0, true, changeChan, stopChan) ; err != nil{
        fmt.Fprintln(os.Stderr,"Cannot watch changes in etcd: ",err)
    }    

}