package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mkorman9/go-commons/configutil"
	"github.com/mkorman9/go-commons/coreutil"
	"github.com/mkorman9/go-commons/httputil"
	"github.com/mkorman9/go-commons/logutil"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"os"
	"path"
)

var (
	//go:embed web/templates
	templatesFS embed.FS

	//go:embed web/static
	staticFS embed.FS
)

func main() {
	configFilePath := flag.String("config", "./config.yml", "path to config.yml file")
	flag.Parse()

	config, err := configutil.LoadConfigFromFile(*configFilePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to load configuration file: %v\n", err)
		os.Exit(1)
	}

	logutil.SetupLogger()

	httpServer := httputil.NewServer(
		httputil.Address(config.String("http.address", "0.0.0.0:8080")),
	)

	htmlTemplates, err := template.ParseFS(templatesFS, "web/templates/**/*.html")
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

	httpServer.GET("/static/*filepath", func(c *gin.Context) {
		c.FileFromFS(path.Join("web/", c.Request.URL.Path), http.FS(staticFS))
	})

	coreutil.StartAndBlock(httpServer)
}
