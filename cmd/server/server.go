package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/core/server/listener"
	"github.com/mss-boot-io/mss-boot/pkg/config/storage"
	"github.com/mss-boot-io/mss-boot/pkg/config/storage/queue"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"log/slog"
	"monitor-server/pkg"

	"monitor-server/config"
	"monitor-server/router"
)

/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2023/10/31 16:37:31
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2023/10/31 16:37:31
 */

var (
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "start server",
		Long:    "start monitor-server server",
		Example: "monitor-server server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func setup() error {
	// setup config
	config.Cfg.Init()

	r := gin.Default()
	router.Init(r.Group("/"))

	runnable := []server.Runnable{
		config.Cfg.Server.Init(
			listener.WithStartedHook(cronJob),
			listener.WithName("monitor-server"),
			listener.WithHandler(r)),
	}

	server.Manage.Add(runnable...)

	return nil
}

func run() error {
	ctx := context.Background()

	return server.Manage.Start(ctx)
}

func cronJob() {
	// cron job 每十分钟执行一次
	var lastNetIn, lastNetOut uint64
	task := cron.New(cron.WithSeconds())
	_, err := task.AddFunc("0 */10 * * * ?", func() {
		slog.Debug("cron job start")
		data, netIn, netOut, err := pkg.CollectMetrics(lastNetIn, lastNetOut, 600)
		if err != nil {
			slog.Error("cron job error", slog.Any("err", err))
			return
		}
		lastNetIn = netIn
		lastNetOut = netOut
		message := &queue.Message{}
		message.SetStream("server-monitor")
		message.SetValues(data)
		err = config.Cfg.GetMQ().Append(storage.WithTopic("server-monitor"), storage.WithMessage(message))
		if err != nil {
			slog.Error("cron job error", slog.Any("err", err))
			return
		}
	})
	if err != nil {
		slog.Error("cron job error", slog.Any("err", err))
		return
	}
	task.Start()
}
