## GO_GAE_PUBSUB
This is GAE Standard application, which enables asynchronyzed tasks with Pub/Sub.
GAE publish messages to Pub/Sub, and 2 subscriptions reveive them.
Each subsctiption send messages to Slack Channel, which you set up tokens in app.yml.

```
                          / GAE -> send_to_slack_ch
reqests -> GAE -> Pub/Sub
                          \ GAE -> send_to_slack_ch
```

## Setup GAE Standard

## Setup Pub/Sub with push type

## Deploy

