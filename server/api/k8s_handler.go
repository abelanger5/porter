package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/porter-dev/porter/internal/kubernetes"
	v1 "k8s.io/api/core/v1"

	"github.com/gorilla/websocket"
	"github.com/porter-dev/porter/internal/forms"
)

// Enumeration of k8s API error codes, represented as int64
const (
	ErrK8sDecode ErrorCode = iota + 600
	ErrK8sValidate
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HandleListNamespaces retrieves a list of namespaces
func (app *App) HandleListNamespaces(w http.ResponseWriter, r *http.Request) {
	vals, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		app.handleErrorFormDecoding(err, ErrReleaseDecode, w)
		return
	}

	// get the filter options
	form := &forms.K8sForm{
		OutOfClusterConfig: &kubernetes.OutOfClusterConfig{
			Repo: app.repo,
		},
	}

	form.PopulateK8sOptionsFromQueryParams(vals, app.repo.Cluster)

	// validate the form
	if err := app.validator.Struct(form); err != nil {
		app.handleErrorFormValidation(err, ErrK8sValidate, w)
		return
	}

	// create a new agent
	var agent *kubernetes.Agent

	if app.testing {
		agent = app.TestAgents.K8sAgent
	} else {
		agent, err = kubernetes.GetAgentOutOfClusterConfig(form.OutOfClusterConfig)
	}

	namespaces, err := agent.ListNamespaces()

	if err != nil {
		app.handleErrorDataRead(err, w)
		return
	}

	if err := json.NewEncoder(w).Encode(namespaces); err != nil {
		app.handleErrorFormDecoding(err, ErrK8sDecode, w)
		return
	}
}

// HandleGetPodLogs returns real-time logs of the pod via websockets
// TODO: Refactor repeated calls.
func (app *App) HandleGetPodLogs(w http.ResponseWriter, r *http.Request) {

	// get session to retrieve correct kubeconfig
	_, err := app.store.Get(r, app.cookieName)

	// get path parameters
	namespace := chi.URLParam(r, "namespace")
	podName := chi.URLParam(r, "name")

	if err != nil {
		app.handleErrorFormDecoding(err, ErrReleaseDecode, w)
		return
	}

	vals, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		app.handleErrorFormDecoding(err, ErrReleaseDecode, w)
		return
	}

	// get the filter options
	form := &forms.K8sForm{
		OutOfClusterConfig: &kubernetes.OutOfClusterConfig{
			Repo: app.repo,
		},
	}

	form.PopulateK8sOptionsFromQueryParams(vals, app.repo.Cluster)

	// validate the form
	if err := app.validator.Struct(form); err != nil {
		app.handleErrorFormValidation(err, ErrK8sValidate, w)
		return
	}

	// create a new agent
	var agent *kubernetes.Agent

	if app.testing {
		agent = app.TestAgents.K8sAgent
	} else {
		agent, err = kubernetes.GetAgentOutOfClusterConfig(form.OutOfClusterConfig)
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade to websocket.
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		app.handleErrorUpgradeWebsocket(err, w)
	}

	err = agent.GetPodLogs(namespace, podName, conn)

	if err != nil {
		app.handleErrorWebsocketWrite(err, w)
		return
	}
}

// HandleListPods returns all pods that match the given selectors
// TODO: Refactor repeated calls.
func (app *App) HandleListPods(w http.ResponseWriter, r *http.Request) {

	// get session to retrieve correct kubeconfig
	_, err := app.store.Get(r, app.cookieName)

	if err != nil {
		app.handleErrorFormDecoding(err, ErrReleaseDecode, w)
		return
	}

	vals, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		app.handleErrorFormDecoding(err, ErrReleaseDecode, w)
		return
	}

	// get the filter options
	form := &forms.K8sForm{
		OutOfClusterConfig: &kubernetes.OutOfClusterConfig{
			Repo: app.repo,
		},
	}

	form.PopulateK8sOptionsFromQueryParams(vals, app.repo.Cluster)

	// validate the form
	if err := app.validator.Struct(form); err != nil {
		app.handleErrorFormValidation(err, ErrK8sValidate, w)
		return
	}

	// create a new agent
	var agent *kubernetes.Agent

	if app.testing {
		agent = app.TestAgents.K8sAgent
	} else {
		agent, err = kubernetes.GetAgentOutOfClusterConfig(form.OutOfClusterConfig)
	}

	pods := []v1.Pod{}
	for _, selector := range vals["selectors"] {
		podsList, err := agent.GetPodsByLabel(selector)

		if err != nil {
			app.handleErrorFormValidation(err, ErrK8sValidate, w)
			return
		}

		for _, pod := range podsList.Items {
			pods = append(pods, pod)
		}
	}

	if err := json.NewEncoder(w).Encode(pods); err != nil {
		app.handleErrorFormDecoding(err, ErrK8sDecode, w)
		return
	}
}

// HandleStreamControllerStatus test calls
// TODO: Refactor repeated calls.
func (app *App) HandleStreamControllerStatus(w http.ResponseWriter, r *http.Request) {

	// get session to retrieve correct kubeconfig
	_, err := app.store.Get(r, app.cookieName)

	if err != nil {
		app.handleErrorFormDecoding(err, ErrReleaseDecode, w)
		return
	}

	vals, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		app.handleErrorFormDecoding(err, ErrReleaseDecode, w)
		return
	}

	// get the filter options
	form := &forms.K8sForm{
		OutOfClusterConfig: &kubernetes.OutOfClusterConfig{
			Repo: app.repo,
		},
	}

	form.PopulateK8sOptionsFromQueryParams(vals, app.repo.Cluster)

	// validate the form
	if err := app.validator.Struct(form); err != nil {
		app.handleErrorFormValidation(err, ErrK8sValidate, w)
		return
	}

	// create a new agent
	var agent *kubernetes.Agent

	if app.testing {
		agent = app.TestAgents.K8sAgent
	} else {
		agent, err = kubernetes.GetAgentOutOfClusterConfig(form.OutOfClusterConfig)
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade to websocket.
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		app.handleErrorUpgradeWebsocket(err, w)
	}

	// get path parameters
	kind := chi.URLParam(r, "kind")
	err = agent.StreamControllerStatus(conn, kind)

	if err != nil {
		app.handleErrorWebsocketWrite(err, w)
		return
	}
}
