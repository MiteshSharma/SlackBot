package app

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/MiteshSharma/SlackBot/pkg/model"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	slackOauthURL = "https://slack.com/api/oauth.v2.access"
)

var (
	//go:embed views/home.json
	homeViewJson []byte
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
		case *slackevents.AppHomeOpenedEvent:
			// Make sure that this user exists
			// Make sure the team exists in DB
			// Update the home page
			blocks := &slack.Blocks{}
			err = blocks.UnmarshalJSON(homeViewJson)
			a.BotNotify.PublishView(ev.User, slack.HomeTabViewRequest{
				Type:            slack.VTHomeTab,
				Blocks:          *blocks,
				PrivateMetadata: "",
				CallbackID:      "ViewHomeCallbackID",
			}, "")
		}
	}
	return err
}

func (a *App) HandleSlashCommand(slashCommand slack.SlashCommand) error {
	var err error = nil

	input, err := a.CreateDialog()
	if err != nil {
		return err
	}

	err = a.BotNotify.SendDialog(context.Background(), input, slashCommand.TriggerID)
	if err != nil {
		return err
	}

	return err
}

func (a *App) ShowDialog(triggerId string) error {
	var err error = nil
	input, err := a.CreateDialog()
	if err != nil {
		return err
	}

	err = a.BotNotify.SendDialog(context.Background(), input, triggerId)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) handleAppMentionEvent(workspaceID string, event slackevents.AppMentionEvent) error {
	// storageResult := a.Repository.Workspace().GetWorkspace(strings.ToLower(workspaceID))
	// if storageResult.Err != nil {
	// 	return fmt.Errorf(storageResult.Err.Message)
	// }
	// workspace := storageResult.Data.(*model.Workspace)
	// api := slack.New(workspace.AccessToken)
	// authResp, err := api.AuthTest()
	// if err != nil {
	// 	a.Log.Error(fmt.Sprint("auth request test failed with err: %w", err))
	// 	return fmt.Errorf("auth request test failed with err: %w", err)
	// }
	// botID := authResp.UserID
	// if !strings.HasPrefix(event.Text, "<@"+botID+">") {
	// 	a.Log.Debug(fmt.Sprint("Ignoring message as it doesn't contain %w prefix", botID))
	// 	return nil
	// }

	// user, err := api.GetUserInfo(event.User)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// 	return nil
	// }
	// fmt.Println(user.RealName)

	// message := strings.TrimPrefix(event.Text, "<@"+botID+">")

	a.BotNotify.SendAttachmentMessage(context.Background(), getSlackMessage("welcome message"), "test-channel")

	return nil
}

func (a *App) CreateDialog() (slack.Dialog, error) {
	dialog := slack.Dialog{}
	dialog.CallbackID = "test"
	dialog.Title = "Request Title"
	dialog.SubmitLabel = "Request"
	dialog.NotifyOnCancel = true

	// Unmarshall and append the text element
	textElement := &slack.TextInputElement{}
	textElement.Label = "testing label"
	textElement.Name = "testng name"
	textElement.Type = "text"
	textElement.Placeholder = "Enter value"
	textElement.Hint = "testing hint"

	dialog.Elements = []slack.DialogElement{
		textElement,
	}

	return dialog, nil
}

func getSlackMessage(message string) slack.Attachment {
	return slack.Attachment{
		CallbackID: "werdd",
		Fields: []slack.AttachmentField{
			{

				Title: "Message",
				Value: message,
				Short: true,
			},
		},
		Actions: []slack.AttachmentAction{
			{
				Name:  "Execute",
				Text:  "Exec Make",
				Type:  "button",
				Style: "primary",
			},
		},
		Footer: "Message from Bot",
	}
}
