package handler_test

import (
	"encoding/json"
	"io"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/artalkjs/artalk/v2/docs/swagger"
	"github.com/stretchr/testify/require"
	"github.com/swaggo/swag"
)

func TestAPISchemaMatchesSource(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	root := filepath.Join(filepath.Dir(filename), "../..")
	parser := swag.New(swag.SetDebugger(log.New(io.Discard, "", 0)))
	parser.RequiredByDefault = true
	require.NoError(t, parser.ParseAPI(root, "server/server.go", 100))

	var published struct {
		Definitions json.RawMessage `json:"definitions"`
	}
	require.NoError(t, json.Unmarshal([]byte(swagger.SwaggerInfo.ReadDoc()), &published))
	generated, err := json.Marshal(parser.GetSwagger().Definitions)
	require.NoError(t, err)

	// Keep generated API definitions in sync with the Go response types.
	require.JSONEq(t, string(published.Definitions), string(generated))
}
