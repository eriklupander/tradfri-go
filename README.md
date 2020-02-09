# tradfri-go(-blind-server)

[![CircleCI](https://circleci.com/gh/eriklupander/tradfri-go.svg?style=svg)](https://circleci.com/gh/eriklupander/tradfri-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/eriklupander/tradfri-go)](https://goreportcard.com/report/github.com/eriklupander/tradfri-go)

Native Go implementation for talking CoAP to a [IKEA Trådfri](https://www.ikea.com/ms/en_US/usearch/?&query=tr%C3%A5dfri) gateway over DTLS 1.2.

**Note: The author is not in any way affiliated or related to IKEA, this is purely a hobby project!**

**Note 2: The application is being developed and have bugs and known issues!**

- CoAP implementation from [github.com/dustin/go-coap](https://github.com/dustin/go-coap)
- DTLS 1.2 support from [github.com/bocajim/dtls](https://github.com/bocajim/dtls)

This application is just stitching together the excellent work of [github.com/dustin](https://github.com/dustin) and [github.com/bocajim](https://github.com/bocajim) into a stand-alone application that can talk to a IKEA Trådfri gateway out-of-the-box without any dependencies on libcoap, openssl or similar libraries.

**Inspired by:**
https://github.com/eriklupander/tradfri-go

### Changelog
- 2020-02-10: Removing all light, remote and outlet control. Adding support for Ikea Blinds
- 2020-01-16: Updated logging to use logrus with configurable log level in config.json.
- 2019-06-19: gRPC support by [https://github.com/Age15990](https://github.com/Age15990)
- 2019-06-08: Configurable HTTP port by [https://github.com/Mirdinus](https://github.com/Mirdinus)
- 2019-04-02: Configuration redone by [https://github.com/Hades32](https://github.com/Hades32)
- 2019-04-01: Fixed issue with -authenticate
- 2019-03-10: Initial release

### Compatibility
tradfri-go(-blind-server) has been tested against the following DTLS-enabled COAP servers:

- IKEA Trådfri Gateway using PSK after token exchange: OK

### Building
Uses go modules.

    export GO111MODULE=on
    go build -o tradfri-go

or Use a golang:1.13 docker container to build tradfri-go(-blind-server)

    ./build.sh .

### PSK exchange
The Trådfri gateway has its Pre-shared key (PSK) printed on the bottom sticker. However, that PSK is only used for making an initial exchange where you specify a unique _Client_id_ and the original PSK, and you get a new PSK in return that you use for subsequent interactions with the gateway.

_tradfri-go(-blind-server)_ supports this operation out of the box using the following command:

    > ./tradfri-go --authenticate --client_id=MyCoolID --psk=TheKeyAtTheBottomOfYourGateway --gateway_ip=<ip to your gateway>

The generated new PSK and settings used are stored in the current directory in the file "config.json", e.g:

    > cat config.json
    {
      "client_id": "MyCoolID",
      "gateway_address": "192.168.1.19:5684",
      "gateway_ip": "192.168.1.19",
      "pre_shared_key": "the generated psk goes here",
      "psk": "the generated psk goes here",
      "loglevel":"info"
    }
    
_tradfri-go(-blind-server)_ will try to read _config.json_ when starting up, and will in that case set the required properties accordingly.

If you don't feel like using _config.json_, you can either specify the configuration as command-line flags or using the following environment variables:

    ./tradfri-go --server --client_id MyCoolID122 --psk mynewkey --gateway_ip=192.168.1.19

or

    > export CLIENT_ID=MyCoolID1122
    > export PRE_SHARED_KEY=mynewkey
    > export GATEWAY_IP=192.168.1.19
    
Configuration is resolved in the following order of precedence:

config.json -> command-line arguments -> environment variables
    
### Determine gateway IP
_tradfri-go(-blind-server)_ has no means of finding out the IP of the Gateway. I suggest checking your Router's list of connected devices and try to find an item starting with "GW-".

### Running in server mode
Server mode connects to your gateway and then publishes a really simple RESTful interface and a gRPC service for querying your gateway or mutating some state on blinds:

    docker-compose up -d
    
Now, you can use the simple RESTful API provided by tradfri-go which returns more human-readable responses than the raw CoAP responses:

    > curl http://localhost:8080/api/device/65541 | jq .
    {
      "deviceMetadata": {
        "id": 65541,
        "name": "Living Room Left Blind",
        "vendor": "IKEA of Sweden",
        "type": "FYRTUR block-out roller blind",
        "battery": 90
      },
      "position": 20
    }
    
Or use one of the declarative endpoints to mutate the state of the blind:

    > curl -X PUT -data '{"positioning": 20}' http://localhost:8080/api/device/65541/position


# LICENSE

Uses MIT license, see [LICENSE](LICENSE)
