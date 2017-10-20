# Odor

A more efficient approach of niji platform.

## Usage

Clone this repository inside directory `$GOPATH/src/github.com/jlorgal/odor`.

Then you can use make to build the project.

```
Usage: make <command>
Commands:
  help:            Show this help information
  dep:             Ensure dependencies with dep tool
  build:           Build the application
  test-acceptance: Pass component tests
  release:         Create a new release (tag and release notes)
  run:             Launch the service with docker-compose (for testing purposes)
  clean:           Clean the project
  pipeline-pull:   Launch pipeline to handle a pull request
  pipeline-dev:    Launch pipeline to handle the merge of a pull request
  pipeline:        Launch the pipeline for the selected environment
  develenv-up:     Launch the development environment with a docker-compose of the service
  develenv-sh:     Access to a shell of a launched development environment
  develenv-down:   Stop the development environment
```

### Development environment

This project provides a development environment based on docker. It provides some benefits:
 - Common environment with required dependencies already installed
 - Pipelines to build the source code, pass the acceptance tests, package a docker image and publish the docker image.

Prerequirements:
 - Docker v17
 - GNU Make

```sh
# Launch the development environment
make develenv-up
# Access to the development environment
make develenv-sh
# Now you can launch any make task (e.g. build)
make build
```

## How to run odor

### Configuring iptables

Execute (as root) the following list of commands to configure IP tables.

```sh
sysctl -w net.ipv4.conf.all.forwarding=1
iptables -F
iptables -X
iptables -t filter -F
iptables -t filter -X
iptables -t raw -F
iptables -t raw -X
iptables -t nat -F
iptables -t nat -X
iptables -t mangle -F
iptables -t mangle -X
iptables -A FORWARD -i vlan1 -o vlan2 -j ACCEPT
iptables -A FORWARD -i vlan2 -o vlan1 -j ACCEPT
iptables -t raw -A PREROUTING -p tcp -m multiport --dports 80,443 --tcp-flags SYN SYN -j NFQUEUE --queue-num 1 --queue-bypass
iptables-save -c
```

The following command is responsible of forwarding all the SYN packets (to establish connections) for traffic HTTP and HTTPS to netfilter queue with number 1. The SYN packets enqueued are processed by odor:

```sh
iptables -t raw -A PREROUTING -p tcp \
         -m multiport --dports 80,443 --tcp-flags SYN SYN \
         -j NFQUEUE --queue-num 1 --queue-bypass
```

### Configuring odor

The file `config.json` contains the configuration for the odor service:

```json
{
    "logLevel": "INFO",
    "address": ":9000",
    "nfqueueID": 1,
    "filters": {
        "malware": [
            "213.211.198.0/24"
        ],
        "parental": [
            "91.109.250.0/24",
            "185.88.181.0/24"
        ],
        "adBlocking": [
            "216.58.192.0/19"
        ]
    }
}
```

The configuration includes the blocking rules:

| Category | Site | IP range |
|----------|------|----------|
| Parental control | www.888poker.es, www.888casino.es | 91.109.250.0/24 |
| Parental control | www.xvideos.com | 185.88.181.0/24 |
| Malware | www.eicar.org | 213.211.198.0/24 |
| Ad-Blocking | www.googleadservices.com | 216.58.192.0/19 |


### Starting up odor

```
cd build/bin
sudo ./odor
```

### Provisioning users in odor

Provision the user profile:

```sh
curl -X PUT \
  http://10.95.61.198:9000/users/34123456789 \
  -H 'content-type: application/json' \
  -d '{
 "msisdn": "34123456789",
 "antiPhising": true,
 "antiMalware": true,
 "parentalControl": true,
 "adBlocking": true,
 "captive": true
}'
```

Provision the mapping between IP and msisdn:

```sh
curl -X PUT \
  http://10.95.61.198:9000/ips/10.95.61.136 \
  -H 'content-type: application/json' \
  -d '{
 "ip": "10.95.61.136",
 "msisdn": "34123456789"
}'
```

It is possible to change the user profile (e.g. to disable all the settings):

```sh
curl -X PUT \
  http://10.95.61.198:9000/users/34123456789 \
  -H 'content-type: application/json' \
  -d '{
 "msisdn": "34123456789",
 "antiPhising": false,
 "antiMalware": false,
 "parentalControl": false,
 "adBlocking": false,
 "captive": false
}'
````
