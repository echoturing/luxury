package pprof

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/echoturing/luxury/logger"
)

func Start() {
	go func() {
		logger.SugarLogger.Infow("start pprof", "err", http.ListenAndServe("localhost:6060", nil))
	}()
}
