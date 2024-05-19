package router

import (
	"net/http"
	"time"

	"github.com/eriklupander/tradfri-go/tradfri"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

const (
	deviceParam = "deviceId"
	groupParam  = "groupId"
)

var tradfriClient *tradfri.Client

// SetupChi sets up our HTTP router/muxer using Chi, a pointer to a Client must be passed.
func SetupChi(client *tradfri.Client, listenAddress string) {
	tradfriClient = client
	logger := middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logrus.StandardLogger(), NoColor: false})
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	// RESTy routes for "api" resource
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

	// Blocks here!
	if err := http.ListenAndServe(listenAddress, r); err != nil {
		logrus.WithError(err).Fatal("error starting HTTP server")
	}
}
