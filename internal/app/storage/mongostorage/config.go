package mongostorage

import (
	"proj/internal/app/server"
)

type Config struct {
	MongoURI string `json:"mongoURI"`
	DBName   string `json:"DBName"`
}

func (c Config) Build() (server.Storager, error) {
	return NewDatabase(c)
}
