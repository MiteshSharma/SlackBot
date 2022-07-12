package app

import (
	"fmt"
	"strings"

	"github.com/MiteshSharma/SlackBot/notify"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func (a *App) HandleSlackMessage(eventsAPIEvent slackevents.EventsAPIEvent) error {
	var err error = nil
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			err = a.handleAppMentionEvent(*ev)
		}
	}
	return err
}

func (a *App) handleAppMentionEvent(event slackevents.AppMentionEvent) error {
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

	teamInfo, err := api.GetTeamInfo()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(teamInfo.ID)
	fmt.Println(event.User)

	user, err := api.GetUserInfo(event.User)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}
	fmt.Println(user.Name)

	message := strings.TrimPrefix(event.Text, "<@"+botID+">")

	a.BotNotify.SendEvent(a.Context, notify.Event{
		Message:      message,
		Channel:      "test-channel",
		IsAttachment: false,
	})

	return nil
}
