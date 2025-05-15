package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/wxpusher/wxpusher-sdk-go"
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
		switch event.Event {
		case MediaPlay, MediaPause, MediaResume, MediaStop, MediaRate, LibraryNew:
			msg := BuildMessage(event)
			msg.AppToken = token
			msg.UIds = []string{uid}
			log.Println(msg.Summary, " - ", fmt.Sprintf("From %s on %s", event.Player.PublicAddress, event.Player.Title))
			_, err := wxpusher.SendMessage(msg)
			if err != nil {
				log.Printf("error push message: %v", err)
			}
			return err
		default:
			log.Printf("ignore event: %s from %s", event.Event, event.Server.Title)
			return nil
		}
	}
}

var icons = map[Event]string{
	MediaPlay:   "▶️",
	MediaPause:  "⏸️",
	MediaResume: "⏯️",
	MediaStop:   "⏹️",
	MediaRate:   "🌟",
	LibraryNew:  "🎞️",
}

func BuildMessage(event PlexEvent) *model.Message {
	var title = fmt.Sprintf("%s %s", icons[Event(event.Event)], event.Metadata.Title)
	var content = fmt.Sprintf(`
<div style="
        width: min(90vw, 400px); 
        background:white;
        border-radius:12px;
        box-shadow:0 4px 12px rgba(0,0,0,0.1);
        padding:min(4vw, 20px); 
        display:flex;
        align-items:center;
        border-left:25px solid #ffa200;
    ">
	<div style="
            display:flex;
            flex-direction:column;
            align-items:center;
            gap:4px;
            width: min(15vw, 60px);
        ">
			<div style="
				width: min(10vw, 44px); 
				height: min(10vw, 44px);
				min-width:36px;
				min-height:36px;
				border-radius:50%%;
				overflow:hidden;
				flex-shrink:0;
			">
				<img src="%s" alt="头像" style="width:100%%;height:100%%;object-fit:cover;">
			</div>
			<div style="
					font-size:clamp(10px, 2.5vw, 12px);
					color:#888;
					text-align:center;
					width:100%%;
					overflow:hidden;
					text-overflow:ellipsis;
					white-space:nowrap;
				">
				%s
			</div>
		</div>
        <div style="flex:1;padding-right:min(3vw, 15px);margin-left:10px;">
            <div style="font-weight:300;font-size:clamp(12px, 4vw, 16px);margin-bottom:2px;color:#000;">%s</div>
            <div style="font-size:clamp(10px, 3.5vw, 14px);color:#666;line-height:1.4;">%s</div>
        </div>

    </div>
`, event.Account.Thumb, event.Account.Title, title, fmt.Sprintf("From %s on %s", event.Player.PublicAddress, event.Player.Title))

	return &model.Message{
		ContentType: 2,
		Summary:     fmt.Sprintf("%s %s: %s", icons[Event(event.Event)], event.Account.Title, event.Metadata.Title),
		Content:     content,
	}
}
