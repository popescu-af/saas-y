package healthy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/popescu-af/saas-y/pkg/log"
)

// NewHealthWatchdog returns a watchdog that wraps and monitors the given service.
func NewHealthWatchdog(name, port, healthPort string, service Service) *HealthWatchdog {
	svc := &HealthWatchdog{
		Handler:    healthcheck.NewHandler(),
		name:       name,
		port:       port,
		healthPort: healthPort,
		service:    service,
	}

	svc.AddLivenessCheck("live-check", service.Live)
	svc.AddReadinessCheck("ready-check", healthcheck.Async(svc.readyLoop, 1*time.Second))
	return svc
}

// Service is the interface for an HTTP service supporting health queries.
type Service interface {
	Initialize() (*http.Server, error)
	Live() error
	Ready() error
}

// HealthWatchdog is a wrapper for healthy services.
// It takes care of the availability reporting and it automatically
// switches between ready and unavailable states based on readiness checks.
type HealthWatchdog struct {
	healthcheck.Handler
	name       string
	port       string
	healthPort string
	service    Service
	server     *http.Server
	ready      bool
}

// Run runs the healthy Service.
func (w *HealthWatchdog) Run() {
	log.InfoCtx("starting service", log.Context{"name": w.name})
	http.ListenAndServe(fmt.Sprintf(":%s", w.healthPort), w.Handler)
}

func (w *HealthWatchdog) readyLoop() error {
	err := w.service.Ready()
	if w.ready && err != nil {
		log.ErrorCtx("ready error", log.Context{"err": err})
		w.switchToUnavailable()
	} else if !w.ready && err == nil {
		w.shutdownServer()
		w.server, err = w.service.Initialize()
		if err == nil {
			w.ready = true
			go func() {
				log.InfoCtx("service ready", log.Context{"name": w.name})
				log.ErrorCtx("error serving", log.Context{"err": w.server.ListenAndServe()})
			}()
		}
	}
	return err
}

func (w *HealthWatchdog) switchToUnavailable() {
	w.shutdownServer()
	w.server = &http.Server{
		Addr: fmt.Sprintf(":%s", w.port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}),
	}

	go func() {
		log.InfoCtx("service unavailable", log.Context{"name": w.name})
		log.ErrorCtx("error serving", log.Context{"err": w.server.ListenAndServe()})
	}()
}

func (w *HealthWatchdog) shutdownServer() {
	if w.server == nil {
		return
	}

	d := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	err := w.server.Shutdown(ctx)
	if err != nil {
		log.WarnCtx("failed to shutdown the server", log.Context{"err": err})
		err = w.server.Close()
		if err != nil {
			log.ErrorCtx("failed to close the server", log.Context{"err": err})
		}
	}

	w.server = nil
	w.ready = false
}
