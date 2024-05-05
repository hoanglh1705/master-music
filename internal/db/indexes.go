package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func (d *Database) CreateIndexes() {
	ctx := context.Background()
	d.createMusicTrackIndexes(ctx)
}

func (d *Database) createMusicTrackIndexes(ctx context.Context) {
	mods := []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{Key: "title", Value: bsonx.Int32(9)}},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bsonx.Doc{{Key: "artist", Value: bsonx.Int32(8)}},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bsonx.Doc{{Key: "album", Value: bsonx.Int32(6)}},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bsonx.Doc{{Key: "genre", Value: bsonx.Int32(3)}},
			Options: options.Index().SetUnique(false),
		},
	}

	if _, err := d.musicTrack.Indexes().CreateMany(ctx, mods); err != nil {
		fmt.Println("createMusicTrackIndexes().CreateMany() ERROR:", err)
	}
}
