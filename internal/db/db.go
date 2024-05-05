package db

import (
	"context"
	"music-master/config"
	"music-master/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Database struct {
	client     *mongo.Client
	musicTrack *mongo.Collection
	playlist   *mongo.Collection
}

func New(cfg *config.Configuration) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(cfg.DBUrl)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	mongoDB := client.Database(cfg.DBName)

	db := &Database{
		client:     client,
		musicTrack: mongoDB.Collection(model.MusicTrack{}.TableName()),
		playlist:   mongoDB.Collection(model.Playlist{}.TableName()),
	}

	db.CreateIndexes()

	return db, nil
}

func (s *Database) Disconnect() {
	s.client.Disconnect(context.Background())
}

func (s *Database) ExecTx(ctx context.Context, fn func(sessionCtx mongo.SessionContext) error) error {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	session, err := s.client.StartSession()

	if err := mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(txnOpts); err != nil {
			return err
		}

		if err := fn(sessionContext); err != nil {
			return err
		}

		if err = session.CommitTransaction(sessionContext); err != nil {
			return err
		}

		return nil
	}); err != nil {
		if abortErr := session.AbortTransaction(ctx); abortErr != nil {
			return err
		}

		return err
	}

	return nil
}
