//go:build !solution

package requestlog

import (
	"crypto/rand"
	"net/http"

	"go.uber.org/zap"
	"github.com/felixge/httpsnoop"
)

func Log(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				var m httpsnoop.Metrics
				
				defer func() {
					v := recover()
					msg := "request finished"
					if v != nil {
						msg = "request panicked"
					}
					requestId := rand.Text()

				    l.Info("request started",
				    	zap.String("method", r.Method),
				    	zap.String("path", r.URL.Path),
				    	zap.String("request_id", requestId),
	 		        )

					l.Info(
						msg,
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("request_id", requestId),
						zap.Int("status_code", m.Code),
						zap.Duration("duration", m.Duration),
					)

					if v != nil { panic(v) }
				}()
				
				m = httpsnoop.CaptureMetrics(next, w, r)
			},
		)
	}
}
