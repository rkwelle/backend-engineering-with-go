package main

import (
	"net/http"
)

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "Ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err == nil {
		app.internalServerError(w, r, err)
		// writeJSONError(w, http.StatusInternalServerError, "err.Error()")
	}
}
