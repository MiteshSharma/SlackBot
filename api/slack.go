package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/MiteshSharma/SlackBot/notify"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func (a *ServerAPI) InitSlackWebhook() {
	a.Router.APIRoot.Handle("/slack", http.HandlerFunc(a.SlackWebhook))
}

// SlackWebhook func is used to handle slack event API callback
func (a *ServerAPI) SlackWebhook(w http.ResponseWriter, r *http.Request) {
	signingSecret := a.Config.SlackConfig.SigningSecret
	fmt.Println(signingSecret)
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

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			a.handleSlackMessage(*ev)
		}
	}
}

func (a *ServerAPI) handleSlackMessage(event slackevents.AppMentionEvent) error {
	api := slack.New(a.Config.SlackConfig.Token)
	authResp, err := api.AuthTest()
	if err != nil {
		a.Log.Error(fmt.Sprint("auth request test failed with err: %w", err))
		return fmt.Errorf("auth request test failed with err: %w", err)
	}
	botID := authResp.UserID
	if !strings.HasPrefix(event.Text, "<@"+botID+">") {
		a.Log.Debug(fmt.Sprint("Ignoring message as it doesn't contain %w prefix", botID))
		return nil
	}

	message := strings.TrimPrefix(event.Text, "<@"+botID+">")

	a.BotNotify.SendEvent(a.Context, notify.Event{
		Message:      message,
		Channel:      "test-channel",
		IsAttachment: false,
	})

	return nil
}
