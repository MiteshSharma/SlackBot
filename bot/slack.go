package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/MiteshSharma/SlackBot/logger"
	"github.com/slack-go/slack"
)

type SlackBot struct {
	log logger.Logger

	Token       string
	ChannelName string
	BotID       string
}

type slackMessage struct {
	log logger.Logger

	Event       *slack.MessageEvent
	BotID       string
	Request     string
	Response    string
	RTM         *slack.RTM
	SlackClient *slack.Client
}

func NewSlackBot(parentLogger logger.Logger, token string, channelName string) *SlackBot {
	return &SlackBot{
		log:         parentLogger,
		Token:       token,
		ChannelName: channelName,
	}
}

func (sb *SlackBot) Start() error {
	ctx := context.Background()
	sb.log.Info("Starting bot")
	api := slack.New(sb.Token)
	authResp, err := api.AuthTest()
	if err != nil {
		sb.log.Error(fmt.Sprint("auth request test failed with err: %w", err))
		return fmt.Errorf("auth request test failed with err: %w", err)
	}
	botID := authResp.UserID
	sb.log.Debug(fmt.Sprint("botId: %w", botID))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case <-ctx.Done():
			sb.log.Info("Context closed, Finishing")
			return rtm.Disconnect()
		case msg, ok := <-rtm.IncomingEvents:
			if !ok {
				sb.log.Info("Incoming events channel closed. Finishing")
				return nil
			}

			switch msgEvent := msg.Data.(type) {
			case *slack.ConnectedEvent:
				sb.log.Info("We are connected to Slack! Yayyy")

			case *slack.MessageEvent:
				sb.log.Info("Slack bot got some message to handle")
				if msgEvent.User == botID {
					// Skip self message
					continue
				}
				sm := slackMessage{
					log:         sb.log,
					Event:       msgEvent,
					BotID:       botID,
					RTM:         rtm,
					SlackClient: api,
				}
				err := sm.HandleMessage(sb)
				if err != nil {
					wrappedErr := fmt.Errorf("while handling message: %w", err)
					sb.log.Error(wrappedErr.Error())
				}

			case *slack.RTMError:
				sb.log.Error(fmt.Sprint("slack RMT error: %w", msgEvent.Error()))

			case *slack.ConnectionErrorEvent:
				sb.log.Error(fmt.Sprint("Slack connection error: %w", msgEvent.Error()))

			case *slack.IncomingEventError:
				sb.log.Error(fmt.Sprint("Slack incoming event error: %w", msgEvent.Error()))

			case *slack.OutgoingErrorEvent:
				sb.log.Error(fmt.Sprint("Slack outgoing event error: %w", msgEvent.Error()))

			case *slack.UnmarshallingErrorEvent:
				sb.log.Error(fmt.Sprint("Slack unmarshalling error: %w", msgEvent.Error()))

			case *slack.RateLimitedError:
				sb.log.Error(fmt.Sprint("Slack rate limiting error: %w", msgEvent.Error()))

			case *slack.InvalidAuthEvent:
				return fmt.Errorf("invalid auth credentials")
			}
		}
	}
}

func (sm *slackMessage) HandleMessage(sb *SlackBot) error {
	info, err := sm.SlackClient.GetConversationInfo(sm.Event.Channel, true)
	var IsMainChannel bool
	if err == nil {
		if info.IsChannel || info.IsPrivate {
			// Message posted in a channel, serve if starts with mention of bot
			if !strings.HasPrefix(sm.Event.Text, "<@"+sm.BotID+">") {
				sm.log.Debug(fmt.Sprint("Ignoring message as it doesn't contain %w prefix", sm.BotID))
				return nil
			}

			// Serve only if current channel is in config
			if sb.ChannelName == info.Name {
				IsMainChannel = true
			}
		}
	}
	// Serve only if current channel is our channel to handle
	if sb.ChannelName == sm.Event.Channel {
		IsMainChannel = true
	}

	sm.Request = strings.TrimPrefix(sm.Event.Text, "<@"+sm.BotID+">")

	if IsMainChannel {
		sm.Response = "Response"
		err = sm.Send()
		if err != nil {
			return fmt.Errorf("while sending message: %w", err)
		}
	}

	return nil
}

func (sm *slackMessage) Send() error {
	sm.log.Debug(fmt.Sprint("Slack incoming Request: %w", sm.Request))
	sm.log.Debug(fmt.Sprint("Slack Response: %w", sm.Response))
	if len(sm.Response) == 0 {
		return fmt.Errorf("slack response: empty response %q", sm.Request)
	}
	// Upload message as a file if too long
	MAX_LIMIT_FOR_DIRECT_MSG := 4000
	if len(sm.Response) >= MAX_LIMIT_FOR_DIRECT_MSG {
		params := slack.FileUploadParameters{
			Filename: sm.Request,
			Title:    sm.Request,
			Content:  sm.Response,
			Channels: []string{sm.Event.Channel},
		}
		_, err := sm.RTM.UploadFile(params)
		if err != nil {
			return fmt.Errorf("error while uploading file: %w", err)
		}
		return nil
	}

	var options = []slack.MsgOption{slack.MsgOptionText(strings.TrimSpace(sm.Response), false), slack.MsgOptionAsUser(true)}

	if sm.Event.ThreadTimestamp != "" {
		//if the message is from thread then add an option to return the response to the thread
		options = append(options, slack.MsgOptionTS(sm.Event.ThreadTimestamp))
	}

	if _, _, err := sm.RTM.PostMessage(sm.Event.Channel, options...); err != nil {
		return fmt.Errorf("error while posting Slack message: %w", err)
	}

	return nil
}
