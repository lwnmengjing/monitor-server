package config

import (
	"embed"
	"log/slog"
	"os"

	"github.com/mss-boot-io/mss-boot/pkg/config"
	"github.com/mss-boot-io/mss-boot/pkg/config/source"
	"github.com/mss-boot-io/mss-boot/pkg/config/storage"
)

/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2023/10/31 16:37:31
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2023/10/31 16:37:31
 */

//go:embed *.yml
var FS embed.FS

var Cfg Config

type Config struct {
	Server config.Listen `yaml:"server" json:"server"`
	Logger config.Logger `yaml:"logger" json:"logger"`
	Queue  config.Queue  `yaml:"queue" json:"queue"`
	mq     storage.AdapterQueue
}

func (e *Config) Init() {
	opts := []source.Option{
		source.WithProvider(source.FS),
		source.WithFrom(FS),
	}
	err := config.Init(e, opts...)
	if err != nil {
		slog.Error("Config init failed", slog.Any("err", err))
		os.Exit(-1)
	}

	e.Queue.Init(func(queue storage.AdapterQueue) {
		e.mq = queue
	})

	e.Logger.Init()
}

func (e *Config) OnChange() {
	e.Logger.Init()
	e.Queue.Init(func(queue storage.AdapterQueue) {
		e.mq = queue
	})

	slog.Info("!!! cfg change and reload")
}

func (e *Config) GetMQ() storage.AdapterQueue {
	return e.mq
}
