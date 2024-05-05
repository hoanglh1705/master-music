package db

import (
	"context"
	"errors"
	"fmt"
	"music-master/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MusicTrackCollection struct {
	db *Database
}

func NewMusicTrackCollection(db *Database) *MusicTrackCollection {
	return &MusicTrackCollection{
		db: db,
	}
}

func (c *MusicTrackCollection) InsertOne(ctx context.Context, data *model.MusicTrack) (*model.MusicTrack, error) {
	result, err := c.db.musicTrack.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("invalid objectId")
	}
	data.ID = objectID

	return data, nil
}

func (c *MusicTrackCollection) FindOne(ctx context.Context, where bson.M) (*model.MusicTrack, error) {
	result := &model.MusicTrack{}
	opts := options.FindOne()
	opts.SetSort(bson.M{"_id": 1})
	if err := c.db.musicTrack.FindOne(ctx, where, opts).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *MusicTrackCollection) RemoveOne(ctx context.Context, where bson.M) error {
	if _, err := c.db.musicTrack.DeleteOne(ctx, where); err != nil {
		return err
	}

	return nil
}

func (c *MusicTrackCollection) UpdateOne(ctx context.Context, where bson.M, updateData bson.M) (*model.MusicTrack, error) {
	result := &model.MusicTrack{}
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	if err := c.db.musicTrack.FindOneAndUpdate(ctx, where, bson.M{"$set": updateData}, opts).Decode(result); err != nil {
		return nil, err
	}
	return result, nil

}

func (c *MusicTrackCollection) FindOneAndUpdate(ctx context.Context, where bson.M, data *model.MusicTrack) (*model.MusicTrack, error) {
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	result := &model.MusicTrack{}
	fmt.Println("FindOneAndUpdate", *data)
	if err := c.db.musicTrack.FindOneAndUpdate(ctx, where, bson.M{"$set": data}, opts).Decode(result); err != nil {
		return result, err
	}

	return result, nil
}

func (c *MusicTrackCollection) Search(ctx context.Context, searchQuery string, page, pageSize int) ([]*model.MusicTrack, error) {
	// Search for music tracks
	songFilter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"artist": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"album": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"genre": bson.M{"$regex": searchQuery, "$options": "i"}},
		},
	}

	// Paging
	myOptions := options.Find()
	if page > 0 && pageSize > 0 {
		skip := (page - 1) * pageSize
		myOptions.SetSkip(int64(skip)).SetLimit(int64(pageSize))
	}

	dataCursor, err := c.db.musicTrack.Find(context.Background(), songFilter, myOptions)
	if err != nil {
		fmt.Println("Error searching for music tracks:", err)
		return nil, err
	}
	defer dataCursor.Close(context.Background())

	// Process the matched music tracks
	var foundData []*model.MusicTrack
	if err := dataCursor.All(context.Background(), &foundData); err != nil {
		fmt.Println("Error decoding music tracks:", err)
		return nil, err
	}

	return foundData, nil
}
