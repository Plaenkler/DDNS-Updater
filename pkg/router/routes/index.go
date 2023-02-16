package routes

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	//go:embed static
	static embed.FS
)

func ProvideIndex(writer http.ResponseWriter, request *http.Request) {
	// Default to 404
	if request.URL.Path != "/" {
		fmt.Fprintf(writer, "[provide-index-1] 404 - Not found")
		return
	}
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusForbidden)
		_, err := writer.Write([]byte("[provide-index-2] 405 - Method not allowed"))
		if err != nil {
			log.Printf("[provide-index-3] could not write http reply - error: %s", err)
		}
		return
	}
	template, err := template.New("index").ParseFS(static,
		"static/html/index.html",
	)
	if err != nil {
		fmt.Fprintf(writer, "[provide-index-4] could not provide template - error: %s", err)
		return
	}
	writer.Header().Add("Content-Type", "text/html")
	err = template.Execute(writer, nil)
	if err != nil {
		fmt.Fprintf(writer, "[provide-index-5] could not execute parsed template - error: %v", err)
	}
}
