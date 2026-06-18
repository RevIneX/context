package detectors

import (
	"strings"

	"github.com/RevIneX/context/internal/data"
)

var containerRules = []struct {
	Name string
	Tech string
}{
	{"dockerfile", "Docker"},
	{"containerfile", "Docker"},
	{".podman.yaml", "Podman"},
	{"buildah.sh", "Buildah"},
	{"docker-compose.yml", "Docker Compose"},
	{"docker-compose.yaml", "Docker Compose"},
}

type ContainerDetector struct{}

func (d *ContainerDetector) Name() string {
	return "container"
}

func (d *ContainerDetector) Detect(filePath string) []data.Finding {
	lastSlash := strings.LastIndexByte(filePath, '/')
	if lastSlash == -1 {
		lastSlash = strings.LastIndexByte(filePath, '\\')
	}
	fileName := filePath[lastSlash+1:]

	for _, rule := range containerRules {
		if strings.EqualFold(fileName, rule.Name) {
			return []data.Finding{
				data.NewFinding(data.CatContainer, rule.Tech, "", filePath, 1, nil),
			}
		}
	}

	if strings.HasSuffix(filePath, ".lxc/config") || strings.HasSuffix(filePath, ".lxc\\config") {
		return []data.Finding{
			data.NewFinding(data.CatContainer, "LXC", "", filePath, 1, nil),
		}
	}

	return []data.Finding{}
}