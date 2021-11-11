package testgrp

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

// Test handler is for development.
func (h Handlers) Test(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string
	}{
		Status: "Ok",
	}
	json.NewEncoder(w).Encode(status)

	h.Log.Infow("test", "statusCode", http.StatusOK, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}
