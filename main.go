package main

import (
	"context"
	"fmt"
	"github.com/wxpusher/wxpusher-sdk-go"
	"log"
	"os"

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
		&cli.BoolFlag{
			Name:    "accesslog",
			Usage:   "Enable accesslog",
			EnvVars: []string{"PLEX_PUSHER_ACCESSLOG"},
			Value:   false,
		},
	}
	ppApp.Action = func(c *cli.Context) error {
		var listen = c.String("listen")
		var token = c.String("token")
		var uid = c.String("uid")
		var accesslog = c.Bool("accesslog")

		log.SetOutput(os.Stdout)

		return NewWebHook(listen).
			EnableAccessLog(accesslog).
			On(Any, PushMessage(token, uid)).
			Serve()
	}

	if err := ppApp.Run(os.Args); nil != err {
		log.Fatal(err)
	}
}

func PushMessage(token, uid string) func(ctx context.Context, event PlexEvent) error {
	return func(ctx context.Context, event PlexEvent) error {
		switch event.Payload.Event {
		case MediaPlay, MediaPause, MediaResume, MediaStop, MediaRate, LibraryNew:

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
		default:
			log.Printf("ignore event: %s from %s", event.Payload.Event, event.Payload.Server.Title)
			return nil
		}
	}
}
