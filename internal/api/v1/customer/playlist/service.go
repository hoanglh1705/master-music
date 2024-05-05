package playlist

import (
	"context"
	"music-master/internal/model"

	"go.mongodb.org/mongo-driver/bson"
)

// New creates new playlist application service
func New(PlaylistCollection PlaylistCollection, converter ModelConverter) *Playlist {
	return &Playlist{
		playlistCollection: PlaylistCollection,
		converter:          converter,
	}
}

// Playlist represents playlist application service
type Playlist struct {
	playlistCollection PlaylistCollection
	converter          ModelConverter
}

type PlaylistCollection interface {
	InsertOne(ctx context.Context, data *model.Playlist) (*model.Playlist, error)
	FindOne(ctx context.Context, where bson.M) (*model.Playlist, error)
	UpdateOne(ctx context.Context, where bson.M, updateData bson.M) (*model.Playlist, error)
	FindOneAndUpdate(ctx context.Context, where bson.M, data *model.Playlist) (*model.Playlist, error)
	RemoveOne(ctx context.Context, where bson.M) error
	DeleteTrackFromPlaylist(ctx context.Context, playlistID, trackID string) error
	Search(ctx context.Context, searchQuery string, page, pageSize int) ([]*model.Playlist, error)
}

type ModelConverter interface {
	FromModel(to interface{}, from interface{})
	ToModel(to interface{}, from interface{})
}
