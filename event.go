package main

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"

	"github.com/limpo1989/gcache"
)

type Event string

const (
	Any Event = "*"

	MediaPlay     = "media.play"     // Media starts playing. An appropriate poster is attached
	MediaPause    = "media.pause"    // Media playback pauses
	MediaResume   = "media.resume"   // Media playback resumes
	MediaStop     = "media.stop"     // Media playback stops
	MediaRate     = "media.rate"     // Media is rated. A poster is also attached to this event
	MediaScrobble = "media.scrobble" // Media is viewed (played past the 90% mark)

	LibraryOnDeck = "library.on.deck" // A new item is added that appears in the user’s On Deck
	LibraryNew    = "library.new"     // A new item is added to a library to which the user has access

	AdminDatabaseBackup    = "admin.database.backup"    // A database backup is completed successfully via Scheduled Tasks
	AdminDatabaseCorrupted = "admin.database.corrupted" // Corruption is detected in the server database
	DeviceNew              = "device.new"               // A device accesses the owner’s server for any reason, which may come from background connection testing and doesn’t necessarily indicate active browsing or playback
	PlaybackStarted        = "playback.started"         // Playback is started by a shared user for the server
)

type Payload struct {
	Event string `json:"event"` // Type of event (e.g., media.resume)
	User  bool   `json:"user"`  // Whether triggered by a user
	Owner bool   `json:"owner"` // Whether triggered by the owner

	Account struct {
		ID    int    `json:"id"`    // Account ID
		Thumb string `json:"thumb"` // URL to account thumbnail
		Title string `json:"title"` // Account display name
	} `json:"Account"`

	Server struct {
		Title string `json:"title"` // Server name
		UUID  string `json:"uuid"`  // Server unique identifier
	} `json:"Server"`

	Player struct {
		Local         bool   `json:"local"`         // Whether playing locally
		PublicAddress string `json:"publicAddress"` // Public IP address of player
		Title         string `json:"title"`         // Player name/type
		UUID          string `json:"uuid"`          // Player unique identifier
	} `json:"Player"`

	Metadata struct {
		LibrarySectionType    string  `json:"librarySectionType"`              // Type of library section (movie/show/etc)
		RatingKey             string  `json:"ratingKey"`                       // Unique key for rating
		Key                   string  `json:"key"`                             // Media key path
		GUID                  string  `json:"guid"`                            // Global unique identifier
		Slug                  string  `json:"slug"`                            // URL-friendly identifier
		Studio                string  `json:"studio,omitempty"`                // Production studio
		Type                  string  `json:"type"`                            // Media type
		Title                 string  `json:"title"`                           // Media title
		LibrarySectionTitle   string  `json:"librarySectionTitle"`             // Library section name
		LibrarySectionID      int     `json:"librarySectionID"`                // Library section ID
		LibrarySectionKey     string  `json:"librarySectionKey"`               // Library section key path
		OriginalTitle         string  `json:"originalTitle,omitempty"`         // Original title
		ContentRating         string  `json:"contentRating,omitempty"`         // Content rating (e.g., PG-13)
		Summary               string  `json:"summary,omitempty"`               // Plot summary
		Rating                float64 `json:"rating,omitempty"`                // Critic rating
		AudienceRating        float64 `json:"audienceRating,omitempty"`        // Audience rating
		ViewOffset            int     `json:"viewOffset"`                      // Playback position in ms
		LastViewedAt          int64   `json:"lastViewedAt"`                    // Last viewed timestamp
		Year                  int     `json:"year,omitempty"`                  // Release year
		Tagline               string  `json:"tagline,omitempty"`               // Tagline/slogan
		Thumb                 string  `json:"thumb,omitempty"`                 // Thumbnail URL
		Art                   string  `json:"art,omitempty"`                   // Artwork URL
		Duration              int     `json:"duration"`                        // Total duration in ms
		OriginallyAvailableAt string  `json:"originallyAvailableAt,omitempty"` // Original release date
		AddedAt               int64   `json:"addedAt"`                         // When added to library
		UpdatedAt             int64   `json:"updatedAt"`                       // When last updated
		AudienceRatingImage   string  `json:"audienceRatingImage,omitempty"`   // Audience rating image identifier
		ChapterSource         string  `json:"chapterSource,omitempty"`         // Source of chapters
		PrimaryExtraKey       string  `json:"primaryExtraKey,omitempty"`       // Key for primary extra content
		RatingImage           string  `json:"ratingImage,omitempty"`           // Rating image identifier

		Image []struct {
			Alt  string `json:"alt"`  // Image description
			Type string `json:"type"` // Image type (poster/background/etc)
			URL  string `json:"url"`  // Image URL
		} `json:"Image,omitempty"`

		UltraBlurColors struct {
			TopLeft     string `json:"topLeft"`     // Top left color hex
			TopRight    string `json:"topRight"`    // Top right color hex
			BottomRight string `json:"bottomRight"` // Bottom right color hex
			BottomLeft  string `json:"bottomLeft"`  // Bottom left color hex
		} `json:"UltraBlurColors,omitempty"`

		Genre []struct {
			ID     int    `json:"id"`              // Genre ID
			Filter string `json:"filter"`          // Filter parameter
			Tag    string `json:"tag"`             // Genre name
			Count  int    `json:"count,omitempty"` // Count of items in this genre
		} `json:"Genre,omitempty"`

		Country []struct {
			ID     int    `json:"id"`              // Country ID
			Filter string `json:"filter"`          // Filter parameter
			Tag    string `json:"tag"`             // Country name
			Count  int    `json:"count,omitempty"` // Count of items from this country
		} `json:"Country,omitempty"`

		Guid []struct {
			ID string `json:"id"` // External ID (IMDB/TMDB/etc)
		} `json:"Guid,omitempty"`

		Ratings []struct {
			Image string  `json:"image"`           // Rating image identifier
			Value float64 `json:"value"`           // Rating value
			Type  string  `json:"type"`            // Rating type (critic/audience)
			Count int     `json:"count,omitempty"` // Number of ratings
		} `json:"Rating,omitempty"`

		Director []struct {
			ID     int    `json:"id"`               // Director ID
			Filter string `json:"filter"`           // Filter parameter
			Tag    string `json:"tag"`              // Director name
			TagKey string `json:"tagKey,omitempty"` // Director unique key
			Count  int    `json:"count,omitempty"`  // Number of works
			Thumb  string `json:"thumb,omitempty"`  // Director thumbnail URL
		} `json:"Director,omitempty"`

		Writer []struct {
			ID     int    `json:"id"`               // Writer ID
			Filter string `json:"filter"`           // Filter parameter
			Tag    string `json:"tag"`              // Writer name
			TagKey string `json:"tagKey,omitempty"` // Writer unique key
			Count  int    `json:"count,omitempty"`  // Number of works
			Thumb  string `json:"thumb,omitempty"`  // Writer thumbnail URL
		} `json:"Writer,omitempty"`

		Role []struct {
			ID     int    `json:"id"`               // Actor ID
			Filter string `json:"filter"`           // Filter parameter
			Tag    string `json:"tag"`              // Actor name
			TagKey string `json:"tagKey,omitempty"` // Actor unique key
			Role   string `json:"role"`             // Character name
			Count  int    `json:"count,omitempty"`  // Number of works
			Thumb  string `json:"thumb,omitempty"`  // Actor thumbnail URL
		} `json:"Role,omitempty"`

		Producer []struct {
			ID     int    `json:"id"`               // Producer ID
			Filter string `json:"filter"`           // Filter parameter
			Tag    string `json:"tag"`              // Producer name
			TagKey string `json:"tagKey,omitempty"` // Producer unique key
			Count  int    `json:"count,omitempty"`  // Number of works
			Thumb  string `json:"thumb,omitempty"`  // Producer thumbnail URL
		} `json:"Producer,omitempty"`
	} `json:"Metadata"`
}

