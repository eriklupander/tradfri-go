package tradfri

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-coap"
	"github.com/eriklupander/tradfri-go/dtlscoap"
	"github.com/eriklupander/tradfri-go/model"
	"strconv"
	"strings"
)

type TradfriClient struct {
	dtlsclient *dtlscoap.DtlsClient
}

func NewTradfriClient() *TradfriClient {
	client := &TradfriClient{}
	client.dtlsclient = dtlscoap.NewDtlsClient()
	return client
}


func (dc *TradfriClient) PutDeviceDimming(deviceId string, dimming int) (model.Result, error)  {
	payload := fmt.Sprintf(`{ "3311": [{ "5851": %d }] }`, dimming)
	fmt.Printf("Payload is: %v", payload)
	resp, err := dc.dtlsclient.Call(dc.dtlsclient.BuildPUTMessage("/15001/"+deviceId, payload))
	if err != nil {
		return model.Result{}, err
	}
	fmt.Printf("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

func (dc *TradfriClient) PutDevicePower(deviceId string, power int) (model.Result, error)  {
	if !(power == 1 || power == 0) {
		return model.Result{}, fmt.Errorf("Invalid value for setting powered state, must be 1 or 0")
	}
	payload := fmt.Sprintf(`{ "3311": [{ "5850": %d }] }`, power)
	fmt.Printf("Payload is: %v", payload)
	resp, err := dc.dtlsclient.Call(dc.dtlsclient.BuildPUTMessage("/15001/"+deviceId, payload))
	if err != nil {
		return model.Result{}, err
	}
	fmt.Printf("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

func (dc *TradfriClient) PutDeviceState(deviceId string, power int, dimmer int, color string) (model.Result, error)  {
	if !(power == 1 || power == 0) {
		return model.Result{}, fmt.Errorf("Invalid value for setting powered state, must be 1 or 0")
	}
	payload := fmt.Sprintf(`{ "3311": [{ "5850": %d, "5851": %d}] }`, power, dimmer) // , "5706": "%s"
	fmt.Printf("Payload is: %v", payload)
	resp, err := dc.dtlsclient.Call(dc.dtlsclient.BuildPUTMessage("/15001/"+deviceId, payload))
	if err != nil {
		return model.Result{}, err
	}
	fmt.Printf("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

func (dc *TradfriClient) PutDeviceColor(deviceId string, x, y int) (model.Result, error)  {
	payload := fmt.Sprintf(`{ "3311": [ {"5709": %d, "5710": %d}] }`, x, y)
	fmt.Printf("Payload is: %v", payload)
	resp, err := dc.dtlsclient.Call(dc.dtlsclient.BuildPUTMessage("/15001/"+deviceId, payload))
	if err != nil {
		return model.Result{}, err
	}
	fmt.Printf("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

func (dc *TradfriClient) PutDeviceColorRGB(deviceId, rgb string) (model.Result, error)  {
	payload := fmt.Sprintf(`{ "3311": [ {"5706": "%s"}] }`, rgb)
	fmt.Printf("Payload is: %v", payload)
	resp, err := dc.dtlsclient.Call(dc.dtlsclient.BuildPUTMessage("/15001/"+deviceId, payload))
	if err != nil {
		return model.Result{}, err
	}
	fmt.Printf("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

func (dc *TradfriClient) ListGroups() ([]model.Group, error) {
	groups := make([]model.Group, 0)

	resp, err := dc.dtlsclient.Call(dc.dtlsclient.BuildGETMessage("/15004"))
	if err != nil {
		fmt.Println("Unable to call Trådfri: " + err.Error())
		return groups, err
	}

	groupIds := make([]int, 0)
	err = json.Unmarshal(resp.Payload, &groupIds)
	if err != nil {
		fmt.Println("Unable to parse groups list into JSON: " + err.Error())
		return groups, err
	}


	for _, id := range groupIds {
		group, _ := dc.GetGroup(strconv.Itoa(id))
		groups = append(groups, group)
	}
	return groups, nil
}

func (dc *TradfriClient) GetGroup(id string) (model.Group, error) {
	resp, err := dc.dtlsclient.Call(dc.dtlsclient.BuildGETMessage("/15004/" + id))
	group := &model.Group{}
	if err != nil {
		return *group, err
	}

	err = json.Unmarshal(resp.Payload, &group)
	if err != nil {
		return *group, err
	}
	return *group, nil
}

func (dc *TradfriClient) GetDevice(id string) (model.Device, error) {
	device := &model.Device{}

	resp, err := dc.dtlsclient.Call(dc.dtlsclient.BuildGETMessage("/15001/" + id))
	if err != nil {
		return *device, err
	}
	err = json.Unmarshal(resp.Payload, &device)
	if err != nil {
		return *device, err
	}
	return *device, nil
}

func (dc *TradfriClient) Get(id string) (coap.Message, error) {
	if !strings.HasPrefix(id, "/") {
		id = "/" + id
	}
	return dc.dtlsclient.Call(dc.dtlsclient.BuildGETMessage(id))
}

func (dc *TradfriClient) Put(id string, payload string) (coap.Message, error) {
	if !strings.HasPrefix(id, "/") {
		id = "/" + id
	}
	return dc.dtlsclient.Call(dc.dtlsclient.BuildPUTMessage(id, payload))
}



// Ref: https://community.openhab.org/t/ikea-tradfri-gateway/26135/148?u=kai
func (dc *TradfriClient) AuthExchange(clientId string) (TokenExchange, error) {

	req := dc.dtlsclient.BuildPOSTMessage("/15011/9063", fmt.Sprintf(`{"9090":"%s"}`, clientId))

	// Send CoAP message for token exchange
	resp, err := dc.dtlsclient.Call(req)

	// Handle response and return
	token := TokenExchange{}
	err = json.Unmarshal(resp.Payload, &token)
	if err != nil {
		panic(err)
	}
	return token, nil
}




type TokenExchange struct {
	Token          string `json:"9091"`
	TypeIdentifier string `json:"9029"`
}