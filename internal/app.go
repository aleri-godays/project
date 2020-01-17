package internal

import (
	"context"
	"github.com/aleri-godays/project"
	"github.com/aleri-godays/project/internal/config"
	"github.com/aleri-godays/project/internal/db"
	"github.com/aleri-godays/project/internal/http"
	"github.com/asdine/storm/v3"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type App struct {
	server      *http.Server
	conf        *config.Config
	stormDB     *storm.DB
	traceCloser io.Closer
}

func NewApp(conf *config.Config) *App {

	configureLogger(conf.LogLevel)

	stormDB := db.NewStormDB(conf.DbPath)
	a := &App{
		conf:        conf,
		traceCloser: initTracer(conf.ServiceName),
		stormDB:     stormDB,
		server:      http.NewServer(conf, db.NewStormRepository(stormDB)),
	}

	return a
}

func (a *App) Run() {
	seedDB(a.stormDB)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	//start webserver
	go func() {
		log.WithFields(log.Fields{
			"port":      a.conf.HTTPPort,
			"log_level": a.conf.LogLevel,
			"version":   a.conf.Version,
		}).Info("starting server")
		if err := a.server.Start(); err != nil {
			log.WithFields(log.Fields{
				"reason": err,
			}).Info("shutting down the server")
			quit <- os.Interrupt
		}
	}()

	<-quit
	log.Info("received a stop signal")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	//close http connections
	if err := a.server.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("errors during http server shutdown")
	} else {
		log.Info("http server closed")
	}

	a.stormDB.Close()
	a.traceCloser.Close()
}

func initTracer(serviceName string) io.Closer {
	defcfg := jaegerconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	cfg, err := defcfg.FromEnv()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not parse jaeger env vars")
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not initialize jaeger tracer")
	}

	opentracing.SetGlobalTracer(tracer)
	return closer
}

func configureLogger(logLevel string) {
	switch strings.ToLower(logLevel) {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}
}

func seedDB(sdb *storm.DB) {
	repo := db.NewStormRepository(sdb)
	if all, err := repo.All(context.Background()); err != nil {
		log.Errorf("could not seed stormDB: %s", err)
	} else if len(all) > 0 {
		return
	}
	log.Info("seeding database")

	repo.Add(context.Background(), &project.Project{
		Name:      "Web development",
		User:      "Gopher",
		DailyRate: 650,
	})
	repo.Add(context.Background(), &project.Project{
		Name:      "Microservices workshop",
		User:      "Martin Radile",
		DailyRate: 600,
	})
	repo.Add(context.Background(), &project.Project{
		Name:      "REST interface design",
		User:      "Christian Sage",
		DailyRate: 625,
	})
}
