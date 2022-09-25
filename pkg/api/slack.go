package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func (a *ServerAPI) InitSlackWebhook() {
	a.Router.APIRoot.Handle("/slack/oauth", http.HandlerFunc(a.OAuthWebhook))
	a.Router.APIRoot.Handle("/slack/command", http.HandlerFunc(a.SlashCommand))
	a.Router.APIRoot.Handle("/slack/interaction", http.HandlerFunc(a.InteractionAction))
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// SlackWebhook func is used to handle slack event API callback
func (a *ServerAPI) SlashCommand(w http.ResponseWriter, r *http.Request) {
	signingSecret := a.Config.SlackConfig.SigningSecret

	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &sv))

	slackSlashCommand, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = a.App.HandleSlashCommand(slackSlashCommand)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// SlackWebhook func is used to handle slack event API callback
func (a *ServerAPI) InteractionAction(w http.ResponseWriter, r *http.Request) {
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

	// Parse request body
	str, _ := url.QueryUnescape(string(body))
	str = strings.Replace(str, "payload=", "", 1)
	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(str), &message); err != nil {
		fmt.Println(err)
		return
	}

	switch message.Type {
	case slack.InteractionTypeDialogSubmission:
		// Receive a notification of a dialog submission
		fmt.Println("Successfully receive a dialog submission.")
	case slack.InteractionTypeInteractionMessage:
		a.App.ShowDialog(message.TriggerID)
		fmt.Println("Successfully interacted with message.")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func getFormParam(r *http.Request, key string) (string, bool) {
	param, ok := r.Form[key]
	if !ok || len(param) != 1 {
		return "", false
	}
	return param[0], true
}
