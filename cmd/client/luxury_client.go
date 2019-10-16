package main

import (
	"context"
	"flag"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/echoturing/luxury/crawlers/hermes"
	"github.com/echoturing/luxury/default_web_navigator"
)

var (
	interval int64
	query    string
)

func initFlag() {
	flag.Int64Var(&interval, "interval", 1, "时间间隔")
	flag.StringVar(&query, "query", "", "query")
	flag.Parse()
}

func main() {
	initFlag()
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	if query == "" {
		sugar.Warnw("query must be set!")
		return
	}
	sugar.Debugw("runtime", "runtime", runtime.GOOS)
	for range ticker.C {
		sugar.Infow("now crawling goods", "query", query)
		response, err := hermes.CrawlGoods(ctx, query)
		if err != nil {
			sugar.Errorw("crawler wrong!", "err", err.Error())
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
						sugar.Errorw("run cmd error", "err", err.Error())
					}
				}()
				sugar.Infow("now open",
					"title", product.Title,
					"url", url,
					"sku", product.SKU,
				)
				urls = append(urls, url)
			}
			sugar.Infow(query+" 有货了！", "count", response.Total, "urls", urls)
			wg.Wait()
			return
		}
	}
}
