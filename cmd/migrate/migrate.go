package migrate

import (
	"log/slog"
	"path/filepath"

	"github.com/mss-boot-io/mss-boot/pkg/config/gormdb"
	"github.com/mss-boot-io/mss-boot/pkg/migration"
	"github.com/mss-boot-io/mss-boot/pkg/migration/models"
	"github.com/spf13/cobra"

	_ "service-http/cmd/migrate/migration/custom"
	_ "service-http/cmd/migrate/migration/system"
	"service-http/config"
)

/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2023/10/31 16:37:31
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2023/10/31 16:37:31
 */

var (
	generate bool
	username string
	password string
	system   bool
	StartCmd = &cobra.Command{
		Use:     "migrate",
		Short:   "Initialize the database",
		Example: "mss-boot-admin migrate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().BoolVarP(&system, "system", "s",
		false, "generate system migration file")
	StartCmd.PersistentFlags().BoolVarP(&generate, "generate", "g",
		false, "generate migration file")
	StartCmd.PersistentFlags().StringVarP(&username, "username", "u",
		"admin", "system super administrator login username")
	StartCmd.PersistentFlags().StringVarP(&password, "password", "p",
		"admin", "system super administrator login password")
}

func Run() error {
	if !generate {
		slog.Info("start init")
		config.Cfg.Init()
		return migrate()
	}
	slog.Info(`generate migration file`)
	return migration.GenFile(system, filepath.Join("cmd", "migrate", "migration"))
}

func migrate() error {
	db := gormdb.DB
	err := db.AutoMigrate(&models.Migration{})
	if err != nil {
		slog.Error("auto migrate error", slog.Any("err", err))
		return err
	}
	migration.Migrate.SetDb(db)
	migration.Migrate.Migrate()
	return err
}
