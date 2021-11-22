package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/mroobert/my-service/business/sys/validate"
	"github.com/mroobert/my-service/foundation/web"
)

// Test handler is for development.
func Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if n := rand.Intn(100); n%2 == 0 {
		//return errors.New("untrusted error")
		return validate.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
		//return web.NewShutdownError("restart service")
	}

	status := struct {
		Status string
	}{
		Status: "Ok",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
