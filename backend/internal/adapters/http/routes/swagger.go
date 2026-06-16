package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterSwagger registers simple Swagger UI and JSON endpoints.
// - /swagger serves a small HTML page that loads Swagger UI from CDN
// - /swagger.json serves the OpenAPI JSON file from ./docs/swagger.json
func RegisterSwagger(router *mux.Router) {
	router.Path("/swagger.json").Handler(http.FileServer(http.Dir("./docs"))).Methods(http.MethodGet)

	router.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Swagger UI</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui-bundle.js"></script>
    <script>
      window.onload = function() {
        const ui = SwaggerUIBundle({
          url: '/swagger.json',
          dom_id: '#swagger-ui',
        })
        window.ui = ui
      }
    </script>
  </body>
</html>`)
	}).Methods(http.MethodGet)
}
