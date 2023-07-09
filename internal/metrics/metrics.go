package metrics

import (
	"batchdispatcher/internal/model"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type MetricsServer struct {
	Logger     *zap.Logger
	Dispatcher Dispatcher
	Port       string
}

type metrics struct {
	lastSuccessDuration prometheus.GaugeVec
}

type Dispatcher interface {
	GetJobsInfo() (jobInfo []model.JobInfo)
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		lastSuccessDuration: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "last_success_duration",
			Help: "shows the between now and last successed running batch",
		},
			[]string{"batch_name"},
		),
	}
	reg.MustRegister(m.lastSuccessDuration)
	return m
}

// refresh: 60 秒ごとに情報を更新
func (s *MetricsServer) refresh(ctx context.Context, m *metrics) {
	go func() {
		for {
			cis := s.Dispatcher.GetJobsInfo()

			for _, ci := range cis {
				pl := prometheus.Labels{"batch_name": ci.Name}

				// 現時刻と、取得した last_success の時間を比較
				nowt := time.Now()
				lastt := ci.LastSuccessStatus
				diff := nowt.Sub(lastt).Seconds()

				m.lastSuccessDuration.With(pl).Set(diff)
			}
			time.Sleep(60 * time.Second)
		}
	}()
	<-ctx.Done()
	s.Logger.Info("refresh routine close")
}

func (s *MetricsServer) Start(ctx context.Context) error {
	s.Logger.Info("cert-metrics server start", zap.String("port", s.Port))
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := NewMetrics(reg)

	go s.refresh(ctx, m)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Port),
		Handler: nil,
	}
	go func() {
		<-ctx.Done()
		s.Logger.Info("shutdown signal catch")
		ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		nerr := server.Shutdown(ctx2)
		if nerr != nil {
			s.Logger.Error("gracefully shutdown error", zap.Error(nerr))
		}
	}()

	s.Logger.Info("start listening", zap.String("port", s.Port))
	err := server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		// expected error
		err = nil
	} else {
		s.Logger.Error("metrics server close error", zap.Error(err))
		return err
	}

	s.Logger.Info("metrics server shutdown")
	return nil
}
