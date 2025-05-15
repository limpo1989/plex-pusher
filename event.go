package main

type Account struct {
	ID    int    `json:"id"`
	Thumb string `json:"thumb"`
	Title string `json:"title"`
}

type Server struct {
	Title string `json:"title"`
	UUID  string `json:"uuid"`
}

type Player struct {
	Local         bool   `json:"local"`
	PublicAddress string `json:"PublicAddress"`
	Title         string `json:"title"`
	UUID          string `json:"uuid"`
}

type Metadata struct {
	LibrarySectionType   string `json:"librarySectionType"`
	RatingKey            string `json:"ratingKey"`
	Key                  string `json:"key"`
	ParentRatingKey      string `json:"parentRatingKey"`
	GrandparentRatingKey string `json:"grandparentRatingKey"`
	GUID                 string `json:"guid"`
	LibrarySectionID     int    `json:"librarySectionID"`
	MediaType            string `json:"type"`
	Title                string `json:"title"`
	GrandparentKey       string `json:"grandparentKey"`
	ParentKey            string `json:"parentKey"`
	GrandparentTitle     string `json:"grandparentTitle"`
	ParentTitle          string `json:"parentTitle"`
	Summary              string `json:"summary"`
	Index                int    `json:"index"`
	ParentIndex          int    `json:"parentIndex"`
	RatingCount          int    `json:"ratingCount"`
	Thumb                string `json:"thumb"`
	Art                  string `json:"art"`
	ParentThumb          string `json:"parentThumb"`
	GrandparentThumb     string `json:"grandparentThumb"`
	GrandparentArt       string `json:"grandparentArt"`
	AddedAt              int    `json:"addedAt"`
	UpdatedAt            int    `json:"updatedAt"`
}

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

type PlexEvent struct {
	Event    string   `json:"event"`
	User     bool     `json:"user"`
	Owner    bool     `json:"owner"`
	Account  Account  `json:"Account"`
	Server   Server   `json:"Server"`
	Player   Player   `json:"Player"`
	Metadata Metadata `json:"Metadata"`
}
