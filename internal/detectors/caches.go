package detectors

import (
	"path/filepath"
	"strings"

	"github.com/RevIneX/context/internal/data"
)

var cacheSignatures = []struct {
	Sig  string
	Tech string
}{
	{"memcached", "Memcached"},
	{"memcache", "Memcached"},
	{"redis-cluster", "Redis"},
	{"redis-sentinel", "Redis"},
	{"redis", "Redis"},
	{"varnish", "Varnish"},
	{"hazelcast", "Hazelcast"},
	{"ehcache", "Ehcache"},
	{"caffeine", "Caffeine"},
	{"infinispan", "Infinispan"},
}

var cacheFiles = map[string]string{
	"redis.conf":     "Redis",
	"memcached.conf": "Memcached",
	"varnish.vcl":    "Varnish",
}

type CacheDetector struct{}

func (d *CacheDetector) Name() string {
	return "cache"
}

func (d *CacheDetector) Detect(filePath string) []data.Finding {
	ext := filepath.Ext(filePath)
	// Zero-allocation проверка расширения через EqualFold
	switch {
	case strings.EqualFold(ext, ".yaml"),
		strings.EqualFold(ext, ".yml"),
		strings.EqualFold(ext, ".conf"),
		strings.EqualFold(ext, ".cfg"),
		strings.EqualFold(ext, ".xml"),
		strings.EqualFold(ext, ".json"),
		strings.EqualFold(ext, ".properties"),
		strings.EqualFold(ext, ".vcl"):
		// ok
	default:
		return nil
	}

	fileName := filepath.Base(filePath)
	lowName := strings.ToLower(fileName)

	if tech, ok := cacheFiles[lowName]; ok {
		return []data.Finding{
			data.NewFinding(data.CatCache, tech, "", filePath, 1, nil),
		}
	}

	for i := range cacheSignatures {
		if strings.Contains(lowName, cacheSignatures[i].Sig) {
			return []data.Finding{
				data.NewFinding(data.CatCache, cacheSignatures[i].Tech, "", filePath, 1, nil),
			}
		}
	}

	return nil
}
