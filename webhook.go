package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"

	"go-spring.dev/web"
)

type EventHandler func(ctx context.Context, event PlexEvent) error

type WebHook struct {
	webServer       *web.Server
	enableAccessLog bool
	observes        map[Event]EventHandler
}

func NewWebHook(addr string) *WebHook {
	return &WebHook{
		webServer: web.NewServer(web.Options{Addr: addr}),
		observes:  make(map[Event]EventHandler),
	}
}

func (wh *WebHook) On(event Event, handler EventHandler) *WebHook {
	wh.observes[event] = handler
	return wh
}

func (wh *WebHook) EnableAccessLog(enable bool) *WebHook {
	wh.enableAccessLog = enable
	return wh
}

func (wh *WebHook) Serve() error {

	if wh.enableAccessLog {
		// access log
		wh.webServer.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if nil != err {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				log.Printf("Access: method: %s, path: %s, address: %s, body: %s", r.Method, r.RequestURI, r.RemoteAddr, body)

				r.Body = io.NopCloser(bytes.NewBuffer(body))
				next.ServeHTTP(w, r)
			})
		})
	}

	wh.webServer.Post("/events", wh.OnEvents)
	log.Printf("Listening on http://%s/events", wh.webServer.Addr())
	return wh.webServer.Run()
}

func (wh *WebHook) Shutdown() {
	wh.webServer.Shutdown(context.Background())
}

func (wh *WebHook) OnEvents(ctx context.Context, event PlexEvent) error {
	// match events.
	if handler, ok := wh.observes[Event(event.Event)]; ok {
		return handler(ctx, event)
	}

	// any events
	if handler, ok := wh.observes[Any]; ok {
		return handler(ctx, event)
	}

	// ignore events.
	log.Printf("ignore event: %s from %s", event.Event, event.Server.Title)
	return nil
}
