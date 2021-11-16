package testgrp

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

// Test handler is for development.
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "Ok",
	}
	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		return err
	}

	h.Log.Infow("test", "statusCode", http.StatusOK, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

	return nil
}
