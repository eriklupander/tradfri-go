package router

import (
	"encoding/json"
	"fmt"
	"github.com/eriklupander/tradfri-go/model"
	"github.com/eriklupander/tradfri-go/tradfri"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var tradfriClient *tradfri.Client

// SetupChi sets up our HTTP router/muxer using Chi, a pointer to a Client must be passed.
func SetupChi(client *tradfri.Client, port int) {
	tradfriClient = client
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

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
		r.Put("/device/{deviceId}/position", setPositioning)

	})

	// Blocks here!
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func setPositioning(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "deviceId")
	body, _ := ioutil.ReadAll(r.Body)

	positioningRequest := model.PositioningRequest{}
	if err := json.Unmarshal(body, &positioningRequest); err != nil {
		badRequest(w, err)
		return
	}
	res, err := tradfriClient.PutDevicePositioning(deviceId, positioningRequest.Positioning)
	respond(w, res, err)
}

func listGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := tradfriClient.ListGroups()
	groupResponses := make([]model.GroupResponse, 0)
	for _, g := range groups {
		groupResponses = append(groupResponses, model.ToGroupResponse(g))
	}
	respond(w, groupResponses, err)
}

func getGroup(w http.ResponseWriter, r *http.Request) {
	group, err := tradfriClient.GetGroup(chi.URLParam(r, "groupId"))
	respond(w, model.ToGroupResponse(group), err)
}

func getDevicesOnGroup(w http.ResponseWriter, r *http.Request) {
	group, _ := tradfriClient.GetGroup(chi.URLParam(r, "groupId"))
	devices := make([]model.BlindResponse, 0)
	for _, deviceID := range group.Content.DeviceList.DeviceIds {
		device, _ := tradfriClient.GetDevice(strconv.Itoa(deviceID))
		devices = append(devices, model.ToDeviceResponse(device))
	}
	respondWithJSON(w, 200, devices)
}

func getDeviceIdsOnGroup(w http.ResponseWriter, r *http.Request) {
	group, _ := tradfriClient.GetGroup(chi.URLParam(r, "groupId"))
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

func badRequest(w http.ResponseWriter, err error) {
	logrus.WithError(err).Error("error processing request body")
	respondWithError(w, http.StatusBadRequest, err.Error())
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	device, _ := tradfriClient.GetDevice(chi.URLParam(r, "deviceId"))
	respondWithJSON(w, 200, model.ToDeviceResponse(device))
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"message": msg})
}

// respondWithJSON write json response format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	logrus.Info(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
