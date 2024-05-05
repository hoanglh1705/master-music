package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// swagger:model Playlist
type Playlist struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name   string             `bson:"name,omitempty" json:"name"`
	Tracks []*MusicTrack      `bson:"tracks,omitempty" json:"tracks"`
}

func (Playlist) TableName() string {
	return "playlists"
}
