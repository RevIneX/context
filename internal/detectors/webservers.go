package detectors

import (
	"strings"

	"github.com/RevIneX/context/internal/data"
)

// webserverRules — правила определения веб-серверов по имени файла.
// Проверка через EqualFold, zero-allocation.
var webserverRules = []struct {
	Name string
	Tech string
}{
	{"nginx.conf", "Nginx"},
	{"apache.conf", "Apache"},
	{"httpd.conf", "Apache"},
	{"caddyfile", "Caddy"},
	{"traefik.yml", "Traefik"},
	{"traefik.yaml", "Traefik"},
	{"lighttpd.conf", "Lighttpd"},
	{"haproxy.cfg", "HAProxy"},
	{"envoy.yaml", "Envoy"},
	{"envoy.yml", "Envoy"},
}

type WebserverDetector struct{}

func (d *WebserverDetector) Name() string {
	return "webserver"
}

func (d *WebserverDetector) Detect(filePath string) []data.Finding {
	// Извлекаем имя файла без аллокаций (срез строки)
	lastSlash := strings.LastIndexByte(filePath, '/')
	if lastSlash == -1 {
		lastSlash = strings.LastIndexByte(filePath, '\\')
	}
	fileName := filePath[lastSlash+1:]

	// Проверяем точное совпадение через EqualFold (без ToLower)
	for _, rule := range webserverRules {
		if strings.EqualFold(fileName, rule.Name) {
			return []data.Finding{
				data.NewFinding(data.CatWebServer, rule.Tech, "", filePath, 1, nil),
			}
		}
	}

	return []data.Finding{}
}
