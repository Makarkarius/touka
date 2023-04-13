package comtradeapi

import (
	"proj/internal/app/server"
)

type Config struct {
	Token  string `json:"token"`
	ApiUrl string `json:"apiUrl"`

	SplitFactor       int `json:"splitFactor"`
	RequestTimeoutSec int `json:"requestTimeoutSec"`
}

func (c Config) Build() (server.Requester, error) {
	return NewComtradeRequester(c)
}
