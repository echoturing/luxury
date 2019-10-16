package main

import (
	"context"
	"flag"
	"runtime"
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
			urls := make([]string, 0, int(response.Total))
			for i := 0; i < int(response.Total) && i < len(response.Products); i++ {
				url := response.Products[i].GetDetailURL()
				cmd := default_web_navigator.OpenURL(runtime.GOOS, url)
				go func() {
					err := cmd.Run()
					if err != nil {
						sugar.Errorw("run cmd error", "err", err.Error())
					}
				}()
				urls = append(urls, url)
			}
			sugar.Infow(query+" 有货了！", "count", response.Total, "urls", urls)
			return
		}
	}
}
