package detectors

import (
	"bufio"
	"os"
	"strings"

	"github.com/RevIneX/context/internal/data"
)

var frameworkSignatures = map[string]string{
	"django":    "Django",
	"flask":     "Flask",
	"fastapi":   "FastAPI",
	"aiohttp":   "AIOHTTP",
	"tornado":   "Tornado",
	"pyramid":   "Pyramid",
	"bottle":    "Bottle",
	"sanic":     "Sanic",
	"express":   "Express",
	"next":      "Next.js",
	"nuxt":      "Nuxt",
	"nest":      "NestJS",
	"fastify":   "Fastify",
	"koa":       "Koa",
	"hapi":      "Hapi",
	"svelte":    "Svelte",
	"react":     "React",
	"vue":       "Vue",
	"angular":   "Angular",
	"gatsby":    "Gatsby",
	"remix":     "Remix",
	"gin":       "Gin",
	"echo":      "Echo",
	"mux":       "Gorilla Mux",
	"fiber":     "Fiber",
	"iris":      "Iris",
	"spring-boot": "Spring Boot",
	"spring":    "Spring",
	"quarkus":   "Quarkus",
	"micronaut": "Micronaut",
	"jakarta":   "Jakarta EE",
	"rails":     "Rails",
	"sinatra":   "Sinatra",
	"hanami":    "Hanami",
	"laravel":   "Laravel",
	"symfony":   "Symfony",
	"slim":      "Slim",
	"lumen":     "Lumen",
	"actix-web": "Actix Web",
	"rocket":    "Rocket",
	"warp":      "Warp",
	"axum":      "Axum",
	"phoenix":   "Phoenix",
	"ktor":      "Ktor",
}

var dependencyFiles = map[string]bool{
	"package.json":      true,
	"package-lock.json": true,
	"yarn.lock":         true,
	"requirements.txt":  true,
	"pipfile":           true,
	"pipfile.lock":      true,
	"pyproject.toml":    true,
	"go.mod":            true,
	"go.sum":            true,
	"pom.xml":           true,
	"build.gradle":      true,
	"gemfile":           true,
	"gemfile.lock":      true,
	"composer.json":     true,
	"composer.lock":     true,
	"cargo.toml":        true,
	"cargo.lock":        true,
	"mix.exs":           true,
	"rebar.config":      true,
}

type FrameworkDetector struct{}

func (d *FrameworkDetector) Name() string {
	return "framework"
}

func (d *FrameworkDetector) Detect(filePath string) []data.Finding {
	lastSlash := strings.LastIndexByte(filePath, '/')
	if lastSlash == -1 {
		lastSlash = strings.LastIndexByte(filePath, '\\')
	}
	fileName := filePath[lastSlash+1:]
	lowName := strings.ToLower(fileName)

	if !dependencyFiles[lowName] {
		return []data.Finding{}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return []data.Finding{}
	}
	defer file.Close()

	var findings []data.Finding
	scanner := bufio.NewScanner(file)
	lineNum := 0
	maxLines := 1000

	for scanner.Scan() && lineNum < maxLines {
		lineNum++
		line := scanner.Text()
		lowLine := strings.ToLower(line)

		words := strings.FieldsFunc(lowLine, func(r rune) bool {
			return r == ' ' || r == '"' || r == ':' || r == ',' || r == '{' || r == '}' || r == '[' || r == ']' || r == '=' || r == '\'' || r == '>' || r == '<' || r == '/'
		})
		for _, word := range words {
			if tech, ok := frameworkSignatures[word]; ok {
				findings = append(findings, data.NewFinding(data.CatFramework, tech, "", filePath, lineNum, nil))
			}
		}
	}

	return findings
}
