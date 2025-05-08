package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mss-boot-io/mss-boot/core/server"
	"github.com/mss-boot-io/mss-boot/core/server/listener"
	"github.com/mss-boot-io/mss-boot/virtual/action"
	"github.com/spf13/cobra"
	"service-http/config"
	"service-http/models"

	"service-http/router"
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
		Long:    "start service-http server",
		Example: "service-http server",
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
			listener.WithName("service-http"),
			listener.WithHandler(r)),
	}

	// init virtual models
	ms, err := models.GetModels()
	if err != nil {
		return err
	}
	for i := range ms {
		action.SetModel(ms[i].Path, ms[i].MakeVirtualModel())
	}
	server.Manage.Add(runnable...)

	return nil
}

func run() error {
	ctx := context.Background()

	return server.Manage.Start(ctx)
}
