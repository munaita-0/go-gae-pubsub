package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"cloud.google.com/go/pubsub"
	"google.golang.org/appengine"
	appLog "google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"golang.org/x/net/context"
	"github.com/nlopes/slack"
)

func main() {
	http.HandleFunc("/pubsub/publish", publishHandler)
	http.HandleFunc("/pubsub/push", pushHandler)
	// http.HandleFunc("/branch/push", branchHandler)
	appengine.Main()
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

type pushRequest struct {
	Message struct {
		Attributes map[string]string
		Data       []byte
		ID         string `json:"message_id"`
	}
	Subscription string
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	msg := &pushRequest{}
	if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	s := NewSlack(w, r, mustGetenv("SLACK_TOKEN"), mustGetenv("IMAGE_URL_1"), mustGetenv("CHANNEL_TOKEN_1"))
	s.SendToChannel()
}


func publishHandler(w http.ResponseWriter, r *http.Request) {
	p := NewPubsub(w, r)
	p.PublishMessage(p.GetTopic())
}

type Pubsub struct {
	ctx       context.Context
	bCtx      context.Context
	r         *http.Request
	w         http.ResponseWriter
	topicName string
}

func NewPubsub(w http.ResponseWriter, r *http.Request) *Pubsub {
	p := &Pubsub{
		ctx:       appengine.NewContext(r),
		bCtx:      context.Background(),
		r:         r,
		w:         w,
		topicName: os.Getenv("PUBSUB_TOPIC"),
	}
	return p
}

func (p Pubsub) GetTopic() *pubsub.Topic {
	client, err := pubsub.NewClient(p.ctx, appengine.AppID(p.ctx))
	if err != nil {
		log.Fatal(err)
	}
	topic, _ := client.CreateTopic(p.ctx, p.topicName)
	return topic
}

func (p Pubsub) PublishMessage(topic *pubsub.Topic) {
	msg := &pubsub.Message{
		Data: []byte(p.r.FormValue("payload")),
	}

	if _, err := topic.Publish(p.bCtx, msg).Get(p.bCtx); err != nil {
		http.Error(p.w, fmt.Sprintf("Could not publish message: %v", err), 500)
		return
	}
}

type Slack struct {
	ctx          context.Context
	r            *http.Request
	w            http.ResponseWriter
	slackToken   string
	channelToken string
	imageUrl     string
}

func NewSlack(w http.ResponseWriter, r *http.Request, slackToken string, channelToken string, imageUrl string) *Slack {
	p := &Slack{
		ctx:          appengine.NewContext(r),
		r:            r,
		w:            w,
		slackToken:   slackToken,
		channelToken: channelToken,
		imageUrl:     imageUrl,
	}
	return p
}

func (s Slack) SendToChannel() {
	api := slack.New(s.slackToken)
	slack.SetHTTPClient(urlfetch.Client(s.ctx))
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Pretext:  "some pretext",
		Text:     "some text",
		ImageURL: s.imageUrl,
	}
	params.Attachments = []slack.Attachment{attachment}
	_, _, err := api.PostMessage(s.channelToken, "title", params)
	if err != nil {
		appLog.Errorf(s.ctx, "%v", err)
		return
	}
}
