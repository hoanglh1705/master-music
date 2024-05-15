package musictrack

import (
	"context"
	"music-master/internal/model"

	"go.mongodb.org/mongo-driver/bson"
)

// New creates new musictrack application service
func New(musicTrackCollection MusicTrackCollection,
	converter ModelConverter,
	musicTrackES MusicTrackES) *MusicTrack {
	return &MusicTrack{
		musicTrackCollection: musicTrackCollection,
		converter:            converter,
		musicTrackES:         musicTrackES,
	}
}

// MusicTrack represents musictrack application service
type MusicTrack struct {
	musicTrackCollection MusicTrackCollection
	converter            ModelConverter
	musicTrackES         MusicTrackES
}

type MusicTrackCollection interface {
	InsertOne(ctx context.Context, data *model.MusicTrack) (*model.MusicTrack, error)
	FindOne(ctx context.Context, where bson.M) (*model.MusicTrack, error)
	UpdateOne(ctx context.Context, where bson.M, updateData bson.M) (*model.MusicTrack, error)
	FindOneAndUpdate(ctx context.Context, where bson.M, data *model.MusicTrack) (*model.MusicTrack, error)
	RemoveOne(ctx context.Context, where bson.M) error
	Search(ctx context.Context, searchQuery string, page, pageSize int) ([]*model.MusicTrack, error)
}

type MusicTrackES interface {
	Search(ctx context.Context)
}

type ModelConverter interface {
	FromModel(to interface{}, from interface{})
	ToModel(to interface{}, from interface{})
}
