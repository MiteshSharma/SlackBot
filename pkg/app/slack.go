package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/MiteshSharma/SlackBot/notify"
	"github.com/MiteshSharma/SlackBot/pkg/model"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	slackOauthURL = "https://slack.com/api/oauth.v2.access"
)

type slackOauthResponse struct {
	OK          bool
	Error       string
	AccessToken string `json:"access_token"`
	Scope       string
	UserID      string `json:"user_id"`
	TeamName    string `json:"team_name"`
	TeamID      string `json:"team_id"`
}

func (a *App) HandleOauth(code string, state string) error {
	form := url.Values{}
	form.Add("client_id", a.Config.SlackConfig.ClientID)
	form.Add("client_secret", a.Config.SlackConfig.ClientSecret)
	form.Add("code", code)

	resp, err := http.PostForm(slackOauthURL, form)
	if err != nil {
		a.Log.Error(fmt.Sprint("error when posting oauth.access: %w", err))
		return err
	}
	defer resp.Body.Close()

	var oauthResponse slackOauthResponse
	err = json.NewDecoder(resp.Body).Decode(&oauthResponse)
	if err != nil {
		a.Log.Error(fmt.Sprint("error when deserializing oauth.access response: %w", err))
		return err
	}

	api := slack.New(oauthResponse.AccessToken)
	authResp, err := api.AuthTest()
	if err != nil {
		a.Log.Error(fmt.Sprint("auth request test failed with err: %w", err))
		return fmt.Errorf("auth request test failed with err: %w", err)
	}
	a.Log.Info(fmt.Sprint("auth test successful with userId: %w and teamId: %w", authResp.UserID, authResp.TeamID))

	storageResult := a.Repository.Workspace().GetWorkspace(strings.ToLower(authResp.TeamID))
	if storageResult.Err != nil {
		return fmt.Errorf(storageResult.Err.Message)
	}
	if storageResult.Data != nil {
		existingWorkspace := storageResult.Data.(*model.Workspace)
		if oauthResponse.AccessToken != "" {
			existingWorkspace.AccessToken = oauthResponse.AccessToken
		}

		storageResult = a.Repository.Workspace().UpdateWorkspace(existingWorkspace)
		if storageResult.Err != nil {
			return fmt.Errorf(storageResult.Err.Message)
		}
		a.Log.Info(fmt.Sprint("update workspace with new access token: %w", authResp.TeamID))
	} else {
		workspace := &model.Workspace{
			WorkspaceID: strings.ToLower(authResp.TeamID),
			OwnerID:     authResp.UserID,
			Name:        authResp.Team,
			AccessToken: oauthResponse.AccessToken,
		}

		storageResult := a.Repository.Workspace().CreateWorkspace(workspace)
		if storageResult.Err != nil {
			return fmt.Errorf(storageResult.Err.Message)
		}
	}

	return nil
}

func (a *App) HandleSlackMessage(eventsAPIEvent slackevents.EventsAPIEvent) error {
	var err error = nil
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			err = a.handleAppMentionEvent(eventsAPIEvent.TeamID, *ev)
		}
	}
	return err
}

func (a *App) handleAppMentionEvent(workspaceID string, event slackevents.AppMentionEvent) error {
	storageResult := a.Repository.Workspace().GetWorkspace(strings.ToLower(workspaceID))
	if storageResult.Err != nil {
		return fmt.Errorf(storageResult.Err.Message)
	}
	workspace := storageResult.Data.(*model.Workspace)
	api := slack.New(workspace.AccessToken)
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

	user, err := api.GetUserInfo(event.User)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}
	fmt.Println(user.RealName)

	message := strings.TrimPrefix(event.Text, "<@"+botID+">")

	a.BotNotify.SendEvent(a.Context, notify.Event{
		Message:      message,
		Channel:      "test-channel",
		IsAttachment: false,
	})

	return nil
}
