package detectors

import (
	"strings"

	"github.com/RevIneX/context/internal/data"
)

var serverNames = []string{
	"server.js", "server.ts", "server.mjs",
	"app.py", "api.py", "main.py",
	"main.go", "server.go",
	"application.java", "app.java",
	"app.js", "index.js", "api.js", "api.ts",
}

type ArchitectureDetector struct{}

func (d *ArchitectureDetector) Name() string {
	return "architecture"
}

func (d *ArchitectureDetector) Detect(filePath string) []data.Finding {
	var findings []data.Finding

	lastSlash := strings.LastIndexByte(filePath, '/')
	lastBackslash := strings.LastIndexByte(filePath, '\\')
	if lastBackslash > lastSlash {
		lastSlash = lastBackslash
	}
	fileName := filePath[lastSlash+1:]

	// docker-compose.yml → микросервисы (где бы ни находился)
	if strings.EqualFold(fileName, "docker-compose.yml") || strings.EqualFold(fileName, "docker-compose.yaml") {
		findings = append(findings, data.NewFinding(data.CatArchitecture, "микросервисы", "", filePath, 1, nil))
	}

	// index.html → клиент
	if strings.EqualFold(fileName, "index.html") || strings.EqualFold(fileName, "index.htm") {
		findings = append(findings, data.NewFinding(data.CatClient, "веб приложение", "", filePath, 1, nil))
	}

	// Серверные файлы
	for _, s := range serverNames {
		if strings.EqualFold(fileName, s) {
			findings = append(findings, data.NewFinding(data.CatServer, "api сервер", "", filePath, 1, nil))
			break
		}
	}

	return findings
}
