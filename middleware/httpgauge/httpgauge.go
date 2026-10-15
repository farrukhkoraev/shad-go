//go:build !solution

package httpgauge

import (
	//	"fmt"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/go-chi/chi/v5"
)


type Gauge struct{
	lock sync.RWMutex
	metrics map[string]int
	nextHandler http.Handler
}

func New() *Gauge {
	return &Gauge{ metrics: map[string]int{} }
}

func (g *Gauge) Snapshot() map[string]int {
	m := map[string]int{}
	g.lock.RLock()
	for k, v := range g.metrics {
		m[k] = v
	}
	g.lock.RUnlock()
	return m
}

// ServeHTTP returns accumulated statistics in text format ordered by pattern.
//
// For example:
//
//	/a 10
//	/b 5
//	/c/{id} 7
func (g *Gauge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	snapshot := g.Snapshot()
	keys := make([]string, 0, len(snapshot))
	for k := range snapshot {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%s %d\n", k, snapshot[k])
	}
	
}

func (g *Gauge) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				v := recover()
				pattern := chi.RouteContext(r.Context()).RoutePattern()
				g.lock.Lock()
				g.metrics[pattern] += 1
				g.lock.Unlock()

				if v != nil { panic(v) }
			}()
			next.ServeHTTP(w, r)
		})

}
