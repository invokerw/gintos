package assets

import (
	"embed"
)

//go:embed frontend
var Frontend embed.FS

//go:embed openapi.yaml
var OpenApiData []byte