type PlexThumb struct {
	ThumbImg  []byte
	OriginImg []byte
	ImgType   string
}

type PlexEvent struct {
	Payload Payload
	Thumb   PlexThumb
}

// 专辑图缓存
var thumbCache = gcache.New[string, PlexThumb](1000).
	LRU().
	Build()

type RawPlexEvent struct {
	Payload string                `form:"payload"`
	Thumb   *multipart.FileHeader `form:"thumb"`
}

func (pe *RawPlexEvent) Parse() (event PlexEvent, err error) {
	if err = json.Unmarshal([]byte(pe.Payload), &event.Payload); nil != err {
		return
	}

	if pe.Thumb != nil && pe.Thumb.Size > 0 {
		var f multipart.File
		f, err = pe.Thumb.Open()
		if nil != err {
			return
		}
		defer f.Close()

		// 获取图片类型
		if event.Thumb.ImgType = pe.Thumb.Header.Get("Content-Type"); len(event.Thumb.ImgType) <= 0 {
			event.Thumb.ImgType = "image/jpeg"
		}

		// 读取图片数据
		if event.Thumb.OriginImg, err = io.ReadAll(f); nil != err {
			return
		}

		// 按照固定高度进行缩放
		event.Thumb.ThumbImg = resizeImage(event.Thumb.ImgType, event.Thumb.OriginImg, 0, 100)

		// 缓存起来
		_ = thumbCache.Set(event.Payload.Metadata.GUID, &event.Thumb)
	} else {
		// find in cache
		if thumb, err := thumbCache.Get(context.Background(), event.Payload.Metadata.GUID); nil == err && nil != thumb {
			event.Thumb = *thumb
		}
	}

	return
}
