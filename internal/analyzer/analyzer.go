package analyzer

import (
	"os"
	"path/filepath"
	"strings"
)

var ignoredDirsMap = map[string]struct{}{
	".git": {}, "node_modules": {}, "vendor": {}, "__pycache__": {},
	".idea": {}, ".vscode": {}, "dist": {}, "build": {},
	".next": {}, ".nuxt": {}, "target": {},
}

var allowedExtsMap = map[string]struct{}{
	".py": {}, ".go": {}, ".java": {}, ".js": {}, ".ts": {}, ".jsx": {}, ".tsx": {},
	".rb": {}, ".php": {}, ".c": {}, ".cpp": {}, ".h": {}, ".hpp": {}, ".cs": {},
	".rs": {}, ".swift": {}, ".kt": {}, ".scala": {}, ".yaml": {}, ".yml": {},
	".toml": {}, ".json": {}, ".xml": {}, ".ini": {}, ".cfg": {}, ".conf": {},
	".tf": {}, ".hcl": {}, ".sql": {}, ".sh": {}, ".bash": {}, ".ps1": {},
	".css": {}, ".scss": {}, ".sass": {}, ".less": {}, ".html": {}, ".htm": {},
	".graphql": {}, ".gql": {}, ".proto": {}, ".wsdl": {},
	".md": {}, ".txt": {}, ".properties": {},
	".cbl": {}, ".cob": {}, ".f90": {}, ".f95": {}, ".f03": {},
	".pas": {}, ".pp": {}, ".ada": {}, ".lisp": {}, ".lsp": {},
	".ml": {}, ".mli": {}, ".groovy": {}, ".gvy": {},
	".jl": {}, ".nim": {}, ".zig": {}, ".v": {}, ".vhdl": {},
	".mat": {}, ".m": {}, ".r": {}, ".rmd": {}, ".ipynb": {},
	".ex": {}, ".exs": {}, ".erl": {}, ".hrl": {},
	".hs": {}, ".lhs": {}, ".clj": {}, ".cljs": {}, ".elm": {},
	".vue": {}, ".svelte": {}, ".astro": {},

	"dockerfile": {}, "containerfile": {},
	"docker-compose.yml": {}, "docker-compose.yaml": {},
	"makefile": {}, "gnumakefile": {},
	".env": {}, ".gitlab-ci.yml": {}, "jenkinsfile": {}, "vagrantfile": {},
	"cmakelists.txt": {}, "rebar.config": {}, "mix.exs": {},
	"build.sbt": {}, "pubspec.yaml": {}, "chart.yaml": {}, "site.yml": {},
	"pulumi.yaml": {}, "cloudformation.json": {}, "terraform.tfvars": {},
	"schema.prisma": {}, "caddyfile": {}, "traefik.yml": {}, "traefik.yaml": {},
	"argo-app.yaml": {}, "spinnaker.yml": {}, "fluxcd.yaml": {},
	"sonar-project.properties": {}, ".eslintrc.json": {},
	".snyk": {}, ".trivyignore": {},
	"tsconfig.json": {}, "pyproject.toml": {}, "pipfile": {}, "poetry.lock": {},
	"gemfile": {}, "gemfile.lock": {}, "cargo.toml": {}, "cargo.lock": {},
	"go.mod": {}, "go.sum": {}, "package.json": {}, "package-lock.json": {},
	"composer.json": {}, "composer.lock": {}, "requirements.txt": {},
	"pom.xml": {}, "build.gradle": {}, "app.csproj": {}, "podfile": {},
	"notebook": {},
}

func WalkProject(root string) ([]string, error) {
	files := make([]string, 0, 1024)

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		name := d.Name()

		if d.IsDir() {
			if _, ignored := ignoredDirsMap[strings.ToLower(name)]; ignored {
				return filepath.SkipDir
			}
			return nil
		}

		lowName := strings.ToLower(name)

		if _, allowed := allowedExtsMap[lowName]; allowed {
			files = append(files, path)
			return nil
		}

		if _, allowed := allowedExtsMap[filepath.Ext(lowName)]; allowed {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
