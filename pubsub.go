// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample pubsub demonstrates use of the cloud.google.com/go/pubsub package from App Engine flexible environment.
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

	"golang.org/x/net/context"

  "github.com/nlopes/slack"
  "google.golang.org/appengine/urlfetch"
)

var (
	topic *pubsub.Topic
)

func main() {
	http.HandleFunc("/pubsub/publish", publishHandler)
	http.HandleFunc("/pubsub/push", pushHandler)
	http.HandleFunc("/branch/push", branchHandler)
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
	Message struct { Attributes map[string]string
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

  imageUrl := "http://entert.jyuusya-yoshiko.com/wp/wp-content/uploads/2017/07/C7T-zq0V0AAmar0.jpg"
  sendSlack(r, imageUrl, mustGetenv("CHANNEL_TOKEN_1"))
}

func branchHandler(w http.ResponseWriter, r *http.Request) {
	msg := &pushRequest{}
	if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}

  imageUrl := "https://blog.golang.org/gopher/gopher.png"
  sendSlack(r, imageUrl, mustGetenv("CHANNEL_TOKEN_2"))
}

func sendSlack(r *http.Request, imageUrl string, channelToken string) {
	ctx := appengine.NewContext(r)
  api := slack.New(mustGetenv("SLACK_TOKEN"))
  slack.SetHTTPClient(urlfetch.Client(ctx))
  params := slack.PostMessageParameters{}
  attachment := slack.Attachment{
    Pretext: "some pretext",
    Text:    "some text",
    ImageURL: imageUrl,
  }
  params.Attachments = []slack.Attachment{attachment}
  _, _, err := api.PostMessage(channelToken, "title", params)
  if err != nil {
    appLog.Errorf(ctx, "%v", err)
    return
  }
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
  topic := getTopic(r)

	ctx := context.Background()

	msg := &pubsub.Message{
		Data: []byte(r.FormValue("payload")),
	}

	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		http.Error(w, fmt.Sprintf("Could not publish message: %v", err), 500)
		return
	}

	fmt.Fprint(w, "Message published.")
}

func getTopic(r *http.Request) *pubsub.Topic {
	ctx := appengine.NewContext(r)
	client, err := pubsub.NewClient(ctx, appengine.AppID(ctx))
	if err != nil {
		log.Fatal(err)
	}
	topic, _ = client.CreateTopic(ctx, mustGetenv("PUBSUB_TOPIC"))
  return topic
}
