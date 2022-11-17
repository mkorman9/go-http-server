package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mkorman9/tiny"
	"github.com/mkorman9/tiny/tinyhttp"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
)

//go:embed web
var webFS embed.FS

func main() {
	configFilePath := flag.String("config", "./config.yml", "path to config.yml file")
	flag.Parse()

	config, err := tiny.LoadConfig(*configFilePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to load configuration file: %v\n", err)
		os.Exit(1)
	}

	tiny.SetupLogger()

	httpServer := tinyhttp.NewServer(
		tinyhttp.Address(config.String("http.address", "0.0.0.0:8080")),
	)

	httpWebFS := http.FS(webFS)
	htmlTemplates, err := template.ParseFS(webFS, "web/templates/**/*.html")
	if err != nil {
		log.Error().Err(err).Msg("Failed to load HTML templates")
		return
	}
	httpServer.SetHTMLTemplate(htmlTemplates)

	httpServer.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "html/index.html", gin.H{
			"message": "Hello world",
		})
	})

	httpServer.GET("/static/*path", func(c *gin.Context) {
		p := c.Param("path")
		if strings.Contains(p, "..") {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.FileFromFS(path.Join("web/static/", p), httpWebFS)
	})

	httpServer.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("web/favicon.ico", httpWebFS)
	})

	tiny.StartAndBlock(httpServer)
}
