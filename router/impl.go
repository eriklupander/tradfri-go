package router

import (
	"encoding/json"
	"github.com/eriklupander/tradfri-go/model"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strconv"
)

func setColorXY(w http.ResponseWriter, r *http.Request) {
	deviceId, err := paramToInt(chi.URLParam(r, deviceParam))
	if err != nil {
		badIdentifierError(w, chi.URLParam(r, deviceParam), err)
		return
	}
	xStr := chi.URLParam(r, "x")
	yStr := chi.URLParam(r, "y")
	x, _ := strconv.Atoi(xStr)
	y, _ := strconv.Atoi(yStr)
	res, err := tradfriClient.PutDeviceColor(deviceId, x, y)
	respond(w, res, err)
}

func setColorRGBHex(w http.ResponseWriter, r *http.Request) {
	deviceId, err := paramToInt(chi.URLParam(r, deviceParam))
	if err != nil {
		badIdentifierError(w, chi.URLParam(r, deviceParam), err)
		return
	}

	body, _ := io.ReadAll(r.Body)

	rgbColorRequest := model.RgbColorRequest{}
	if err := json.Unmarshal(body, &rgbColorRequest); err != nil {
		badRequest(w, err)
		return
	}
	result, err := tradfriClient.PutDeviceColorRGB(deviceId, rgbColorRequest.RGBcolor)
	respond(w, result, err)
}

func setDimming(w http.ResponseWriter, r *http.Request) {
	deviceId, err := paramToInt(chi.URLParam(r, deviceParam))
	if err != nil {
		badIdentifierError(w, chi.URLParam(r, deviceParam), err)
		return
	}

	body, _ := io.ReadAll(r.Body)

	dimmingRequest := model.DimmingRequest{}
	if err := json.Unmarshal(body, &dimmingRequest); err != nil {
		badRequest(w, err)
		return
	}
	res, err := tradfriClient.PutDeviceDimming(deviceId, dimmingRequest.Dimming)
	respond(w, res, err)
}

func setPower(w http.ResponseWriter, r *http.Request) {
	deviceId, err := paramToInt(chi.URLParam(r, deviceParam))
	if err != nil {
		badIdentifierError(w, chi.URLParam(r, deviceParam), err)
		return
	}
	body, _ := io.ReadAll(r.Body)

	powerRequest := model.PowerRequest{}
	if err := json.Unmarshal(body, &powerRequest); err != nil {
		badRequest(w, err)
		return
	}
	res, err := tradfriClient.PutDevicePower(deviceId, powerRequest.Power)
	respond(w, res, err)
}

func setState(w http.ResponseWriter, r *http.Request) {
	deviceId, err := paramToInt(chi.URLParam(r, deviceParam))
	if err != nil {
		badIdentifierError(w, chi.URLParam(r, deviceParam), err)
		return
	}
	body, _ := io.ReadAll(r.Body)

	stateReq := model.StateRequest{}
	if err := json.Unmarshal(body, &stateReq); err != nil {
		badRequest(w, err)
		return
	}
	res, err := tradfriClient.PutDeviceState(deviceId, stateReq.Power, stateReq.Dimmer)
	respond(w, res, err)
}

func setPositioning(w http.ResponseWriter, r *http.Request) {
	deviceId, err := paramToInt(chi.URLParam(r, deviceParam))
	if err != nil {
		badIdentifierError(w, chi.URLParam(r, deviceParam), err)
		return
	}

	body, _ := io.ReadAll(r.Body)
	positioningReq := model.PositioningRequest{}
	if err := json.Unmarshal(body, &positioningReq); err != nil {
		badRequest(w, errors.Wrap(err, "unmarshalling of positioning JSON body failed"))
		return
	}

	res, err := tradfriClient.PutDevicePositioning(deviceId, positioningReq.Positioning)
	respond(w, res, err)
}

func listGroups(w http.ResponseWriter, _ *http.Request) {
	groups, err := tradfriClient.ListGroups()
	groupResponses := make([]model.GroupResponse, 0)
	for _, g := range groups {
		groupResponses = append(groupResponses, model.ToGroupResponse(g))
	}
	respond(w, groupResponses, err)
}

func getGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := paramToInt(chi.URLParam(r, groupParam))
	if err != nil {
		badIdentifierError(w, chi.URLParam(r, groupParam), err)
		return
	}

	group, err := tradfriClient.GetGroup(groupId)
	respond(w, model.ToGroupResponse(group), err)
}

func getDevicesOnGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := paramToInt(chi.URLParam(r, groupParam))
	if err != nil {
		badIdentifierError(w, chi.URLParam(r, groupParam), err)
		return
	}

	group, _ := tradfriClient.GetGroup(groupId)
	devices := make([]interface{}, 0)
	for _, deviceID := range group.Content.DeviceList.DeviceIds {
		device, _ := tradfriClient.GetDevice(deviceID)
		devices = append(devices, model.ToDeviceResponse(device))
	}
	respondWithJSON(w, 200, devices)
}

func getDeviceIdsOnGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := paramToInt(chi.URLParam(r, groupParam))
	if err != nil {
		badIdentifierError(w, chi.URLParam(r, groupParam), err)
		return
	}

	group, _ := tradfriClient.GetGroup(groupId)
	deviceIds := make([]int, 0)
	deviceIds = append(deviceIds, group.Content.DeviceList.DeviceIds...)
	respondWithJSON(w, 200, deviceIds)
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	deviceId, err := paramToInt(chi.URLParam(r, deviceParam))
	if err != nil {
		badIdentifierError(w, chi.URLParam(r, deviceParam), err)
		return
	}
	device, _ := tradfriClient.GetDevice(deviceId)
	respondWithJSON(w, 200, model.ToDeviceResponse(device))
}
