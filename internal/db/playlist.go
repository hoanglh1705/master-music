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

type PlaylistCollection struct {
	db *Database
}

func NewPlaylistCollection(db *Database) *PlaylistCollection {
	return &PlaylistCollection{
		db: db,
	}
}

func (c *PlaylistCollection) InsertOne(ctx context.Context, data *model.Playlist) (*model.Playlist, error) {
	result, err := c.db.playlist.InsertOne(ctx, data)
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

func (c *PlaylistCollection) FindOne(ctx context.Context, where bson.M) (*model.Playlist, error) {
	result := &model.Playlist{}
	opts := options.FindOne()
	opts.SetSort(bson.M{"_id": 1})
	if err := c.db.playlist.FindOne(ctx, where, opts).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *PlaylistCollection) RemoveOne(ctx context.Context, where bson.M) error {
	if _, err := c.db.playlist.DeleteOne(ctx, where); err != nil {
		return err
	}

	return nil
}

func (c *PlaylistCollection) UpdateOne(ctx context.Context, where bson.M, updateData bson.M) (*model.Playlist, error) {
	result := &model.Playlist{}
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	if err := c.db.playlist.FindOneAndUpdate(ctx, where, bson.M{"$set": updateData}, opts).Decode(result); err != nil {
		return nil, err
	}
	return result, nil

}

func (c *PlaylistCollection) FindOneAndUpdate(ctx context.Context, where bson.M, data *model.Playlist) (*model.Playlist, error) {
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	result := &model.Playlist{}
	fmt.Println("FindOneAndUpdate", *data)
	if err := c.db.playlist.FindOneAndUpdate(ctx, where, bson.M{"$set": data}, opts).Decode(result); err != nil {
		return result, err
	}

	return result, nil
}

// DeleteTrackFromPlaylist removes a track from a playlist in MongoDB
func (c *PlaylistCollection) DeleteTrackFromPlaylist(ctx context.Context, playlistID, trackID string) error {
	// Define filter to identify the playlist
	playlistObjectID, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		fmt.Println("Error parse playlist object id", err)
		return err
	}
	filter := bson.M{"_id": playlistObjectID}

	// Define the update to remove the track from the playlist
	musicTrackObjectID, err := primitive.ObjectIDFromHex(trackID)
	if err != nil {
		fmt.Println("Error parse music track object id", err)
		return err
	}

	update := bson.M{"$pull": bson.M{"tracks": bson.M{"_id": musicTrackObjectID}}}
	// Perform the update operation
	_, err = c.db.playlist.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("error deleting track from playlist: %v", err)
	}

	return nil
}

func (c *PlaylistCollection) Search(ctx context.Context, searchQuery string, page, pageSize int) ([]*model.Playlist, error) {
	// Define the filter for full-text search
	// Define the filter to search within the Tracks field
	filter := bson.M{
		"tracks": bson.M{
			"$elemMatch": bson.M{
				"$or": []bson.M{
					{"title": bson.M{"$regex": searchQuery, "$options": "i"}},
					{"artist": bson.M{"$regex": searchQuery, "$options": "i"}},
					{"album": bson.M{"$regex": searchQuery, "$options": "i"}},
					{"genre": bson.M{"$regex": searchQuery, "$options": "i"}},
				},
			},
		},
	}

	// Paging
	myOptions := options.Find()
	if page > 0 && pageSize > 0 {
		skip := (page - 1) * pageSize
		myOptions.SetSkip(int64(skip)).SetLimit(int64(pageSize))
	}

	// Perform the search
	cursor, err := c.db.playlist.Find(ctx, filter, myOptions)
	if err != nil {
		fmt.Println("Error searching for playlists:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// Process the matched playlists
	var foundData []*model.Playlist
	if err := cursor.All(ctx, &foundData); err != nil {
		fmt.Println("Error decoding playlists:", err)
		return nil, err
	}

	for _, playlist := range foundData {
		fmt.Println(playlist.Name)
	}

	return foundData, nil
}
