package tradfri

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-coap"
	"github.com/eriklupander/tradfri-go/dtlscoap"
	"github.com/eriklupander/tradfri-go/model"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// Client provides a declarative API for sending CoAP messages to the gateway over DTLS.
type Client struct {
	dtlsclient *dtlscoap.DtlsClient
}

// NewTradfriClient creates a new instance of Client, including initiating the DTLS client.
func NewTradfriClient(gatewayAddress, clientID, psk string) *Client {
	client := &Client{}
	client.dtlsclient = dtlscoap.NewDtlsClient(gatewayAddress, clientID, psk)
	return client
}

// PutDeviceDimming sets the dimming property (0-255) of the specified device.
// The device must be a bulb supporting dimming, otherwise the call if ineffectual.
func (tc *Client) PutDeviceDimming(deviceId string, dimming int) (model.Result, error) {
	payload := fmt.Sprintf(`{ "3311": [{ "5851": %d }] }`, dimming)
	logrus.Infof("Payload is: %v", payload)
	resp, err := tc.Call(tc.dtlsclient.BuildPUTMessage("/15001/"+deviceId, payload))
	if err != nil {
		return model.Result{}, err
	}
	logrus.Infof("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

// PutDevicePower switches the power state of the specified device to on (1) or off (0)
func (tc *Client) PutDevicePower(deviceId string, power int) (model.Result, error) {
	if !(power == 1 || power == 0) {
		return model.Result{}, fmt.Errorf("invalid value for setting power state, must be 1 or 0")
	}
	payload := fmt.Sprintf(`{ "3311": [{ "5850": %d }] }`, power)
	logrus.Infof("Payload is: %v", payload)
	resp, err := tc.Call(tc.dtlsclient.BuildPUTMessage("/15001/"+deviceId, payload))
	if err != nil {
		return model.Result{}, err
	}
	logrus.Infof("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

// PutDeviceState allows changing both power (1 or 0) and dimmer (0-255) for a given device with one command.
func (tc *Client) PutDeviceState(deviceId string, power int, dimmer int, color string) (model.Result, error) {
	if !(power == 1 || power == 0) {
		return model.Result{}, fmt.Errorf("invalid value for setting power state, must be 1 or 0")
	}
	payload := fmt.Sprintf(`{ "3311": [{ "5850": %d, "5851": %d}] }`, power, dimmer) // , "5706": "%s"
	logrus.Infof("Payload is: %v", payload)
	resp, err := tc.Call(tc.dtlsclient.BuildPUTMessage("/15001/"+deviceId, payload))
	if err != nil {
		return model.Result{}, err
	}
	logrus.Infof("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

// PutDeviceColor sets the CIE 1931 color space x/y color, x and y must be between 0-65536 but note that
// many combinations won't work. See CIE 1931 for more details.
func (tc *Client) PutDeviceColor(deviceId string, x, y int) (model.Result, error) {
	payload := fmt.Sprintf(`{ "3311": [ {"5709": %d, "5710": %d}] }`, x, y)
	logrus.Infof("Payload is: %v", payload)
	resp, err := tc.Call(tc.dtlsclient.BuildPUTMessage("/15001/"+deviceId, payload))
	if err != nil {
		return model.Result{}, err
	}
	logrus.Infof("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

// PutDeviceColorRGB sets the color of the bulb using RGB hex string such as 8f2686 (purple). Note that
// many colors doesn't seem to work. Not sure how the IKEA bulbs with color support works.
func (tc *Client) PutDeviceColorRGB(deviceId, rgb string) (model.Result, error) {
	payload := fmt.Sprintf(`{ "3311": [ {"5706": "%s"}] }`, rgb)
	logrus.Infof("Payload is: %v", payload)
	resp, err := tc.Call(tc.dtlsclient.BuildPUTMessage("/15001/"+deviceId, payload))
	if err != nil {
		return model.Result{}, err
	}
	logrus.Infof("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

// PutDevicePositioning sets the positioning property (0-100) of the specified device.
func (tc *Client) PutDevicePositioning(deviceId string, positioning float32) (model.Result, error) {
	payload := fmt.Sprintf(`{ "15015": [{ "5536": %f }] }`, positioning)
	logrus.Infof("Payload is: %v", payload)
	resp, err := tc.Call(tc.dtlsclient.BuildPUTMessage("/15001/"+deviceId, payload))
	if err != nil {
		return model.Result{}, err
	}
	logrus.Infof("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

// ListGroups lists all groups
func (tc *Client) ListGroups() ([]model.Group, error) {
	groups := make([]model.Group, 0)

	resp, err := tc.Call(tc.dtlsclient.BuildGETMessage("/15004"))
	if err != nil {
		logrus.WithError(err).Error("Unable to call Trådfri Gateway")
		return groups, err
	}

	groupIds := make([]int, 0)
	err = json.Unmarshal(resp.Payload, &groupIds)
	if err != nil {
		logrus.Info("Unable to parse groups list into JSON: " + err.Error())
		return groups, err
	}

	for _, id := range groupIds {
		group, _ := tc.GetGroup(strconv.Itoa(id))
		groups = append(groups, group)
	}
	return groups, nil
}

// GetGroup gets the JSON representation of the specified group.
func (tc *Client) GetGroup(id string) (model.Group, error) {
	resp, err := tc.Call(tc.dtlsclient.BuildGETMessage("/15004/" + id))
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

// GetDevice gets the JSON representation of the specified device.
func (tc *Client) GetDevice(id string) (model.Device, error) {
	device := &model.Device{}

	resp, err := tc.Call(tc.dtlsclient.BuildGETMessage("/15001/" + id))
	if err != nil {
		return *device, err
	}
	err = json.Unmarshal(resp.Payload, &device)
	if err != nil {
		return *device, err
	}
	return *device, nil
}

// Get gets whatever is identified by the passed ID string.
func (tc *Client) Get(id string) (coap.Message, error) {
	if !strings.HasPrefix(id, "/") {
		id = "/" + id
	}
	return tc.Call(tc.dtlsclient.BuildGETMessage(id))
}

// Put puts the payload for whatever is identified by the passed ID string.
func (tc *Client) Put(id string, payload string) (coap.Message, error) {
	if !strings.HasPrefix(id, "/") {
		id = "/" + id
	}
	return tc.Call(tc.dtlsclient.BuildPUTMessage(id, payload))
}

// AuthExchange performs the initial PSK exchange.
// see ref: https://community.openhab.org/t/ikea-tradfri-gateway/26135/148?u=kai
func (tc *Client) AuthExchange(clientId string) (model.TokenExchange, error) {

	req := tc.dtlsclient.BuildPOSTMessage("/15011/9063", fmt.Sprintf(`{"9090":"%s"}`, clientId))

	// Send CoAP message for token exchange
	resp, err := tc.Call(req)
	if err != nil {
		logrus.WithError(err).Fatal("error performing call to Gateway for token exchange")
	}

	// Handle response and return
	token := model.TokenExchange{}
	err = json.Unmarshal(resp.Payload, &token)
	if err != nil {
		logrus.WithError(err).Fatal("error unmarhsalling response from Gateway for token exchange")
	}
	return token, nil
}

// Call is just a proxy to the underlying DtlsClient Call
func (tc *Client) Call(msg coap.Message) (coap.Message, error) {
	return tc.dtlsclient.Call(msg)
}
