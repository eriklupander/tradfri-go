package dtlscoap

import (
	"fmt"
	"github.com/dustin/go-coap"
	"github.com/eriklupander/dtls"
	"github.com/spf13/viper"
	"time"
)

// DtlsClient provides an domain-agnostic CoAP-client with DTLS transport.
type DtlsClient struct {
	peer  *dtls.Peer
	msgId uint16
}

// NewDtlsClient acts as factory function, returns a pointer to a connected (or will panic) DtlsClient.
func NewDtlsClient() *DtlsClient {
	client := &DtlsClient{}
	client.connect()
	return client
}

func (dc *DtlsClient) connect() {
	setupKeystore()

	listener, err := dtls.NewUdpListener(":0", time.Second*900)
	if err != nil {
		panic(err.Error())
	}

	gatewayIp := viper.GetString("GATEWAY_IP") + ":5684"

	peerParams := &dtls.PeerParams{
		Addr:             gatewayIp,
		Identity:         viper.GetString("CLIENT_ID"),
		HandshakeTimeout: time.Second * 30}
	fmt.Printf("Connecting to peer at %v\n", gatewayIp)

	dc.peer, err = listener.AddPeerWithParams(peerParams)
	if err != nil {
		panic("Unable to connect to Gateway at " + gatewayIp + ": " + err.Error())
	}
	dc.peer.UseQueue(true)
	fmt.Printf("DTLS connection established to %v\n", gatewayIp)
}

// Call writes the supplied coap.Message to the peer
func (dc *DtlsClient) Call(req coap.Message) (coap.Message, error) {
	fmt.Printf("Calling %v %v", req.Code.String(), req.PathString())
	data, err := req.MarshalBinary()
	if err != nil {
		return coap.Message{}, err
	}
	err = dc.peer.Write(data)

	if err != nil {
		return coap.Message{}, err
	}

	respData, err := dc.peer.Read(time.Second)
	if err != nil {
		return coap.Message{}, err
	}

	msg, err := coap.ParseMessage(respData)
	if err != nil {
		return coap.Message{}, err
	}

	fmt.Printf("\nMessageID: %v\n", msg.MessageID)
	fmt.Printf("Type: %v\n", msg.Type)
	fmt.Printf("Code: %v\n", msg.Code)
	fmt.Printf("Token: %v\n", msg.Token)
	fmt.Printf("Payload: %v\n", string(msg.Payload))

	return msg, nil
}

// BuildGETMessage produces a CoAP GET message with the next msgId set.
func (dc *DtlsClient) BuildGETMessage(path string) coap.Message {
	dc.msgId++
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: dc.msgId,
	}
	req.SetPathString(path)
	return req
}

//req.SetOption(coap.ETag, "weetag")
//req.SetOption(coap.MaxAge, 3)

// BuildPUTMessage produces a CoAP PUT message with the next msgId set.
func (dc *DtlsClient) BuildPUTMessage(path string, payload string) coap.Message {
	dc.msgId++

	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.PUT,
		MessageID: dc.msgId,
		Payload:   []byte(payload),
	}
	req.SetPathString(path)

	return req
}

// BuildPOSTMessage produces a CoAP POST message with the next msgId set.
func (dc *DtlsClient) BuildPOSTMessage(path string, payload string) coap.Message {
	dc.msgId++

	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.POST,
		MessageID: dc.msgId,
		Payload:   []byte(payload),
	}
	req.SetPathString(path)

	return req
}

func setupKeystore() {
	mks := dtls.NewKeystoreInMemory()
	dtls.SetKeyStores([]dtls.Keystore{mks})
	mks.AddKey(viper.GetString("CLIENT_ID"), []byte(viper.GetString("PRE_SHARED_KEY")))
}
