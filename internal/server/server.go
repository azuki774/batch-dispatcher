package server

import (
	"batchdispatcher/internal/model"
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
	Logger     *zap.Logger
	Dispatcher dispatcher
}

type dispatcher interface {
	Run(ctx context.Context, jobname string) (err error)
	GetJobsInfo() (jobInfo []model.JobInfo)
}

func (s *Server) Start(ctx context.Context, port string) (err error) {
	s.Logger.Info("server start")
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
		outputJson, err := json.Marshal(s.Dispatcher.GetJobsInfo())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(outputJson))
	})

	r.Post("/jobs/{jobname}/run", func(w http.ResponseWriter, r *http.Request) {
		jobName := chi.URLParam(r, "jobname")
		fmt.Println(jobName)
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)

	return nil
}
