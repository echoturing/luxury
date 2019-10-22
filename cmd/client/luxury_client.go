package main

import (
	"context"
	"flag"
	"runtime"
	"sync"
	"time"

	"github.com/echoturing/luxury/conf"
	"github.com/echoturing/luxury/crawlers/hermes"
	"github.com/echoturing/luxury/default_web_navigator"
	"github.com/echoturing/luxury/logger"
	"github.com/echoturing/luxury/pprof"
)

var (
	interval int64
	query    string
	env      string
)

func initFlag() {
	flag.Int64Var(&interval, "interval", 1, "时间间隔")
	flag.StringVar(&query, "query", "", "query")
	flag.StringVar(&env, "env", "prod", "env")
	flag.Parse()
}

func main() {
	initFlag()
	ctx := context.Background()
	logger.Init(conf.Env(env))
	defer logger.Sync()

	pprof.Start()
	log := logger.SugarLogger
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	if query == "" {
		log.Warnw("query must be set!")
		return
	}
	log.Debugw("runtime", "runtime", runtime.GOOS)
	for range ticker.C {
		log.Infow("now crawling goods", "query", query)
		response, err := hermes.CrawlGoods(ctx, query)
		if err != nil {
			log.Errorw("crawler wrong!", "err", err.Error())
			continue
		}
		if response.Total > 0 {
			wg := sync.WaitGroup{}
			urls := make([]string, 0, int(response.Total))
			for i := 0; i < int(response.Total) && i < len(response.Products); i++ {
				product := response.Products[i]
				url := product.GetDetailURL()
				cmd := default_web_navigator.OpenURL(runtime.GOOS, url)
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := cmd.Run()
					if err != nil {
						log.Errorw("run cmd error", "err", err.Error())
					}
				}()
				log.Infow("now open",
					"title", product.Title,
					"url", url,
					"sku", product.SKU,
				)
				urls = append(urls, url)
			}
			log.Infow(query+" 有货了！", "count", response.Total, "urls", urls)
			wg.Wait()
			return
		}
	}
}
