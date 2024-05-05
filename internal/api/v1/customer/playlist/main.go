package playlist

import (
	"context"
	"encoding/json"
	"fmt"
	"music-master/internal/model"

	httputil "music-master/internal/util/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create creates a new Playlist
func (s *Playlist) Create(ctx context.Context, authUsr *model.AuthUser, data CreationData) (*model.Playlist, error) {
	rec := &model.Playlist{}

	s.converter.ToModel(rec, data)
	result, err := s.playlistCollection.InsertOne(ctx, rec)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// View returns single Playlist
func (s *Playlist) View(ctx context.Context, authUsr *model.AuthUser, id string) (*model.Playlist, error) {
	rec := new(model.Playlist)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error when parse object id", err)
		return nil, err
	}

	rec, err = s.playlistCollection.FindOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	return rec, nil
}

// Search returns single Playlist
func (s *Playlist) Search(ctx context.Context, authUsr *model.AuthUser, lq *httputil.ListRequest) ([]*model.Playlist, error) {
	query := map[string]interface{}{}
	if err := json.Unmarshal([]byte(lq.Filter), &query); err != nil {
		return nil, err
	}

	queryStr := query["query"].(string)
	rec, err := s.playlistCollection.Search(ctx, queryStr, lq.Page, lq.Limit)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

// Update updates Playlist information
func (s *Playlist) Update(ctx context.Context, authUsr *model.AuthUser, id string, data UpdateData) (*model.Playlist, error) {
	// * do validation
	curr, err := s.View(ctx, authUsr, id)
	if err != nil || curr == nil {
		return nil, err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	rec := new(model.Playlist)
	s.converter.ToModel(&curr, &data)

	rec, err = s.playlistCollection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, curr)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

// Delete deletes a Playlist
func (s *Playlist) Delete(ctx context.Context, authUsr *model.AuthUser, id string) error {
	// * do validation
	curr, err := s.View(ctx, authUsr, id)
	if err != nil || curr == nil {
		return err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = s.playlistCollection.RemoveOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a musictrack in Playlist
func (s *Playlist) DeleteMusicTrack(ctx context.Context, authUsr *model.AuthUser, id string, data DeleteMusicTrack) error {
	// * do validation
	curr, err := s.View(ctx, authUsr, id)
	if err != nil || curr == nil {
		fmt.Println("Error when view Play list", err)
		return err
	}

	err = s.playlistCollection.DeleteTrackFromPlaylist(ctx, id, data.MusicTrackID)
	if err != nil {
		return err
	}

	return nil
}
