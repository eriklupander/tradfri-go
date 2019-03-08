package main

import (
	"flag"
	"fmt"
	"github.com/eriklupander/tradfri-go/router"
	"github.com/eriklupander/tradfri-go/tradfri"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var serverMode, authenticate *bool

func init() {
	// dtls.SetLogLevel("debug")

	// read clientId / PSK from file if possible
	ok := resolveClientIdAndPSKFromFile()
	viper.AutomaticEnv()
	if !ok {

		if viper.GetString("PRE_SHARED_KEY") == "" {
			panic("Unable to resolve PRE_SHARED_KEY from env-var or psk.key file")
		}
		if viper.GetString("CLIENT_ID") == "" {
			panic("Unable to resolve CLIENT_ID from env-var or psk.key file")
		}
	}

}

func resolveClientIdAndPSKFromFile() bool {
	data, err := ioutil.ReadFile("psk.key")
	if err != nil {
		fmt.Println("Could not find psk.key in current directory, trying to use env-vars CLIENT_ID and PRE_SHARED_KEY instead")
		return false
	}
	str := string(data)
	if !strings.Contains(str, "=") {
		panic("Invalid psk.key, must contain a single line with [clientId]=[PSK]. No = found")
	}
	parts := strings.Split(str, "=")
	if len(parts) != 2 {
		panic("Invalid psk.key, must contain a single line with [clientId]=[PSK], could not split into two parts")
	}
	viper.Set("CLIENT_ID", parts[0])
	viper.Set("PRE_SHARED_KEY", parts[1])
	return true
}

func main() {
	// Really ugly flag / config handling... :(
	serverMode = flag.Bool("server", false, "Start in server mode?")
	gatewayIp := flag.String("gateway", "", "ip to your gateway. No protocol or port here!")
	authenticate = flag.Bool("authenticate", false, "Perform PSK exchange")
	psk := flag.String("psk", "", "Pre-shared key on bottom of Gateway")
	clientId := flag.String("clientId", "", "Your client id, make something up or use the NNN-NNN-NNN on the bottom of your Gateway")
	get := flag.String("get", "", "URL to GET")
	put := flag.String("put", "", "URL to PUT")
	payload := flag.String("payload", "", "payload for PUT")
	flag.Parse()

	resolveGatewayIP(gatewayIp)

	handleSigterm(nil)
	if *serverMode {
		fmt.Println("Running in server mode on :8080")
		go router.SetupChi(tradfri.NewTradfriClient())

		wg := sync.WaitGroup{}
		wg.Add(1)
		wg.Wait()
	} else {
		// client mode
		if *authenticate {
			performTokenExchange(clientId, psk)
		} else if *get != "" {
			resp, _ := tradfri.NewTradfriClient().Get(*get)
			fmt.Printf("%v", string(resp.Payload))
		} else if *put != "" {
			resp, _ := tradfri.NewTradfriClient().Put(*put, *payload)
			fmt.Printf("%v", string(resp.Payload))
		} else {
			fmt.Println("No client operation was specified, supported one(s) are: authenticate")
		}
	}

}

func resolveGatewayIP(gatewayIp *string) {
	if viper.GetString("GATEWAY_IP") == "" {
		if *gatewayIp != "" {
			viper.Set("GATEWAY_IP", *gatewayIp)
		} else {
			panic("Unable to resolve gateway IP from either of env var GATEWAY_IP or flag -gateway")
		}
	}
}

func performTokenExchange(clientId *string, psk *string) {
	if len(*clientId) < 1 || len(*psk) < 10 {
		panic("Both clientId and psk args must be specified")
	}

	// Note that we hard-code the Client_identity here before creating the DTLS client,
	// required when performing token exchange
	viper.Set("CLIENT_ID", "Client_identity")
	viper.Set("PRE_SHARED_KEY", *psk)
	dtlsClient := tradfri.NewTradfriClient()

	authToken, err := dtlsClient.AuthExchange(*clientId)
	if err != nil {
		panic(err.Error())
	}
	d1 := []byte(*clientId + "=" + authToken.Token)
	err = ioutil.WriteFile("psk.key", d1, 0644)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Your new PSK and clientId has been written to psk.key, keep this file safe!")
}

// Handles Ctrl+C or most other means of "controlled" shutdown gracefully. Invokes the supplied func before exiting.
func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if handleExit != nil {
			handleExit()
		}
		os.Exit(1)
	}()
}
