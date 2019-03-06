package tradfri

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-coap"
	"github.com/eriklupander/dtls"
	"github.com/eriklupander/tradfri-go/model"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

type DtlsClient struct {
	peer  *dtls.Peer
	msgId uint16
}

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

func (dc *DtlsClient) PutDeviceDimming(deviceId string, dimming int) error {
	payload := fmt.Sprintf(`{ "3311": [{ "5851": %d }] }`, dimming)
	fmt.Printf("Payload is: %v", payload)
	resp := dc.call(dc.buildPUTMessage("/15001/"+deviceId, payload))
	fmt.Printf("Response: %+v", resp)
	return nil
}

func (dc *DtlsClient) PutDeviceColor(deviceId string, x, y int) error {
	payload := fmt.Sprintf(`{ "3311": [ {"5709": %d, "5710": %d}] }`, x, y)
	fmt.Printf("Payload is: %v", payload)
	resp := dc.call(dc.buildPUTMessage("/15001/"+deviceId, payload))
	fmt.Printf("Response: %+v", resp)
	return nil
}

func (dc *DtlsClient) PutDeviceColorRGB(deviceId, rgb string) error {
	payload := fmt.Sprintf(`{ "3311": [ {"5706": "%s"}] }`, rgb)
	fmt.Printf("Payload is: %v", payload)
	resp := dc.call(dc.buildPUTMessage("/15001/"+deviceId, payload))
	fmt.Printf("Response: %+v", resp)
	return nil
}

func (dc *DtlsClient) ListGroups() ([]model.Group, error) {
	resp := dc.call(dc.buildGETMessage("/15004"))
	groupIds := make([]int, 0)
	err := json.Unmarshal(resp.Payload, &groupIds)
	if err != nil {
		fmt.Println("Unable to parse groups list into JSON: " + err.Error())
		return nil, err
	}
	groups := make([]model.Group, 0)

	for _, id := range groupIds {
		group, _ := dc.GetGroup(strconv.Itoa(id))
		groups = append(groups, group)
	}
	return groups, nil
}

func (dc *DtlsClient) GetGroup(id string) (model.Group, error) {
	resp := dc.call(dc.buildGETMessage("/15004/" + id))
	group := &model.Group{}
	err := json.Unmarshal(resp.Payload, &group)
	if err != nil {
		return *group, err
	}
	return *group, nil
}

func (dc *DtlsClient) GetDevice(id string) (model.Device, error) {
	resp := dc.call(dc.buildGETMessage("/15001/" + id))
	device := &model.Device{}
	err := json.Unmarshal(resp.Payload, &device)
	if err != nil {
		return *device, err
	}
	return *device, nil
}

func (dc *DtlsClient) Get(id string) (coap.Message, error) {
	if !strings.HasPrefix(id, "/") {
		id = "/" + id
	}

	coapMessage := dc.call(dc.buildGETMessage(id))
	return coapMessage, nil
}

func (dc *DtlsClient) Put(id string, payload string) (coap.Message, error) {
	if !strings.HasPrefix(id, "/") {
		id = "/" + id
	}

	coapMessage := dc.call(dc.buildPUTMessage(id, payload))
	return coapMessage, nil
}

func (dc *DtlsClient) call(req coap.Message) coap.Message {
	fmt.Printf("Calling %v %v", req.Code.String(), req.PathString())
	data, err := req.MarshalBinary()
	if err != nil {
		panic(err)
	}
	err = dc.peer.Write(data)

	if err != nil {
		panic(err.Error())
	}

	respData, err := dc.peer.Read(time.Second)
	if err != nil {
		panic(err.Error())
	}
	msg, _ := coap.ParseMessage(respData)

	fmt.Printf("\nMessageID: %v\n", msg.MessageID)
	fmt.Printf("Type: %v\n", msg.Type)
	fmt.Printf("Code: %v\n", msg.Code)
	fmt.Printf("Token: %v\n", msg.Token)
	fmt.Printf("Payload: %v\n", string(msg.Payload))

	return msg
}

// Ref: https://community.openhab.org/t/ikea-tradfri-gateway/26135/148?u=kai
func (dc *DtlsClient) AuthExchange(clientId string) (TokenExchange, error) {

	// Prepare request
	dc.msgId++
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.POST,
		MessageID: dc.msgId,
		Payload:   []byte(fmt.Sprintf(`{"9090":"%s"}`, clientId)),
	}
	req.SetPathString("/15011/9063")

	// Send CoAP message for token exchange
	resp := dc.call(req)

	// Handle response and return
	token := TokenExchange{}
	err := json.Unmarshal(resp.Payload, &token)
	if err != nil {
		panic(err)
	}
	return token, nil
}

func (dc *DtlsClient) buildGETMessage(path string) coap.Message {
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

func (dc *DtlsClient) buildPUTMessage(path string, payload string) coap.Message {
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

func setupKeystore() {
	mks := dtls.NewKeystoreInMemory()
	dtls.SetKeyStores([]dtls.Keystore{mks})
	mks.AddKey(viper.GetString("CLIENT_ID"), []byte(viper.GetString("PRE_SHARED_KEY")))
}

type TokenExchange struct {
	Token          string `json:"9091"`
	TypeIdentifier string `json:"9029"`
}
