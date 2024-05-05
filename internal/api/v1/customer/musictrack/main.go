package musictrack

import (
	"context"
	"encoding/json"
	"music-master/internal/model"
	httputil "music-master/internal/util/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create creates a new MusicTrack account
func (s *MusicTrack) Create(ctx context.Context, authUsr *model.AuthUser, data CreationData) (*model.MusicTrack, error) {
	rec := &model.MusicTrack{}

	s.converter.ToModel(rec, data)
	result, err := s.musicTrackCollection.InsertOne(ctx, rec)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// View returns single MusicTrack
func (s *MusicTrack) View(ctx context.Context, authUsr *model.AuthUser, id string) (*model.MusicTrack, error) {
	rec := new(model.MusicTrack)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	rec, err = s.musicTrackCollection.FindOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	return rec, nil
}

// Search returns single MusicTrack
func (s *MusicTrack) Search(ctx context.Context, authUsr *model.AuthUser, lq *httputil.ListRequest) ([]*model.MusicTrack, error) {
	query := map[string]interface{}{}
	if err := json.Unmarshal([]byte(lq.Filter), &query); err != nil {
		return nil, err
	}

	queryStr := query["query"].(string)
	rec, err := s.musicTrackCollection.Search(ctx, queryStr, lq.Page, lq.Limit)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

// Update updates MusicTrack information
func (s *MusicTrack) Update(ctx context.Context, authUsr *model.AuthUser, id string, data UpdateData) (*model.MusicTrack, error) {
	// * do validation
	curr, err := s.View(ctx, authUsr, id)
	if err != nil || curr == nil {
		return nil, err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	rec := new(model.MusicTrack)
	s.converter.ToModel(&curr, &data)

	rec, err = s.musicTrackCollection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, curr)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

// Delete deletes a MusicTrack
func (s *MusicTrack) Delete(ctx context.Context, authUsr *model.AuthUser, id string) error {
	// * do validation
	curr, err := s.View(ctx, authUsr, id)
	if err != nil || curr == nil {
		return err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = s.musicTrackCollection.RemoveOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}
