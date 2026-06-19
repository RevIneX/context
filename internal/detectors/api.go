package detectors

import (
	"strings"

	"github.com/RevIneX/context/internal/data"
)

const (
	apiGraphQL   = "GraphQL"
	apiGRPC      = "gRPC"
	apiSOAP      = "SOAP"
	apiWebSocket = "WebSocket"
)

var apiSignatures = map[string]string{
	".graphql": apiGraphQL,
	".gql":     apiGraphQL,
	".proto":   apiGRPC,
	".wsdl":    apiSOAP,
}

var apiFiles = map[string]string{
	"websocket.html": apiWebSocket,
	"schema.graphql": apiGraphQL,
	"service.proto":  apiGRPC,
	"service.wsdl":   apiSOAP,
}

type APIDetector struct{}

func (d *APIDetector) Name() string {
	return "api"
}

func (d *APIDetector) Detect(filePath string) []data.Finding {
	if filePath == "" {
		return nil
	}

	lastSlash := strings.LastIndexByte(filePath, '/')
	lastBackslash := strings.LastIndexByte(filePath, '\\')
	if lastBackslash > lastSlash {
		lastSlash = lastBackslash
	}
	baseName := filePath[lastSlash+1:]

	if baseName == "" {
		return nil
	}

	if tech, ok := apiFiles[baseName]; ok {
		return []data.Finding{data.NewFinding(data.CatAPI, tech, "", filePath, 1, nil)}
	}

	lowName := strings.ToLower(baseName)
	if tech, ok := apiFiles[lowName]; ok {
		return []data.Finding{data.NewFinding(data.CatAPI, tech, "", filePath, 1, nil)}
	}

	dot := strings.LastIndexByte(lowName, '.')
	if dot != -1 {
		ext := lowName[dot:]
		if tech, ok := apiSignatures[ext]; ok {
			return []data.Finding{data.NewFinding(data.CatAPI, tech, "", filePath, 1, nil)}
		}
	}

	return nil
}