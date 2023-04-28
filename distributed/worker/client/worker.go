package client

import (
	"go_crawler/distributed/config"
	"go_crawler/distributed/worker"
	"go_crawler/engine"
	"net/rpc"
)

func CreateProcessor(clientChan chan *rpc.Client) engine.Processor {

	return func(req engine.Request) (engine.ParseResult, error) {

		sReq := worker.SerializeRequest(req)

		var sResult worker.ParseResult

		cli := <-clientChan
		err := cli.Call(config.CrawlServiceRpc, sReq, &sResult)
		if err != nil {
			return engine.ParseResult{}, err
		}

		return worker.DeserializeResult(sResult), nil
	}
}
