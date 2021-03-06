package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/porter-dev/porter/internal/forms"

	ints "github.com/porter-dev/porter/internal/models/integrations"
)

// HandleListClusterIntegrations lists the cluster integrations available to the
// instance
func (app *App) HandleListClusterIntegrations(w http.ResponseWriter, r *http.Request) {
	clusters := ints.PorterClusterIntegrations

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&clusters); err != nil {
		app.handleErrorFormDecoding(err, ErrProjectDecode, w)
		return
	}
}

// HandleListRegistryIntegrations lists the image registry integrations available to the
// instance
func (app *App) HandleListRegistryIntegrations(w http.ResponseWriter, r *http.Request) {
	registries := ints.PorterRegistryIntegrations

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&registries); err != nil {
		app.handleErrorFormDecoding(err, ErrProjectDecode, w)
		return
	}
}

// HandleListRepoIntegrations lists the repo integrations available to the
// instance
func (app *App) HandleListRepoIntegrations(w http.ResponseWriter, r *http.Request) {
	repos := ints.PorterRepoIntegrations

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&repos); err != nil {
		app.handleErrorFormDecoding(err, ErrProjectDecode, w)
		return
	}
}

// HandleCreateGCPIntegration creates a new GCP integration in the DB
func (app *App) HandleCreateGCPIntegration(w http.ResponseWriter, r *http.Request) {
	session, err := app.store.Get(r, app.cookieName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, _ := session.Values["user_id"].(uint)

	projID, err := strconv.ParseUint(chi.URLParam(r, "project_id"), 0, 64)

	if err != nil || projID == 0 {
		app.handleErrorFormDecoding(err, ErrProjectDecode, w)
		return
	}

	form := &forms.CreateGCPIntegrationForm{
		UserID:    userID,
		ProjectID: uint(projID),
	}

	// decode from JSON to form value
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		app.handleErrorFormDecoding(err, ErrProjectDecode, w)
		return
	}

	// validate the form
	if err := app.validator.Struct(form); err != nil {
		app.handleErrorFormValidation(err, ErrProjectValidateFields, w)
		return
	}

	// convert the form to a gcp integration
	gcp, err := form.ToGCPIntegration()

	if err != nil {
		app.handleErrorFormDecoding(err, ErrProjectDecode, w)
		return
	}

	// handle write to the database
	gcp, err = app.repo.GCPIntegration.CreateGCPIntegration(gcp)

	if err != nil {
		app.handleErrorDataWrite(err, w)
		return
	}

	app.logger.Info().Msgf("New gcp integration created: %d", gcp.ID)

	w.WriteHeader(http.StatusCreated)

	gcpExt := gcp.Externalize()

	if err := json.NewEncoder(w).Encode(gcpExt); err != nil {
		app.handleErrorFormDecoding(err, ErrProjectDecode, w)
		return
	}
}

// HandleCreateAWSIntegration creates a new AWS integration in the DB
func (app *App) HandleCreateAWSIntegration(w http.ResponseWriter, r *http.Request) {
	session, err := app.store.Get(r, app.cookieName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, _ := session.Values["user_id"].(uint)

	projID, err := strconv.ParseUint(chi.URLParam(r, "project_id"), 0, 64)

	if err != nil || projID == 0 {
		app.handleErrorFormDecoding(err, ErrProjectDecode, w)
		return
	}

	form := &forms.CreateAWSIntegrationForm{
		UserID:    userID,
		ProjectID: uint(projID),
	}

	// decode from JSON to form value
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		app.handleErrorFormDecoding(err, ErrProjectDecode, w)
		return
	}

	// validate the form
	if err := app.validator.Struct(form); err != nil {
		app.handleErrorFormValidation(err, ErrProjectValidateFields, w)
		return
	}

	// convert the form to a aws integration
	aws, err := form.ToAWSIntegration()

	if err != nil {
		app.handleErrorFormDecoding(err, ErrProjectDecode, w)
		return
	}

	// handle write to the database
	aws, err = app.repo.AWSIntegration.CreateAWSIntegration(aws)

	if err != nil {
		app.handleErrorDataWrite(err, w)
		return
	}

	app.logger.Info().Msgf("New aws integration created: %d", aws.ID)

	w.WriteHeader(http.StatusCreated)

	awsExt := aws.Externalize()

	if err := json.NewEncoder(w).Encode(awsExt); err != nil {
		app.handleErrorFormDecoding(err, ErrProjectDecode, w)
		return
	}
}
