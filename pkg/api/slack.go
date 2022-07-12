package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func (a *ServerAPI) InitSlackWebhook() {
	a.Router.APIRoot.Handle("/slack/oauth", http.HandlerFunc(a.OAuthWebhook))
	a.Router.APIRoot.Handle("/slack", http.HandlerFunc(a.SlackWebhook))
}

// SlackWebhook func is used to handle slack event API callback
func (a *ServerAPI) OAuthWebhook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	oAuthErr, ok := getFormParam(r, "error")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if oAuthErr != "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	code, ok := getFormParam(r, "code")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	state, ok := getFormParam(r, "state")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.App.HandleOauth(code, state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// SlackWebhook func is used to handle slack event API callback
func (a *ServerAPI) SlackWebhook(w http.ResponseWriter, r *http.Request) {
	signingSecret := a.Config.SlackConfig.SigningSecret
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}

	err = a.App.HandleSlackMessage(eventsAPIEvent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getFormParam(r *http.Request, key string) (string, bool) {
	param, ok := r.Form[key]
	if !ok || len(param) != 1 {
		return "", false
	}
	return param[0], true
}
