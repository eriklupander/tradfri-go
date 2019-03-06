package router

import (
	"encoding/json"
	"fmt"
	"github.com/eriklupander/tradfri-go/model"
	"github.com/eriklupander/tradfri-go/tradfri"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var dtlsClient *tradfri.DtlsClient

func SetupChi(client *tradfri.DtlsClient) {
	dtlsClient = client

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// RESTy routes for "api" resource
	r.Route("/api", func(r chi.Router) {
		r.Get("/groups", listGroups)
		r.Get("/groups/{groupId}", getGroup)
		r.Get("/groups/{groupId}/deviceIds", getDeviceIdsOnGroup)
		r.Get("/groups/{groupId}/devices", getDevicesOnGroup)
		r.Get("/device/{deviceId}", getDevice)
		r.Put("/device/{deviceId}/color/{x}/{y}", setColorXY)
		r.Put("/device/{deviceId}/rgb", setColorRGBHex)
		r.Put("/device/{deviceId}/dimmer/{dimming}", setDimming)

	})
	http.ListenAndServe(":8080", r)

}

func setColorXY(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "deviceId")
	xStr := chi.URLParam(r, "x")
	yStr := chi.URLParam(r, "y")
	x, _ := strconv.Atoi(xStr)
	y, _ := strconv.Atoi(yStr)
	err := dtlsClient.PutDeviceColor(deviceId, x, y)
	respond(w, "{}", err)
}

type RgbColorRequest struct {
	RGBcolor string `json:"rgbcolor"`
}

func setColorRGBHex(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "deviceId")
	body, _ := ioutil.ReadAll(r.Body)

	rgbColorRequest := RgbColorRequest{}
	err := json.Unmarshal(body, &rgbColorRequest)
	err = dtlsClient.PutDeviceColorRGB(deviceId, rgbColorRequest.RGBcolor)
	respond(w, "{}", err)
}

func setDimming(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "deviceId")
	dimmingStr := chi.URLParam(r, "dimming")
	dimming, _ := strconv.Atoi(dimmingStr)

	err := dtlsClient.PutDeviceDimming(deviceId, dimming)
	respond(w, "{}", err)
}

func listGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := dtlsClient.ListGroups()
	respond(w, groups, err)
}

func getGroup(w http.ResponseWriter, r *http.Request) {
	group, err := dtlsClient.GetGroup(chi.URLParam(r, "groupId"))
	respond(w, group, err)
}

func getDevicesOnGroup(w http.ResponseWriter, r *http.Request) {
	group, _ := dtlsClient.GetGroup(chi.URLParam(r, "groupId"))
	devices := make([]model.BulbResponse, 0)
	for _, deviceID := range group.Content.DeviceList.DeviceIds {
		device, _ := dtlsClient.GetDevice(strconv.Itoa(deviceID))
		devices = append(devices, model.ToDeviceResponse(device))
	}
	respondWithJSON(w, 200, devices)
}

func getDeviceIdsOnGroup(w http.ResponseWriter, r *http.Request) {
	group, _ := dtlsClient.GetGroup(chi.URLParam(r, "groupId"))
	deviceIds := make([]int, 0)
	for _, deviceID := range group.Content.DeviceList.DeviceIds {
		deviceIds = append(deviceIds, deviceID)
	}
	respondWithJSON(w, 200, deviceIds)
}

func respond(w http.ResponseWriter, payload interface{}, err error) {
	if err != nil {
		respondWithError(w, 500, err.Error())
	} else {
		respondWithJSON(w, 200, payload)
	}
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	device, _ := dtlsClient.GetDevice(chi.URLParam(r, "deviceId"))
	respondWithJSON(w, 200, model.ToDeviceResponse(device))
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"message": msg})
}

// respondWithJSON write json response format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	fmt.Println(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
