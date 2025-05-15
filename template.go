package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
)

var noticeTemplate = `
<style>
  .card-container {
    width: min(90vw, 400px); 
    background: white;
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
    padding: min(4vw, 20px); 
    display: flex;
    flex-direction: column;
    gap: 12px;
    border-top: 10px solid #ffa200;
    border-bottom: 10px solid #ffa200;
  }
  
  .header-container {
    display: flex;
    align-items: center;
  }
  
  .avatar-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    width: min(15vw, 60px);
  }
  
  .avatar-image {
    width: min(14vw, 64px); 
    height: min(14vw, 64px);
    min-width: 36px;
    min-height: 36px;
    border-radius: 50%;
    overflow: hidden;
    flex-shrink: 0;
  }
  
  .avatar-title {
    font-size: clamp(10px, 2.5vw, 12px);
    color: #888;
    text-align: center;
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .content-container {
    flex: 1;
    padding-right: min(3vw, 15px);
    margin-left: 10px;
  }
  
  .content-title {
    font-weight: bold;
    font-size: clamp(10px, 4vw, 14px);
    margin-bottom: 2px;
    color: #000;
  }
  
  .content-subtitle {
    font-size: clamp(5px, 3.5vw, 10px);
    color: #666;
    line-height: 1.4;
  }
  
  .cover-image {
	width: auto;
    height: min(18vw, 80px);
    aspect-ratio: 2/3; 
    border-radius: 4px;
    overflow: hidden;
    flex-shrink: 0;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  }
  
  .media-info {
    background: #f8f8f8;
    border-radius: 8px;
    padding: 10px;
    font-size: 12px;
    color: #555;
    line-height: 1.5;
  }
  
  .summary-text {
    color: #666;
  }
  
  .people-container {
    display: flex;
    gap: 10px;
    overflow-x: auto;
    padding-bottom: 8px;
  }
  
  .person-item {
    flex-shrink: 0;
    width: 60px;
    text-align: center;
  }
  
  .person-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    overflow: hidden;
    margin: 0 auto 4px;
    background: #eee;
  }
  
  .person-name {
    font-size: 10px;
    color: #666;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .person-role {
    font-size: 8px;
    color: #999;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .progress-info {
    display: flex;
    justify-content: space-between;
    font-size: 10px;
    color: #777;
    margin-bottom: 4px;
  }
  
  .progress-bar {
    height: 4px;
    background: #eee;
    border-radius: 2px;
    overflow: hidden;
  }
  
  .progress-fill {
    height: 100%;
    background: #ffa200;
  }
  
  .placeholder-text {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #999;
    font-size: 12px;
  }
  
  .image-fit {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
</style>

<div class="card-container">
    <!-- Header Section -->
    <div class="header-container">
        <div class="avatar-container">
            <div class="avatar-image">
                <img src="{{.payload.Account.Thumb}}" alt="avatar" class="image-fit">
            </div>
            <div class="avatar-title">
                {{.payload.Account.Title}}
            </div>
        </div>
        <div class="content-container">
            <div class="content-title">{{.icon}} {{.payload.Metadata.LibrarySectionTitle}}: {{.payload.Metadata.Title}}</div>
            <div class="content-subtitle">From {{.payload.Player.PublicAddress}} on {{.payload.Player.Title}}</div>
        </div>
        {{if .thumb}}
        <div class="cover-image">
            <img src="data:{{.thumbType}};base64,{{.thumb}}" alt="cover" class="image-fit">
        </div>
        {{end}}
    </div>

    <!-- Media Info Section -->
    <div class="media-info">
        {{if .payload.Metadata.Year}}<div><strong>Year:</strong> {{.payload.Metadata.Year}}</div>{{end}}
        {{if .payload.Metadata.ContentRating}}<div><strong>Rating:</strong> {{.payload.Metadata.ContentRating}}</div>{{end}}
        {{if .payload.Metadata.Rating}}<div><strong>Score:</strong> {{printf "%.1f" .payload.Metadata.Rating}}/10</div>{{end}}
        {{if .payload.Metadata.Duration}}<div><strong>Duration:</strong> {{divide .payload.Metadata.Duration 60000}} minutes</div>{{end}}
        {{if .payload.Metadata.Summary}}
        <div style="margin-top:6px;">
            <div class="summary-text">{{.payload.Metadata.Summary}}</div>
        </div>
        {{end}}
    </div>

    <!-- Cast & Crew Section -->
    {{if or (and .payload.Metadata.Director (gt (len .payload.Metadata.Director) 0)) 
            (and .payload.Metadata.Writer (gt (len .payload.Metadata.Writer) 0))
            (and .payload.Metadata.Role (gt (len .payload.Metadata.Role) 0))}}
    <div style="margin-top:4px;">
        <div class="people-container">
            {{range .payload.Metadata.Director}}
            <div class="person-item">
                <div class="person-avatar">
                    {{if and .Thumb (ne .Thumb "")}}
                    <img src="{{.Thumb}}" alt="{{.Tag}}" class="image-fit">
                    {{else}}
                    <div class="placeholder-text">Director</div>
                    {{end}}
                </div>
                <div class="person-name">{{.Tag}}</div>
				<div class="person-role">Director</div>
            </div>
            {{end}}

            {{range .payload.Metadata.Writer}}
            <div class="person-item">
                <div class="person-avatar">
                    {{if and .Thumb (ne .Thumb "")}}
                    <img src="{{.Thumb}}" alt="{{.Tag}}" class="image-fit">
                    {{else}}
                    <div class="placeholder-text">Writer</div>
                    {{end}}
                </div>
                <div class="person-name">{{.Tag}}</div>
				<div class="person-role">Writer</div>
            </div>
            {{end}}

            {{range .payload.Metadata.Role}}
            <div class="person-item">
                <div class="person-avatar">
                    {{if and .Thumb (ne .Thumb "")}}
                    <img src="{{.Thumb}}" alt="{{.Tag}}" class="image-fit">
                    {{else}}
                    <div class="placeholder-text">Actor</div>
                    {{end}}
                </div>
                <div class="person-name">{{.Tag}}</div>
                <div class="person-role">{{.Role}}</div>
            </div>
            {{end}}
        </div>
    </div>
    {{end}}

    <!-- Progress Bar -->
    {{if .payload.Metadata.Duration}}
    <div style="margin-top:8px;">
        <div class="progress-info">
            <span>{{divide .payload.Metadata.ViewOffset 60000}}m</span>
            <span>{{divide (subtract .payload.Metadata.Duration .payload.Metadata.ViewOffset) 60000}}m left</span>
            <span>{{divide .payload.Metadata.Duration 60000}}m</span>
        </div>
        <div class="progress-bar">
            <div class="progress-fill" style="width:{{percent .payload.Metadata.ViewOffset .payload.Metadata.Duration}}%;"></div>
        </div>
    </div>
    {{end}}
</div>
`

var notify = template.Must(template.New("plex").
	Funcs(template.FuncMap{
		"divide": func(a, b int) int {
			return a / b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"percent": func(a, b int) int {
			return int(float64(a) / float64(b) * 100)
		},
	}).
	Parse(noticeTemplate))

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
	summary = fmt.Sprintf("%s %s %s: %s", icon, ev.Payload.Account.Title, ev.Payload.Metadata.LibrarySectionTitle, ev.Payload.Metadata.Title)
	content = buffer.String()
	return
}
