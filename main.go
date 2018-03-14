package main

//
//                       _oo0oo_
//                      o8888888o
//                      88" . "88
//                      (| -_- |)
//                      0\  =  /0
//                    ___/`---'\___
//                  .' \\|     |// '.
//                 / \\|||  :  |||// \
//                / _||||| -:- |||||- \
//               |   | \\\  -  /// |   |
//               | \_|  ''\---/''  |_/ |
//               \  .-\__  '-'  ___/-. /
//             ___'. .'  /--.--\  `. .'___
//          ."" '<  `.___\_<|>_/___.' >' "".
//         | | :  `- \`.;`\ _ /`;.`/ - ` : | |
//         \  \ `_.   \_ __\ /__ _/   .-` /  /
//     =====`-.____`.___ \_____/___.-`___.-'=====
//                       `=---='
//
//
//     ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//
//               佛祖保佑         永无BUG
//

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"runtime/debug"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
)

func parseCommandLineParams() {
	flag.StringVar(&configPath, "c", "./config.yml", "Path to config.yml")
	flag.Parse()
}

func initConfigs() {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln("error reading config", err)
	}

	NodesList = make(map[string]nodeElement)
	for _, element := range config.Nodes {
		NodesList[element.NodeName] = nodeElement{NodeName: element.NodeName, NodeURL: element.NodeURL, AESKey: []byte(element.AESKey)}
	}
}

func initRuntime() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	Logger.Printf("Init runtime to use %d CPUs and %d threads", numCPU, config.System.MaxThreads)
	debug.SetMaxThreads(config.System.MaxThreads)
}

func main() {
	fmt.Printf("Version:    [%s]\nBuild:      [%s]\nBuild Date: [%s]\n", version, build, buildDate)
	parseCommandLineParams()
	initConfigs()
	initLogger()
	initRuntime()

	if config.Clickhouse.IsEnabled {
		connectClickDB()
	}

	router := routing.New()

	router.Post("/<node>/*", StartHandler, ProxyHandler, FinishHandler)
	router.Post("*", OptionsHandler)
	router.Get("*", OptionsHandler)
	router.Options("*", OptionsHandler)

	server := fasthttp.Server{
		Name:    ServerName,
		Handler: router.HandleRequest,
		Logger:  Logger,
	}

	Logger.Println("Started fasthttp server on port", config.System.ListenOn)
	Logger.Fatal(server.ListenAndServe(config.System.ListenOn))
}
