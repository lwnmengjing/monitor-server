package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/mss-boot-io/mss-boot/pkg"
	"github.com/spf13/cobra"

	"monitor-server/cmd/migrate"
	"monitor-server/cmd/server"
)

/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2023/10/31 16:37:31
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2023/10/31 16:37:31
 */

var rootCmd = &cobra.Command{
	Use:          "monitor-server",
	Short:        "monitor-server",
	SilenceUsage: true,
	Long:         `monitor-server is a background management system developed by the mss-boot framework`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New(pkg.Red("requires at least one arg"))
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	usageStr := `欢迎使用 ` + pkg.Green(`monitor-server 0.0.1`) + ` 可以使用 ` + pkg.Red(`-h`) + ` 查看命令`
	usageStr1 := `也可以参考 https://doc.mss-boot-io.top 的相关内容`
	fmt.Printf("%s\n", usageStr)
	fmt.Printf("%s\n", usageStr1)
}

func init() {
	rootCmd.AddCommand(server.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
}

// Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
