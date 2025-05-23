package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/samber/lo"
	"github.com/wxpusher/wxpusher-sdk-go"

	"github.com/urfave/cli/v2"
	"github.com/wxpusher/wxpusher-sdk-go/model"
)

func main() {
	ppApp := cli.NewApp()
	ppApp.Name = "plex-pusher"
	ppApp.Usage = "plex event pusher"
	ppApp.Version = "0.0.1"
	ppApp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "listen",
			Usage:   "Specify the host:port to run pusher",
			Value:   "0.0.0.0:9876",
			EnvVars: []string{"PLEX_PUSHER_LISTEN"},
		},
		&cli.StringFlag{
			Name:     "token",
			Usage:    "WxPusher app token",
			Required: true,
			EnvVars:  []string{"PLEX_PUSHER_TOKEN"},
		},
		&cli.StringFlag{
			Name:     "uid",
			Usage:    "WxPusher app uid",
			Required: true,
			EnvVars:  []string{"PLEX_PUSHER_UID"},
		},
		&cli.StringSliceFlag{
			Name:    "events",
			Usage:   "Specify the events to be pushed",
			EnvVars: []string{"PLEX_PUSH_EVENT"},
			Value:   cli.NewStringSlice("*"),
		},
	}
	ppApp.Action = func(c *cli.Context) error {
		var listen = c.String("listen")
		var token = c.String("token")
		var uid = c.String("uid")
		var events = c.StringSlice("events")

		log.SetOutput(os.Stdout)

		return NewWebHook(listen).
			On(Any, PushMessage(token, uid, events)).
			Serve()
	}

	if err := ppApp.Run(os.Args); nil != err {
		log.Fatal(err)
	}
}

func PushMessage(token, uid string, events []string) func(ctx context.Context, event PlexEvent) error {
	return func(ctx context.Context, event PlexEvent) error {

		if !lo.Contains(events, event.Payload.Event) && !lo.Contains(events, string(Any)) {
			log.Printf("ignore event: %s from %s", event.Payload.Event, event.Payload.Server.Title)
			return nil
		}

		// render message
		summary, content, err := RenderNotice(event)
		if nil != err {
			log.Printf("error render message: %v", err)
			return err
		}

		// build message
		msg := model.NewMessage(token)
		msg.AddUId(uid)
		msg.SetContentType(2)
		msg.SetSummary(summary)
		msg.SetContent(content)

		// push message
		if _, err = wxpusher.SendMessage(msg); err != nil {
			log.Printf("error push message: %v", err)
			return err
		}

		log.Println(msg.Summary, " - ", fmt.Sprintf("From %s on %s", event.Payload.Player.PublicAddress, event.Payload.Player.Title))
		return nil
	}
}
