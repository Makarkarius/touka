package mongostorage

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

type database struct {
	cfg Config

	client           *mongo.Client
	importCollection *mongo.Collection
	exportCollection *mongo.Collection

	logger *zap.Logger
}

func NewDatabase(cfg Config) (*database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}
	db := client.Database(cfg.DBName)

	impCol := db.Collection("import")
	expCol := db.Collection("export")

	return &database{
		cfg:              cfg,
		client:           client,
		importCollection: impCol,
		exportCollection: expCol,
		logger:           zap.NewNop(),
	}, nil
}

func (d *database) UseLogger(logger *zap.Logger) {
	d.logger = logger
}
