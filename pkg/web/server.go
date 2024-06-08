// This package includes web api
package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didip/tollbooth/v7"
	gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/williampsena/bugs-channel/pkg/config"
	"github.com/williampsena/bugs-channel/pkg/storage"
)

// Represents the HTTP handlers and router instances.
type Server struct {
	Router  http.Handler
	Srv     *http.Server
	log     *logrus.Logger
	Context *ServerContext
}

// The web server context
type ServerContext struct {
	context.Context
	Queue storage.Queue
}

// Creates and returns a new instance of Server
func NewServer(c *ServerContext, handler http.Handler, log *logrus.Logger) *Server {
	ch := gorilla.CORS(gorilla.AllowedOrigins([]string{"*"}))

	return &Server{
		Context: c,
		Router:  handler,
		log:     log,
		Srv: &http.Server{
			Addr:         fmt.Sprintf(":%v", config.ApiPort()),
			Handler:      ch(handler),
			IdleTimeout:  time.Second * 5,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		},
	}
}

// Waits for an interrupt signal before gracefully shutting down the handlers.
func (s *Server) GraceFulShutDown(killTime time.Duration) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	ctx, cancel := context.WithTimeout(s.Context, killTime)

	defer cancel()

	s.log.Print("🛑 The bugs channel server has been shut down.")

	if err := s.Srv.Shutdown(ctx); err != nil {
		s.log.Fatalf("❌ Unexpected interruption to the bugs channel server's listening: %s\n", err)
	}

	s.log.Print("❎ The bugs channel server exited properly")

}

// Turn on the HTTP handlers and listen in.
func (s *Server) ListenAndServe() error {
	return s.Srv.ListenAndServe()
}

// Shutdown terminates the handlers for HTTP.
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Print("🛑 The bugs channel server was shut down.")
	return s.Srv.Shutdown(ctx)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return gorilla.LoggingHandler(os.Stdout, next)
}

func maybeUseRatelimitHandler(r *mux.Router) {
	rateLimit := config.RateLimit()

	if rateLimit == 0 {
		logrus.Warnf("💡 RateLimit middleware is disabled %v", config.RateLimit())
	}

	limiter := BuildRateLimitMiddleware(rateLimit)

	handler := func(next http.Handler) http.Handler {
		return tollbooth.LimitFuncHandler(limiter, func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	r.Use(handler)
}

func buildRouter() (*mux.Router, error) {
	r := mux.NewRouter()

	r.PathPrefix("/health").HandlerFunc(HealthCheckEndpoint).Methods("GET")

	r.PathPrefix("/").HandlerFunc(NoRouteEndpoint)

	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(loggingMiddleware)
	maybeUseRatelimitHandler(r)

	return r, nil
}

// Setup the bugs channel web server
func SetupServer(context *ServerContext) (*Server, error) {
	r, err := buildRouter()

	if err != nil {
		return nil, err
	}

	srv := NewServer(context, r, logrus.StandardLogger())

	go srv.ListenAndServe()

	srv.log.Infof("🐛 Bugs Channel Sever listening at port %v...", config.ApiPort())

	srv.GraceFulShutDown(time.Second * 5)

	return srv, nil
}
