package router

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/eriklupander/tradfri-go/model"
	"github.com/eriklupander/tradfri-go/tradfri"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	deviceParam = "deviceId"
	groupParam  = "groupId"
)

// TradfriClient defines the gateway operations used by the HTTP handlers.
type TradfriClient interface {
	GetDevice(deviceId int) (model.Device, error)
	GetGroup(groupId int) (model.Group, error)
	ListGroups() ([]model.Group, error)
	PutDeviceColor(deviceId int, x, y int) (model.Result, error)
	PutDeviceColorRGB(deviceId int, rgb string) (model.Result, error)
	PutDeviceDimming(deviceId int, dimming int) (model.Result, error)
	PutDevicePower(deviceId int, power int) (model.Result, error)
	PutDeviceState(deviceId int, power int, dimmer int) (model.Result, error)
	PutDevicePositioning(deviceId int, positioning float32) (model.Result, error)
}

var tradfriClient TradfriClient

// newRouter builds and returns the chi router wired to the provided client.
func newRouter(client TradfriClient) chi.Router {
	tradfriClient = client
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/groups", listGroups)
		r.Get("/groups/{groupId}", getGroup)
		r.Get("/groups/{groupId}/deviceIds", getDeviceIdsOnGroup)
		r.Get("/groups/{groupId}/devices", getDevicesOnGroup)
		r.Get("/device/{deviceId}", getDevice)
		r.Put("/device/{deviceId}/color", setColorXY)
		r.Put("/device/{deviceId}/rgb", setColorRGBHex)
		r.Put("/device/{deviceId}/dimmer", setDimming)
		r.Put("/device/{deviceId}/power", setPower)
		r.Put("/device/{deviceId}", setState)
		r.Put("/device/{deviceId}/position", setPositioning)
	})
	return r
}

// SetupChi sets up our HTTP router/muxer using Chi, a pointer to a Client must be passed.
func SetupChi(client *tradfri.Client, listenAddress string) {
	r := newRouter(client)
	// Blocks here!
	if err := http.ListenAndServe(listenAddress, r); err != nil {
		slog.Error("error starting HTTP server", slog.Any("error", err))
		os.Exit(1)
	}
}
