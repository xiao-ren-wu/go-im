package middleware

import (
	"context"
	"github.com/xiao-ren-wu/go-im/model"
	"github.com/xiao-ren-wu/loggo"
	"time"
)

var L *loggo.Loggers

func InitLoggo(conf *model.LoggoConf) {
	var err error
	L, err = loggo.NewLogger(
		loggo.WithCtxValue(func(ctx context.Context) map[string]interface{} {
			return map[string]any{
				"LOGID": ctx.Value("LOGID"),
			}
		}),
		loggo.WithReportCaller(),
		loggo.WithRotateLogs(&loggo.RotateLogsConfig{
			LogFilePrefix: conf.LogFilePrefix,
			RotationSize:  conf.RotationSize,
			RotationTime:  time.Duration(conf.RotationTime) * time.Hour,
			MaxAge:        time.Duration(conf.MaxAge) * time.Hour,
		}),
	)
	if err != nil {
		panic(err)
	}
}
