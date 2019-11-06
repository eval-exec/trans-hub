package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	explorerUrl string
	crawlerUrl  string
	PushHome    string
	role        string
)

var (
	lambda_crawler_is_alive                       string
	lambda_crawler_latest_height_from_chain       string
	lambda_explorer_latest_height_from_db         string
	lambda_explorer_homepage_visited_count        string
	lambda_explorer_is_alive                      string
	lambda_explorer_validators_list_visited_count string
)

var (
	lambda_crawler_job  string
	lambda_explorer_job string
)

func init() {

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	logs.SetLogger(logs.AdapterFile, `{"filename":"trans_hub.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)

	explorerUrl = viper.GetString("explorer_url")
	crawlerUrl = viper.GetString("crawler_url")
	PushHome = viper.GetString("push_home")
	role = viper.GetString("role")

}

func initVal() {
	lambda_crawler_is_alive = role + "lambda_crawler_is_alive"
	lambda_crawler_latest_height_from_chain = role + "lambda_crawler_latest_height_from_chain"
	lambda_explorer_latest_height_from_db = role + "lambda_explorer_latest_height_from_db"
	lambda_explorer_homepage_visited_count = role + "lambda_explorer_homepage_visited_count"
	lambda_explorer_is_alive = role + "lambda_explorer_is_alive"
	lambda_explorer_validators_list_visited_count = role + "lambda_explorer_validators_list_visited_count"

	lambda_crawler_job = role + "lambda_crawler_job"
	lambda_explorer_job = role + "lambda_explorer_job"

}

func initPro() {

	ExplorerProAlive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: lambda_explorer_is_alive,
		Help: "check lambda_explorer is alive",
	})

	ExplorerProHomeVisited = promauto.NewCounter(prometheus.CounterOpts{
		Name: lambda_explorer_homepage_visited_count,
		Help: "homepage visited count, when visited, count++ ",
	})

	ExplorerProValidatorsVisited = promauto.NewCounter(prometheus.CounterOpts{
		Name: lambda_explorer_validators_list_visited_count,
		Help: "validator lists visited count",
	})

	ExplorerProDBHeight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: lambda_explorer_latest_height_from_db,
		Help: "get latest block height from database",
	})

	CrawlerProChainHeight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: lambda_crawler_latest_height_from_chain,
		Help: "get latest block height from block chain",
	})

	CrawlerProAlive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: lambda_crawler_is_alive,
		Help: "check lambda_crawler is alive",
	})
}

var explorerMetrics ExplorerMetrics
var crawlerMetics CrawlerMetrics

func getExplorerMetrisc() {
	resp, err := http.Get(explorerUrl)
	if err != nil {
		logs.Error("get explore metrics error")
		return
	}

	defer func(){
		err := resp.Body.Close()
		if err != nil {
			logs.Error("fetch explorer metrics erorr %s",err.Error())
		}
	}()

	exbyte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("read metrics error")
		return

	}
	exfetch := string(exbyte)
	metrics := strings.Split(exfetch, "\n")
	for _, x := range metrics {
		if strings.HasPrefix(x, "#") {
			continue
		}
		if strings.Contains(x, lambda_explorer_latest_height_from_db) {
			y := strings.Split(x, " ")
			explorerMetrics.LambdaExplorerLatestHeightFromDb = y[1]
		}
		if strings.Contains(x, lambda_explorer_homepage_visited_count) {
			y := strings.Split(x, " ")
			explorerMetrics.LambdaExplorerHomepageVisitedCount = y[1]
		}
		if strings.Contains(x, lambda_explorer_is_alive) {
			y := strings.Split(x, " ")
			explorerMetrics.LambdaExplorerIsAlive = y[1]
		}
		if strings.Contains(x, lambda_explorer_validators_list_visited_count) {
			y := strings.Split(x, " ")
			explorerMetrics.LambdaExplorerValidatorsListVisitedCount = y[1]
		}
	}
	logs.Debug("explorer get metrics :%s", explorerMetrics)

	ExplorerProAlive.Set(getFloat64(explorerMetrics.LambdaExplorerIsAlive))
	ExplorerProHomeVisited.Add(getFloat64(explorerMetrics.LambdaExplorerIsAlive))
	ExplorerProDBHeight.Set(getFloat64(explorerMetrics.LambdaExplorerLatestHeightFromDb))
	ExplorerProValidatorsVisited.Add(getFloat64(explorerMetrics.LambdaExplorerIsAlive))

}

func getFloat64(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

func getCrawerMetrisc() {
	resp, err := http.Get(crawlerUrl)
	if err != nil {
		logs.Error("get crawler metrics error")
		return
	}


	defer func(){
		err := resp.Body.Close()
		if err != nil {
			logs.Error("fetch explorer metrics erorr %s",err.Error())
		}
	}()
	exbyte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("read metrics error")
	}
	exfetch := string(exbyte)
	metrics := strings.Split(exfetch, "\n")
	for _, x := range metrics {
		if strings.HasPrefix(x, "#") {
			continue
		}
		if strings.Contains(x, lambda_crawler_is_alive) {
			y := strings.Split(x, " ")
			crawlerMetics.LambdaCrawlerIsAlive = y[1]
		}
		if strings.Contains(x, lambda_crawler_latest_height_from_chain) {
			y := strings.Split(x, " ")
			crawlerMetics.LambdaCrawlerLatestHeightFromChain = y[1]
		}

	}
	logs.Debug("crawler get metrics :%s", crawlerMetics)

	CrawlerProAlive.Set(getFloat64(crawlerMetics.LambdaCrawlerIsAlive))
	CrawlerProChainHeight.Set(getFloat64(crawlerMetics.LambdaCrawlerLatestHeightFromChain))
}

func pushAllToGateway() {
	// explorer
	if err := push.New(PushHome, lambda_explorer_job).Collector(ExplorerProAlive).Grouping("alive", "is_alive").Push(); err != nil {
		logs.Error("push alive error", err.Error())
	}
	if err := push.New(PushHome, lambda_explorer_job).Collector(ExplorerProHomeVisited).Grouping("home", "home_visited_count").Push(); err != nil {
		logs.Error("push homepage visited error", err.Error())
	}
	if err := push.New(PushHome, lambda_explorer_job).Collector(ExplorerProValidatorsVisited).Grouping("validator", "validator_visited_count").Push(); err != nil {
		logs.Error("push validators visited error", err.Error())
	}

	if err := push.New(PushHome, lambda_explorer_job).Collector(ExplorerProDBHeight).Grouping("database", "database_block_height").Push(); err != nil {
		logs.Error("push database block height error", err.Error())
	}

	// crawler
	if err := push.New(PushHome, lambda_crawler_job).Collector(CrawlerProAlive).Grouping("alive", "is_alive").Push(); err != nil {
		logs.Error("prometheus push alive error", err.Error())
	}

	if err := push.New(PushHome, lambda_crawler_job).Collector(CrawlerProChainHeight).Grouping("database", "chain_latest_height").Push(); err != nil {
		logs.Error("prometheus push database block height error", err.Error())
	}
}
func main() {
	initVal()
	initPro()
	logs.Info("explorerUrl is %s", explorerUrl)
	logs.Info("crawler urlis %s", crawlerUrl)
	logs.Info("pushhomwis %s", PushHome)
	for {
		getExplorerMetrisc()
		logs.Info("explorer fetch finish")
		getCrawerMetrisc()
		logs.Info("crawler fetch finish")
		pushAllToGateway()
		logs.Info("push all finish ")
		time.Sleep(20 * time.Second)
	}
}
