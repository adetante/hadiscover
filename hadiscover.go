package main

import (
    "text/template"
    "os"
    "flag"
)

var filename = flag.String("config","","Template config file used for HAproxy")
var name     = flag.String("name","","Base name used for backends")    

type Backend struct{
    Name string
    Ip string
    Port string
}

func main(){
    flag.Parse()

    backends := []Backend{
        Backend{*name + "1","127.0.0.1","8001"},
        Backend{*name + "2","127.0.0.1","8002"},
        Backend{*name + "3","127.0.0.1","8003"},
    }

    tpl, _ := template.ParseFiles(*filename)
    tpl.Execute(os.Stdout, backends)

}