package main

import (
	"flag"
	"fmt"
	"go_crawler/distributed/rpcsupport"
	"go_crawler/distributed/worker"
	"log"
)

//命令行参数
var port = flag.Int("port", 0, "the port for me to listen on")

//var port = fmt.Sprintf(":%d", config.WorkerPort0)

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify port")
		return
	}
	log.Fatal(rpcsupport.ServeRpc(
		fmt.Sprintf(":%d", *port),
		worker.CrawlService{}))
}
