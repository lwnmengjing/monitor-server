package config

import (
	"log/slog"
	"os"

	"github.com/mss-boot-io/mss-boot/pkg/config"
	"github.com/mss-boot-io/mss-boot/pkg/config/gormdb"
	"github.com/mss-boot-io/mss-boot/pkg/config/source"
)

/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2023/10/31 16:37:31
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2023/10/31 16:37:31
 */

var Cfg Config

type Config struct {
	Server   config.Listen    `yaml:"server" json:"server"`
	Database *gormdb.Database `yaml:"database" json:"database"`
	Logger   config.Logger    `yaml:"logger" json:"logger"`
}

func (e *Config) Init() {
	opts := []source.Option{
		source.WithDir("config"),
		source.WithProvider(source.Local),
	}
	err := config.Init(e, opts...)
	if err != nil {
		slog.Error("Config init failed", slog.Any("err", err))
		os.Exit(-1)
	}

	e.Logger.Init()
	e.Database.Init()
}

func (e *Config) OnChange() {
	e.Logger.Init()
	e.Database.Init()
	slog.Info("!!! cfg change and reload")
}
