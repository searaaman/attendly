// Package docs serves the hand-written OpenAPI spec and a Swagger UI page for it.
package docs

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed openapi.yaml
var spec []byte

const uiPage = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Attendly API</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>SwaggerUIBundle({ url: "docs/openapi.yaml", dom_id: "#swagger-ui" });</script>
</body>
</html>`

func SpecHandler(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "application/yaml", spec)
}

func UIHandler(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(uiPage))
}
