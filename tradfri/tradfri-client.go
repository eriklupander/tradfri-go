package tradfri

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/dustin/go-coap"
	"github.com/eriklupander/tradfri-go/dtlscoap"
	"github.com/eriklupander/tradfri-go/model"
	"github.com/sirupsen/logrus"
)

const (
	/** A "normal" remote */
	DeviceTypeRemote = iota
	/**
	 * A remote which has been paired with another remote.
	 * See https://www.reddit.com/r/tradfri/comments/6x1miq for details
	 */
	DeviceTypeSlaveRemote
	/** A lightbulb */
	DeviceTypeLightbulb
	/** A smart plug */
	DeviceTypePlug
	/** A motion sensor (currently unsupported) */
	DeviceTypeMotionSensor
	/** A signal repeater */
	DeviceTypeSignalRepeater
	/** A smart blind */
	DeviceTypeBlind
	/** Symfonisk Remote */
	DeviceTypeSoundRemote
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

// Yo who decided it would be a good idea to have the deviceId be an int in all the models, but here every function wants it as a string, its stupid
// But i aint changin it because it could break code that depends on these functions (It would be a great change tho)

// PutDeviceDimming sets the dimming property (0-255) of the specified device.
// The device must be a bulb supporting dimming, otherwise the call if ineffectual.
func (tc *Client) PutDeviceDimming(deviceId int, dimming int) (model.Result, error) {
	payload := fmt.Sprintf(`{ "3311": [{ "5851": %d }] }`, dimming)
	logrus.Infof("Payload is: %v", payload)
	resp, err := tc.Call(tc.dtlsclient.BuildPUTMessage(toDeviceUri(deviceId), payload))
	if err != nil {
		return model.Result{}, err
	}
	logrus.Infof("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

// PutDevicePower switches the power state of the specified device to on (1) or off (0)
func (tc *Client) PutDevicePower(deviceId int, power int) (model.Result, error) {
	if !(power == 1 || power == 0) {
		return model.Result{}, fmt.Errorf("invalid value for setting power state, must be 1 or 0")
	}
	payload := fmt.Sprintf(`{ "3311": [{ "5850": %d }] }`, power)
	logrus.Infof("Payload is: %v", payload)
	resp, err := tc.Call(tc.dtlsclient.BuildPUTMessage(toDeviceUri(deviceId), payload))
	if err != nil {
		return model.Result{}, err
	}
	logrus.Infof("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

// PutDeviceState allows changing both power (1 or 0) and dimmer (0-255) for a given device with one command.
func (tc *Client) PutDeviceState(deviceId int, power int, dimmer int) (model.Result, error) {
	if !(power == 1 || power == 0) {
		return model.Result{}, fmt.Errorf("invalid value for setting power state, must be 1 or 0")
	}
	payload := fmt.Sprintf(`{ "3311": [{ "5850": %d, "5851": %d}] }`, power, dimmer) // , "5706": "%s"
	logrus.Infof("Payload is: %v", payload)
	resp, err := tc.Call(tc.dtlsclient.BuildPUTMessage(toDeviceUri(deviceId), payload))
	if err != nil {
		return model.Result{}, err
	}
	logrus.Infof("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

// PutDeviceColor sets the CIE 1931 color space x/y color, x and y must be between 0-65536 but note that
// many combinations won't work. See CIE 1931 for more details.
// It is not recommended to use these values to set colors, as it is often not supported by the gateway and is intended for internal use.
func (tc *Client) PutDeviceColor(deviceId int, x, y int) (model.Result, error) {
	return tc.PutDeviceColorTimed(deviceId, x, y, 500)
}

// PutDeviceColorTimed does the same as PutDeviceColor but it gives you the ability to change the speed at which the color changes
func (tc *Client) PutDeviceColorTimed(deviceId int, x, y int, transitionTimeMS int) (model.Result, error) {
	payload := fmt.Sprintf(`{ "3311": [ {"5709": %d, "5710": %d, "5712": %d}] }`, x, y, transitionTimeMS/100)
	logrus.Infof("Payload is: %v", payload)
	resp, err := tc.Call(tc.dtlsclient.BuildPUTMessage(toDeviceUri(deviceId), payload))
	if err != nil {
		return model.Result{}, err
	}
	logrus.Infof("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

// PutDeviceColorRGB sets the color of the bulb using RGB hex string such as 8f2686 (purple). Note that
// It does not use the built in rgb hex parameter as that does not work reliably, so the rgb is converted to hsl and that is sent
func (tc *Client) PutDeviceColorRGB(deviceId int, rgb string) (model.Result, error) {
	return tc.PutDeviceColorRGBTimed(deviceId, rgb, 500)
}

// PutDeviceColorRGBTimed does the same as PutDeviceColorRGB but it gives you the ability to change the speed at which the color changes
func (tc *Client) PutDeviceColorRGBTimed(deviceId int, rgb string, transitionTimeMS int) (model.Result, error) {
	r, g, b, err := hexStringToRgb(rgb)
	if err != nil {
		return model.Result{}, err
	}

	return tc.PutDeviceColorRGBIntTimed(deviceId, r, g, b, transitionTimeMS)
}

// PutDeviceColorRGBInt does about the same as PutDeviceColorRGB except you can directly pass the rgb instead of a hex string
func (tc *Client) PutDeviceColorRGBInt(deviceId int, r, g, b int) (model.Result, error) {
	return tc.PutDeviceColorRGBIntTimed(deviceId, r, g, b, 500)
}

// PutDeviceColorRGBIntTimed does the same as PutDeviceColorRGBInt but it gives you the ability to change the speed at which the color changes
func (tc *Client) PutDeviceColorRGBIntTimed(deviceId int, r, g, b int, transitionTimeMS int) (model.Result, error) {
	h, s, l := rgbToHsl(r, g, b)

	return tc.PutDeviceColorHSLTimed(deviceId, h, s, l, transitionTimeMS)
}

// PutDeviceColorHSL sets the color of the bulb using the HSL color notation
// This is more effictive than RGB because RGB is always at full brightness, ("000000" is the same as "ffffff")
func (tc *Client) PutDeviceColorHSL(deviceId int, hue float64, saturation float64, lightness float64) (model.Result, error) {
	return tc.PutDeviceColorHSLTimed(deviceId, hue, saturation, lightness, 500)
}

// PutDeviceColorHSLTimed does the same as PutDeviceColorHSL but it gives you the ability to change the speed at which the color changes
func (tc *Client) PutDeviceColorHSLTimed(deviceId int, hue float64, saturation float64, lightness float64, transitionTimeMS int) (model.Result, error) {
	hueInt := int(mapRange(hue, 0, 360, 0, 65279))
	saturationInt := int(mapRange(saturation, 0, 100, 0, 65279))
	lightnessInt := int(mapRange(lightness, 0, 100, 0, 254))

	payload := fmt.Sprintf(`{ "3311": [ {"5707": %d, "5708": %d, "5851": %d, "5712": %d}] }`, hueInt, saturationInt, lightnessInt, transitionTimeMS/100)
	logrus.Infof("Payload is: %v", payload)
	resp, err := tc.Call(tc.dtlsclient.BuildPUTMessage(toDeviceUri(deviceId), payload))
	if err != nil {
		return model.Result{}, err
	}
	logrus.Infof("Response: %+v", resp)
	return model.Result{Msg: resp.Code.String()}, nil
}

// PutDevicePositioning sets the positioning property (0-100) of the specified device.
func (tc *Client) PutDevicePositioning(deviceId int, positioning float32) (model.Result, error) {
	payload := fmt.Sprintf(`{ "15015": [{ "5536": %f }] }`, positioning)
	logrus.Infof("Payload is: %v", payload)
	resp, err := tc.Call(tc.dtlsclient.BuildPUTMessage(toDeviceUri(deviceId), payload))
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

	for _, groupId := range groupIds {
		group, _ := tc.GetGroup(groupId)
		groups = append(groups, group)
	}
	return groups, nil
}

// GetGroup gets the JSON representation of the specified group.
func (tc *Client) GetGroup(groupId int) (model.Group, error) {
	resp, err := tc.Call(tc.dtlsclient.BuildGETMessage(toGroupUri(groupId)))
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
func (tc *Client) GetDevice(deviceId int) (model.Device, error) {
	device := &model.Device{}

	resp, err := tc.Call(tc.dtlsclient.BuildGETMessage(toDeviceUri(deviceId)))
	if err != nil {
		return *device, err
	}

	err = json.Unmarshal(resp.Payload, &device)
	if err != nil {
		return *device, err
	}
	return *device, nil
}

// ListDeviceIds gives you a list of all connected device id's
func (tc *Client) ListDeviceIds() ([]int, error) {
	var devices []int

	resp, err := tc.Call(tc.dtlsclient.BuildGETMessage("/15001/"))
	if err != nil {
		return devices, err
	}

	err = json.Unmarshal(resp.Payload, &devices)
	if err != nil {
		return devices, err
	}
	return devices, nil
}

// ListDevices gives you a list of all devices
func (tc *Client) ListDevices() ([]model.Device, error) {
	var devices []model.Device

	resp, err := tc.ListDeviceIds()
	if err != nil {
		return devices, err
	}

	devices = make([]model.Device, len(resp))

	for i, id := range resp {
		device, err := tc.GetDevice(id)
		if err != nil {
			return devices, err
		}

		devices[i] = device
	}

	return devices, nil
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

func mapRange(x, inMin, inMax, outMin, outMax float64) float64 {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func rgbToHsl(rInt int, gInt int, bInt int) (float64, float64, float64) {
	var r float64 = float64(rInt) / 255
	var g float64 = float64(gInt) / 255
	var b float64 = float64(bInt) / 255

	var maximum float64 = math.Max(r, math.Max(g, b))
	var minimum float64 = math.Min(r, math.Min(g, b))

	var h, s, l float64
	h = (maximum + minimum) / 2
	l = h

	if maximum == minimum {
		h = 0
		s = 0
	} else {
		d := maximum - minimum

		if l > 0.5 {
			s = d / (2 - maximum - minimum)
		} else {
			s = d / (maximum + minimum)
		}

		switch maximum {
		case r:
			if g < b {
				h = (g-b)/d + 6
			} else {
				h = (g-b)/d + 0
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h /= 6
	}

	h *= 360
	s *= 100
	l *= 100

	return h, s, l
}

func hexStringToRgb(hexString string) (int, int, int, error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return 0, 0, 0, err
	}

	return int(bytes[0]), int(bytes[1]), int(bytes[2]), nil
}

func toDeviceUri(deviceId int) string {
	return fmt.Sprintf("/15001/%d", deviceId)
}

func toGroupUri(groupId int) string {
	return fmt.Sprintf("/15004/%d", groupId)
}
