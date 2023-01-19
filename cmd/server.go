package cmd

import (
	"fmt"
	"io"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/server"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"serve"},
	Short:   "Start the server",
	Long:    Banner,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		core.LoadCore(cfgFile, workDir)

		fmt.Println(Banner)
		fmt.Print("-------------------------------\n\n")

		// create fiber app
		app := fiber.New(fiber.Config{
			// @see https://github.com/gofiber/fiber/issues/426
			// @see https://github.com/gofiber/fiber/issues/185
			Immutable: true,
		})

		// logger
		app.Use(logger.New(logger.Config{
			Format: "[${status}] ${method} ${path} ${latency} ${ip} ${reqHeader:X-Request-ID} ${referer} ${ua}\n",
			Output: io.Discard,
			Done: func(c *fiber.Ctx, logString []byte) {
				statusOK := c.Response().StatusCode() >= 200 && c.Response().StatusCode() <= 302
				if !statusOK {
					logrus.StandardLogger().WriterLevel(logrus.ErrorLevel).Write(logString)
				} else {
					logrus.StandardLogger().WriterLevel(logrus.DebugLevel).Write(logString)
				}
			},
		}))

		// init router
		server.Init(app)

		// listen
		listenAddr := fmt.Sprintf("%s:%d", config.Instance.Host, config.Instance.Port)
		if config.Instance.SSL.Enabled {
			app.ListenTLS(listenAddr, config.Instance.SSL.CertPath, config.Instance.SSL.KeyPath)
		} else {
			app.Listen(listenAddr)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	flagPV(serverCmd, "host", "", "0.0.0.0", "Listening IP")
	flagPV(serverCmd, "port", "", 23366, "Listening port")
}
