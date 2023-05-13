package server

import (
	"batchdispatcher/internal/dispatcher"
	"batchdispatcher/internal/job"
	"batchdispatcher/internal/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type Server struct {
	l          *zap.Logger
	dispatcher *dispatcher.Dispatcher
}

var sampleJobs = []job.Job{
	job.NewJob("ls", "ls -la"),
	job.NewJob("sleep", "sleep 15s"),
}

func NewServer() (srv *Server, err error) {
	srv = &Server{}

	// set logger
	l, err := logger.NewLogger()
	if err != nil {
		return &Server{}, err
	}
	srv.l = l

	// set dispatcher
	l.Info("set dispatcher")
	srv.dispatcher = &dispatcher.Dispatcher{
		Logger: l,
		Jobs:   sampleJobs, // TODO
	}

	return srv, nil
}

func (s *Server) Start(ctx context.Context, port string) (err error) {
	s.l.Info("server start")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Get("/jobs", func(w http.ResponseWriter, r *http.Request) {
		outputJson, err := json.Marshal(&s.dispatcher.Jobs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(outputJson))
	})

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)

	return nil
}
