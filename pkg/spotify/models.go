package spotify

import (
	"github.com/go-pg/pg/v10/orm"
	"time"
)

func init() {
	// This is necessary for many2many relationships starting in go-pg v10.
	orm.RegisterTable((*TrackArtists)(nil))
}

type PlaylistSnapshot struct {
	ID         string
	PlaylistID string
	SnapshotID string
	ImageURL   string
	CreatedAt  *time.Time
}

type PlaylistSnapshotDetail struct {
	PlaylistSnapshotID string
	Rank               int `pg:",use_zero"`
	TrackID            string
	Track              *Track `pg:"rel:has-one"`
}

type Track struct {
	ID          string
	Name        string
	Artists     []*Artist `pg:"many2many:track_artists"`
	AlbumID     string
	Album       *Album `pg:"rel:has-one"`
	DiscNumber  int
	TrackNumber int
	Explicit    bool
	DurationMS  int
}

type Artist struct {
	ID   string
	Name string
}

type Album struct {
	ID     string
	Name   string
	Images []*AlbumImage `pg:"rel:has-many"`
}

type AlbumImage struct {
	AlbumID string
	Height  int
	Width   int
	URL     string
}

type TrackArtists struct {
	TrackID  string
	Track    *Track `pg:"rel:has-one"`
	ArtistID string
	Artist   *Artist `pg:"rel:has-one"`
}
