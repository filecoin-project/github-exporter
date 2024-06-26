package main

import (
	conf "github.com/filecoin-project/github-exporter/config"
	"github.com/filecoin-project/github-exporter/exporter"
	"github.com/filecoin-project/github-exporter/http"
	"github.com/infinityworks/go-common/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	log            *logrus.Logger
	applicationCfg conf.Config
	mets           map[string]*prometheus.Desc
)

func init() {
	applicationCfg = conf.Init()
	mets = exporter.AddMetrics()
	log = logger.Start(&applicationCfg)
}

func main() {
	log.Info("Starting Exporter")

	exp := exporter.Exporter{
		APIMetrics: mets,
		Config:     applicationCfg,
	}

	http.NewServer(exp).Start()
}
