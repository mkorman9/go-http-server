package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/mkorman9/go-commons/coreutil"
	"github.com/mkorman9/go-commons/httputil"
	"github.com/mkorman9/go-commons/logutil"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
)

//go:embed templates
var templatesFS embed.FS

func main() {
	logutil.SetupLogger()

	httpServer := httputil.NewServer()

	htmlTemplates, err := template.ParseFS(templatesFS, "templates/**/*.html")
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

	coreutil.StartAndBlock(httpServer)
}
