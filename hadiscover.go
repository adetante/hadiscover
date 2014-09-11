package main

import (
    "github.com/coreos/go-etcd/etcd"
    "text/template"
    "os"
    "os/exec"
    "flag"
)

var filename = flag.String("config","","Template config file used for HAproxy")
var name     = flag.String("name","","Base name used for backends")
var haproxy  = flag.String("ha","","Path to the `haproxy` executable")
var pidFile  = flag.String("pid","","Path to the haproxy pid file")
var etcdHost = flag.String("etcd","","etcd server(s)")
var etcdKey  = flag.String("key","","etcd root key to look for")

var configFile = ".haproxy.cfg"

func main(){
    flag.Parse()

    var etcdClient = etcd.NewClient([]string{*etcdHost})

    cfgFile,_ := os.Create(configFile)
    defer func() {
        cfgFile.Close();
        os.Remove(configFile)
    }()

    backends,_ := GetBackends(etcdClient,*etcdKey,*name)

    tpl, err := template.ParseFiles(*filename)
    if (err != nil){
        panic(err)
    }
    tpl.Execute(cfgFile, backends)

    cmd := exec.Command(*haproxy,"-f",configFile)
    cmd.Stdout = os.Stdout
    cmd.Run()
}