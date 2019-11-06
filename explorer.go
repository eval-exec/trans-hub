package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ExplorerMetrics struct {
	LambdaExplorerIsAlive                    string `json:"lambda_explorer_is_alive"`
	LambdaExplorerHomepageVisitedCount       string `json:"lambda_explorer_homepage_visited_count"`
	LambdaExplorerLatestHeightFromDb         string `json:"lambda_explorer_latest_height_from_db"`
	LambdaExplorerValidatorsListVisitedCount string `json:"lambda_explorer_validators_list_visited_count"`
}

type CrawlerMetrics struct {
	LambdaCrawlerIsAlive               string `json:"lambda_crawler_is_alive"`
	LambdaCrawlerLatestHeightFromChain string `json:"lambda_crawler_latest_height_from_chain"`
}

var (
	ExplorerProAlive             prometheus.Gauge
	CrawlerProChainHeight        prometheus.Gauge
	ExplorerProHomeVisited       prometheus.Counter
	ExplorerProValidatorsVisited prometheus.Counter
	ExplorerProDBHeight          prometheus.Gauge
	CrawlerProAlive              prometheus.Gauge
)
//
//var ExplorerProAlive = promauto.NewGauge(prometheus.GaugeOpts{
//	Name: lambda_explorer_is_alive,
//	Help: "check lambda_explorer is alive",
//})
//
//var ExplorerProHomeVisited = promauto.NewCounter(prometheus.CounterOpts{
//	Name: lambda_explorer_homepage_visited_count,
//	Help: "homepage visited count, when visited, count++ ",
//})
//
//var ExplorerProValidatorsVisited = promauto.NewCounter(prometheus.CounterOpts{
//	Name: lambda_explorer_validators_list_visited_count,
//	Help: "validator lists visited count",
//})
//
//var ExplorerProDBHeight = promauto.NewGauge(prometheus.GaugeOpts{
//	Name: lambda_explorer_latest_height_from_db,
//	Help: "get latest block height from database",
//})
//
//var CrawlerProChainHeight = promauto.NewGauge(prometheus.GaugeOpts{
//	Name: lambda_crawler_latest_height_from_chain,
//	Help: "get latest block height from block chain",
//})
//
//var CrawlerProAlive = promauto.NewGauge(prometheus.GaugeOpts{
//	Name: lambda_crawler_is_alive,
//	Help: "check lambda_crawler is alive",
//})
