package main

import (
	"net"
	"os"
	"os/signal"

	bot "github.com/MiteshSharma/SlackBot/bot"
	"github.com/MiteshSharma/SlackBot/config"
	"github.com/MiteshSharma/SlackBot/logger"
)

func main() {
	config := config.GetConfig()
	logger := logger.NewLogger(config.LoggerConfig)

	err := sendSystemdNotification()

	if err != nil {
		panic("error")
	}

	slackBot := bot.NewSlackBot(logger, config.SlackConfig.Token, config.SlackConfig.ChannelName)
	slackBot.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	os.Exit(0)
}

func sendSystemdNotification() error {
	notifySocket := os.Getenv("NOTIFY_SOCKET")
	if notifySocket != "" {
		state := "READY=1"
		socketAddr := &net.UnixAddr{
			Name: notifySocket,
			Net:  "unixgram",
		}
		conn, err := net.DialUnix(socketAddr.Net, nil, socketAddr)
		if err != nil {
			return err
		}
		defer conn.Close()
		_, err = conn.Write([]byte(state))
		return err
	}
	return nil
}
