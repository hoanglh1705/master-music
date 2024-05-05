package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// swagger:model MusicTrack
type MusicTrack struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title,omitempty" json:"title"`
	Artist      string             `bson:"artist,omitempty" json:"artist"`
	Album       string             `bson:"album,omitempty" json:"album"`
	Genre       string             `bson:"genre,omitempty" json:"genre"`
	ReleaseYear int                `bson:"release_year,omitempty" json:"release_year"`
	Duration    int                `bson:"duration,omitempty" json:"duration"` // Duration in seconds
	MP3File     []byte             `bson:"mp3_file,omitempty" json:"mp3_file"` // Binary data of the MP3 file
}

func (MusicTrack) TableName() string {
	return "music_tracks"
}
