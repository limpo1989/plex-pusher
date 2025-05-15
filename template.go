package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
)

var noticeTemplate = `
<div style="
        width: min(90vw, 400px); 
        background:white;
        border-radius:12px;
        box-shadow:0 4px 12px rgba(0,0,0,0.1);
        padding:min(4vw, 20px); 
        display:flex;
        align-items:center;
        border-left:20px solid #ffa200;
    ">
	<div style="
            display:flex;
            flex-direction:column;
            align-items:center;
            gap:2px;
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
				<img src="{{.payload.Account.Thumb}}" alt="头像" style="width:100%%;height:100%%;object-fit:cover;">
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
				{{.payload.Account.Title}}
			</div>
		</div>
        <div style="flex:1;padding-right:min(3vw, 15px);margin-left:10px;">
            <div style="font-weight:bold;font-size:clamp(10px, 4vw, 14px);margin-bottom:2px;color:#000;">{{.icon}} {{.payload.Metadata.LibrarySectionTitle}}: {{.payload.Metadata.Title}}</div>
            <div style="font-size:clamp(5px, 3.5vw, 10px);color:#666;line-height:1.4;">From {{.payload.Player.PublicAddress}} on {{.payload.Player.Title}}</div>
        </div>
		{{if .thumb}}
      	<div style="
            width: min(14vw, 60px);
            height: min(14vw, 60px);
            border-radius:2px;
            overflow:hidden;
            flex-shrink:0;
            box-shadow:0 2px 8px rgba(0,0,0,0.1);
        ">
            <img src="data:{{.thumbType}};base64,{{.thumb}}" alt="专辑封面" style="width:100%;height:100%;object-fit:cover;">
		</div>
		{{end}}
    </div>
`

var notify = template.Must(template.New("").Parse(noticeTemplate))
var icons = map[Event]string{
	MediaPlay:   "▶️",
	MediaPause:  "⏸️",
	MediaResume: "⏯️",
	MediaStop:   "⏹️",
	MediaRate:   "🌟",
	LibraryNew:  "🎞️",
}

func RenderNotice(ev PlexEvent) (summary, content string, err error) {
	var icon = icons[Event(ev.Payload.Event)]

	var args = map[string]any{
		"title":     fmt.Sprintf("%s %s", icon, ev.Payload.Metadata.Title),
		"thumb":     base64.StdEncoding.EncodeToString(ev.Thumb.Img),
		"thumbType": ev.Thumb.ImgType,
		"icons":     icons,
		"icon":      icon,
		"payload":   ev.Payload,
	}

	var buffer bytes.Buffer
	if err = notify.Execute(&buffer, args); nil != err {
		return
	}
	summary = fmt.Sprintf("%s %s: %s", icon, ev.Payload.Account.Title, ev.Payload.Metadata.Title)
	content = buffer.String()
	return
}
