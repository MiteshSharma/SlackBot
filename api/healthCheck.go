package api

import (
	"io"
	"net/http"
)

func (a *ServerAPI) InitHealthCheck() {
	a.Router.APIRoot.Handle("/healthCheck", http.HandlerFunc(a.HealthCheck))
}

// HealthCheck func is used to check health check status
func (a *ServerAPI) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check for status of DB or cache (Redis) in future
	io.WriteString(w, `{"alive": true}`)
}
