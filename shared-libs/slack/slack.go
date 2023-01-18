package slack

import (
	"context"
	"time"

	"github.com/slack-go/slack"
)

type (
	SlackOption struct {
		ServiceName string
		Token       string
		Debug       bool
		MaxWait     time.Duration
	}

	SlackMessage struct {
		AuthorName string
		Color      string
		Pretext    string
		Text       string
	}

	slackClient struct {
		serviceName string
		api         *slack.Client
		maxWait     time.Duration
	}
)

func NewSlackApi(opt *SlackOption) *slackClient {
	return &slackClient{
		serviceName: opt.ServiceName,
		api:         slack.New(opt.Token, slack.OptionDebug(opt.Debug)),
		maxWait:     opt.MaxWait,
	}
}

func (s *slackClient) PostMessageOnChannel(ctx context.Context, channelID string, msg SlackMessage) (err error) {
	ctxWT, cancel := context.WithTimeout(ctx, s.maxWait)
	defer cancel()

	_, _, err = s.api.PostMessageContext(
		ctxWT,
		channelID,
		slack.MsgOptionAsUser(true),
		slack.MsgOptionAttachments(slack.Attachment{
			ServiceName: s.serviceName,
			AuthorName:  msg.AuthorName,
			Color:       msg.Color,
			Pretext:     msg.Pretext,
			Text:        msg.Text,
		}),
	)

	return
}
