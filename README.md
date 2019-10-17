# tradfri-go

[![CircleCI](https://circleci.com/gh/eriklupander/tradfri-go.svg?style=svg)](https://circleci.com/gh/eriklupander/tradfri-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/eriklupander/tradfri-go)](https://goreportcard.com/report/github.com/eriklupander/tradfri-go)

Native Go implementation for talking CoAP to a [IKEA Trådfri](https://www.ikea.com/ms/en_US/usearch/?&query=tr%C3%A5dfri) gateway over DTLS 1.2.

**Note: The author is not in any way affiliated or related to IKEA, this is purely a hobby project!**

**Note 2: The application is being developed and have bugs and known issues!**

- CoAP implementation from [github.com/dustin/go-coap](https://github.com/dustin/go-coap)
- DTLS 1.2 support from [github.com/bocajim/dtls](https://github.com/bocajim/dtls)

This application is just stitching together the excellent work of [github.com/dustin](https://github.com/dustin) and [github.com/bocajim](https://github.com/bocajim) into a stand-alone application that can talk to a IKEA Trådfri gateway out-of-the-box without any dependencies on libcoap, openssl or similar libraries.

**Inspired by:**
- https://learn.pimoroni.com/tutorial/sandyj/controlling-ikea-tradfri-lights-from-your-pi
- https://github.com/glenndehaan/ikea-tradfri-coap-docs
- https://bitsex.net/software/2017/coap-endpoints-on-ikea-tradfri/

### Changelog
- 2019-06-19: gRPC support by [https://github.com/Age15990](https://github.com/Age15990)
- 2019-06-08: Configurable HTTP port by [https://github.com/Mirdinus](https://github.com/Mirdinus)
- 2019-04-02: Configuration redone by [https://github.com/Hades32](https://github.com/Hades32)
- 2019-04-01: Fixed issue with -authenticate
- 2019-03-10: Initial release

### Compatibility
tradfri-go has been tested against the following DTLS-enabled COAP servers:

- IKEA Trådfri Gateway using PSK after token exchange: OK
- Californicum COAP server on Scandium DTLS: OK

### Building
Uses go modules.

    export GO111MODULE=on
    go build -o tradfri-go
    
### PSK exchange
The Trådfri gateway has its Pre-shared key (PSK) printed on the bottom sticker. However, that PSK is only used for making an initial exchange where you specify a unique _Client_id_ and the original PSK, and you get a new PSK in return that you use for subsequent interactions with the gateway.

_tradfri-go_ supports this operation out of the box using the following command:

    > ./tradfri-go --authenticate --client_id=MyCoolID --psk=TheKeyAtTheBottomOfYourGateway --gateway_ip=<ip to your gateway>

The generated new PSK and settings used are stored in the current directory in the file "config.json", e.g:

    > cat config.json
    {
      "client_id": "MyCoolID",
      "gateway_address": "192.168.1.19:5684",
      "gateway_ip": "192.168.1.19",
      "pre_shared_key": "the generated psk goes here",
      "psk": "the generated psk goes here"
    }
    
_tradfri-go_ will try to read _config.json_ when starting up, and will in that case set the required properties accordingly.

If you don't feel like using _config.json_, you can either specify the configuration as command-line flags or using the following environment variables:

    ./tradfri-go --server --client_id MyCoolID122 --psk mynewkey --gateway_ip=192.168.1.19

or

    > export CLIENT_ID=MyCoolID1122
    > export PRE_SHARED_KEY=mynewkey
    > export GATEWAY_IP=192.168.1.19
    
Configuration is resolved in the following order of precedence:

config.json -> command-line arguments -> environment variables
    
### Determine gateway IP
_tradfri-go_ has no means of finding out the IP of the Gateway. I suggest checking your Router's list of connected devices and try to find an item starting with "GW-".

### Running in server mode
Server mode connects to your gateway and then publishes a really simple RESTful interface and a gRPC service for querying your gateway or mutating some state on bulbs etc:

    ./tradfri-go --server
    
Now, you can use the simple RESTful API provided by tradfri-go which returns more human-readable responses than the raw CoAP responses:

    > curl http://localhost:8080/api/device/65538 | jq .
    {
      "deviceMetadata": {
        "id": 65538,
        "name": "Färgglad",
        "vendor": "IKEA of Sweden",
        "type": "TRADFRI bulb E27 CWS opal 600lm"
      },
      "dimmer": 100,
      "xcolor": 30015,
      "ycolor": 26870,
      "rgbcolor": "f1e0b5",
      "power": true
    }
    
Or use one of the declarative endpoints to mutate the state of the bulb:

    > curl -X PUT -data '{"rgbcolor":"f1e0b5"}' http://localhost:8080/api/device/65538/rgb

If you want to use the gRPC service, implement your client like this:

    var client pb.TradfriServiceClient
	{
		conn, err := grpc.Dial("localhost:8081",
			grpc.WithInsecure(),
		)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		client = pb.NewTradfriServiceClient(conn)
	}

while importing `pb "github.com/eriklupander/tradfri-go/grpc_server/golang"`.

And, as simple as calling every other method in your code, request the server like this:

    resp, err := client.ListGroups(context.Background(), &pb.ListGroupsRequest{})

You can also install `grpcurl` and query the server via the command line:

    > grpcurl -plaintext localhost:8081 grpc_server.TradfriService/ListGroups
    > grpcurl -plaintext -d '{"id": 65552}' localhost:8081 grpc_server.TradfriService/TurnDeviceOn
    
Just like the client mode, the application will try to use clientId/PSK from _psk.key_ or using env vars.

### Running in client mode

Client mode lets you GET and PUT raw coap payloads to your gateway using the "-get" and "-put" args.

A few examples:

GET my bulb at /15001/65538:

    ./tradfri-go --get /15001/65538
    {"9019":1,"9001":"Färgglad","9002":1550336061,"9020":1551721891,"9003":65538,"9054":0,"5750":2,"3":{"0":"IKEA of Sweden","1":"TRADFRI bulb E27 CWS opal 600lm","2":"","3":"1.3.009","6":1},"3311":[{"5708":65279,"5850":1,"5851":100,"5707":53953,"5709":20316,"5710":8520,"5706":"8f2686","9003":0}]}

PUT that turns off the bulb at /15001/65538:
    
    ./tradfri-go --put /15001/65538 --payload '{ "3311": [{ "5850": 0 }] }'
    
PUT that turns on the bulb at /15001/65538 and sets dimmer to 200:
    
    ./tradfri-go --put /15001/65538 --payload '{ "3311": [{ "5850": 1, "5851": 200 }] }'
    
PUT that sets color of the bulb at /15001/65538 to purple and the dimmer to 100:
        
    ./tradfri-go --put /15001/65538 --payload '{ "3311": [{ "5706": "8f2686", "5851": 100 }] }'
    
The colors possible to set on the bulbs varies. The colors are in the CIE 1931 color space whose x/y values _in theory_ can be set using the 5709 and 5710 codes to values between 0 and 65535. You can't set arbitrary values due to how the CIE 1931 (yes, it's a standard from 1931!) works. Play around with the values, I havn't broken my full-color "TRADFRI bulb E27 CWS opal 600lm" yet...

# LICENSE

Uses MIT license, see [LICENSE](LICENSE)
