package server

import "go.uber.org/zap"

type Config struct {
	Host              string `json:"host"`
	Port              int    `json:"port"`
	RequestQueueSize  int    `json:"requestQueueSize"`
	ResponseQueueSize int    `json:"responseQueueSize"`
	ReadTimeoutSec    int    `json:"readTimeoutSec"`
	WriteTimeoutSec   int    `json:"writeTimeoutSec"`

	RequesterCfg RequesterConfiger `json:"requesterCfg"`
	StorageCfg   StorageConfiger   `json:"storageCfg"`

	LoggerCfg zap.Config `json:"loggerCfg"`

	RabbitURI               string `json:"rabbitURI"`
	RabbitPublishTimeoutSec int    `json:"rabbitPublishTimeoutSec"`
	RabbitExchangeName      string `json:"rabbitExchangeName"`
	RabbitKey               string `json:"rabbitKey"`
}

func (c Config) Build() (*server, error) {
	return NewServer(c)
}
