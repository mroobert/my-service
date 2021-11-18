package testgrp

import (
	"context"
	"net/http"

	"github.com/mroobert/my-service/foundation/web"
)

// Test handler is for development.
func Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "Ok",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
