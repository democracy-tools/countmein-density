package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/democracy-tools/countmein-density/internal"
	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/gorilla/mux"
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

func main() {

	const (
		observations     = "/observations"
		users            = "/users/{user-id}"
		volunteers       = "/volunteers/{user-id}"
		volunteerPolygon = "/volunteers/{user-id}/polygons/{polygon}"
		polygons         = "/polygons"
		whatsappWebhook  = "/whatsapp"
	)

	handle := internal.NewHandle(
		ds.NewClientWrapper(env.Project),
		whatsapp.NewClientWrapper(),
	)
	serve(
		[]string{
			observations, observations, observations,
			users, users,
			volunteers, volunteers,
			volunteerPolygon, volunteerPolygon,
			polygons, polygons,
			whatsappWebhook, whatsappWebhook,
		}, []string{
			http.MethodPost, http.MethodGet, http.MethodOptions,
			http.MethodDelete, http.MethodOptions,
			http.MethodGet, http.MethodOptions,
			http.MethodPut, http.MethodOptions,
			http.MethodGet, http.MethodOptions,
			http.MethodGet, http.MethodPost,
		}, []func(http.ResponseWriter, *http.Request){
			access(handle.CreateObservation), access(handle.GetObservations), options([]string{http.MethodPost, http.MethodGet}),
			access(handle.DeleteUser), options([]string{http.MethodDelete}),
			access(handle.GetVolunteer), options([]string{http.MethodGet}),
			access(handle.ChangePolygon), options([]string{http.MethodPut}),
			access(handle.GetAvailablePolygons), options([]string{http.MethodGet}),
			handle.WhatsAppVerification, handle.WhatsAppEventHandler,
		},
	)
}

func access(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next(w, r)
	}
}

func options(methods []string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
	}
}

func serve(path []string, method []string,
	handle []func(http.ResponseWriter, *http.Request), mwf ...mux.MiddlewareFunc) {

	router := mux.NewRouter()
	router.Use(mwf...)
	for i := 0; i < len(path); i++ {
		router.HandleFunc(path[i], handle[i]).Methods(method[i])
	}

	serveMulti([]*mux.Router{router}, []string{"8080"})
}

func serveMulti(routers []*mux.Router, ports []string) {

	initLogger()
	logVersion()

	var servers []*http.Server
	for i := 0; i < len(ports); i++ {
		servers = append(servers, &http.Server{
			Addr: fmt.Sprintf("%s:%s", "0.0.0.0", ports[i]),
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      routers[i],
		})
		go func(server *http.Server, port string) {
			logrus.Debugf("listening on port %q", port)
			if err := server.ListenAndServe(); err != nil {
				logrus.Error(err)
			}
		}(servers[i], ports[i])
	}
	c := make(chan os.Signal, 1)
	// Graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	<-c

	for _, currServer := range servers {
		shutdown(currServer)
	}

	logrus.Info("exit")
	os.Exit(0)
}

func shutdown(server *http.Server) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("failed to shutdown server with %q", err)
	}
}

func logVersion() {

	logrus.Infof("%s/%s, %s", runtime.GOOS, runtime.GOARCH, runtime.Version())
}

func initLogger() {

	// logrus.SetReportCaller(true)
	initLoggerOutput()
	logrus.SetLevel(getLogLevel())
}

func initLoggerOutput() {

	logrus.SetOutput(io.Discard) // Send all logs to nowhere by default - this is required to avoid duplicate log messages
	logrus.AddHook(filename.NewHook())
	logrus.AddHook(&writer.Hook{ // Send logs with level higher than warning to stderr
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		},
	})
	logrus.AddHook(&writer.Hook{ // Send info and debug logs to stdout
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
			logrus.TraceLevel,
		},
	})
}

func getLogLevel() logrus.Level {

	level := os.Getenv("LOG_LEVEL")
	if strings.EqualFold(level, "debug") {
		return logrus.DebugLevel
	} else if strings.EqualFold(level, "warn") {
		return logrus.WarnLevel
	} else if strings.EqualFold(level, "error") {
		return logrus.ErrorLevel
	}
	return logrus.InfoLevel
}
