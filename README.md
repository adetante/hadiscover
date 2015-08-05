# hadiscover

This tool generates a [HAproxy](www.haproxy.org) configuration file based on [etcd](https://coreos.com/using-coreos/etcd/), and then reloads gracefully HAproxy.

hadiscover is listening on a specific directory in etcd, and for each changes it  re-generates the configuration and reloads graceully the server (using the `-sf` HAproxy flag).

It have been created to be used in parallel of my [Dockreg](https://github.com/adetante/dockreg) tool which does Docker container registration in etcd (see [my blog post](adetante.github.io/articles/service-discovery-with-docker-2)).

For more information and for build instruction, please read my post about [Service Discovery with HAproxy](http://adetante.github.io/articles/service-discovery-haproxy).

## Config file

hadiscover uses a [go text template](http://golang.org/pkg/text/template) to generate the haproxy configuration. For example:

```
global
    maxconn 4096

defaults
    log global
    mode    http
    option  httplog
    option  dontlognull
    retries 3
    redispatch
    maxconn 2000
    contimeout  5000
    clitimeout  50000
    srvtimeout  50000

frontend http-in
    bind *:8000
    default_backend http

backend http
{{range .}}     server {{.Name}} {{.Ip}}:{{.Port}} maxconn 32
{{end}}
```

The `backend http` part will be replaced by the list of available services retrieved in etcd.

The key name in etcd must have be formatted with the form `host:port`, for example:
`http://my-etcd-server:4001/keys/services/192.168.0.1:8000`


## Command line usage

```
haproxy --config templatePath --etcd etcdServersList --ha pathToHAcommand --key etcdKey
```

Where:

* **templatePath** is the path to the configuration template
* **etcdServersList** is the list of etcd servers, like `--etcd http://localhost:4001`
* **ha** is the path to the HAproxy executable
* **key** is the key to watch changes for
