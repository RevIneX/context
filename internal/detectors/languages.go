package detectors

import (
	"github.com/RevIneX/context/internal/data"
	"path/filepath"
	"strings"
)

var languageMap = map[string]string{
	".py": "Python", ".go": "Go", ".java": "Java", ".js": "JavaScript",
	".mjs": "JavaScript", ".cjs": "JavaScript", ".ts": "TypeScript",
	".jsx": "JavaScript", ".tsx": "TypeScript", ".rb": "Ruby", ".php": "PHP",
	".c": "C", ".h": "C", ".cpp": "C++", ".cc": "C++", ".cxx": "C++",
	".hpp": "C++", ".hh": "C++", ".hxx": "C++", ".cs": "C#", ".rs": "Rust",
	".swift": "Swift", ".kt": "Kotlin", ".kts": "Kotlin", ".scala": "Scala",
	".pl": "Perl", ".pm": "Perl", ".sh": "Shell", ".bash": "Shell",
	".lua": "Lua", ".dart": "Dart", ".ex": "Elixir", ".exs": "Elixir",
	".erl": "Erlang", ".hrl": "Erlang", ".hs": "Haskell", ".lhs": "Haskell",
	".clj": "Clojure", ".cljs": "ClojureScript", ".elm": "Elm", ".vue": "Vue",
	".svelte": "Svelte", ".astro": "Astro", ".tf": "HCL", ".hcl": "HCL",
	".yaml": "YAML", ".yml": "YAML", ".toml": "TOML", ".json": "JSON",
	".xml": "XML", ".ini": "INI", ".cfg": "INI", ".conf": "INI",
	".md": "Markdown", ".txt": "Text", ".css": "CSS", ".scss": "SCSS",
	".sass": "Sass", ".less": "Less", ".html": "HTML", ".htm": "HTML",
	".cbl": "COBOL", ".cob": "COBOL",
	".f90": "Fortran", ".f95": "Fortran", ".f03": "Fortran",
	".pas": "Pascal", ".pp": "Pascal",
	".ada": "Ada",
	".lisp": "Lisp", ".lsp": "Lisp",
	".ml": "OCaml", ".mli": "OCaml",
	".groovy": "Groovy", ".gvy": "Groovy",
	".jl": "Julia",
	".nim": "Nim",
	".zig": "Zig",
	".v": "Verilog", ".vhdl": "VHDL",
	".mat": "MATLAB", ".m": "MATLAB",
	".ipynb": "Jupyter",
	"dockerfile":            "Docker",
	"containerfile":         "Docker",
	"docker-compose.yml":    "Docker Compose",
	"docker-compose.yaml":   "Docker Compose",
	"makefile":              "Make",
	"gnumakefile":           "Make",
	".env":                  "Env",
	".gitlab-ci.yml":        "GitLab CI",
	"jenkinsfile":           "Jenkins",
	"vagrantfile":           "Vagrant",
	"cmakelists.txt":        "CMake",
	"rebar.config":          "Erlang",
	"mix.exs":               "Elixir",
	"build.sbt":             "Scala",
	"pubspec.yaml":          "Dart",
	"chart.yaml":            "Helm",
	"site.yml":              "Ansible",
	"pulumi.yaml":           "Pulumi",
	"cloudformation.json":   "CloudFormation",
	"terraform.tfvars":      "Terraform",
	"schema.prisma":         "Prisma",
	"caddyfile":             "Caddy",
	"traefik.yml":           "Traefik",
	"traefik.yaml":          "Traefik",
	"argo-app.yaml":         "ArgoCD",
	"spinnaker.yml":         "Spinnaker",
	"fluxcd.yaml":           "FluxCD",
	"sonar-project.properties": "SonarQube",
	".eslintrc.json":        "ESLint",
	".snyk":                 "Snyk",
	".trivyignore":          "Trivy",
	"tsconfig.json":         "TypeScript",
	"pyproject.toml":        "Python",
	"pipfile":               "Python",
	"poetry.lock":           "Python",
	"gemfile":               "Ruby",
	"gemfile.lock":          "Ruby",
	"cargo.toml":            "Rust",
	"cargo.lock":            "Rust",
	"go.mod":                "Go",
	"go.sum":                "Go",
	"package.json":          "Node.js",
	"package-lock.json":     "Node.js",
	"composer.json":         "PHP",
	"composer.lock":         "PHP",
	"requirements.txt":      "Python",
	"pom.xml":               "Java",
	"build.gradle":          "Java",
	"app.csproj":            "C#",
	"podfile":               "iOS",
	"notebook":              "Jupyter",
}

type LanguageDetector struct{}

func (d *LanguageDetector) Name() string {
	return "language"
}

func (d *LanguageDetector) Detect(filePath string) []data.Finding {
	ext := filepath.Ext(filePath)
	if ext != "" {
		if lang, ok := languageMap[strings.ToLower(ext)]; ok {
			return []data.Finding{
				data.NewFinding(data.CatLanguage, lang, "", filePath, 1, nil),
			}
		}
	}

	lastSlash := strings.LastIndexByte(filePath, '/')
	if lastSlash == -1 {
		lastSlash = strings.LastIndexByte(filePath, '\\')
	}

	nameIdx := lastSlash + 1
	fileName := filePath[nameIdx:]
	lowName := strings.ToLower(fileName)
	if lang, ok := languageMap[lowName]; ok {
		return []data.Finding{
			data.NewFinding(data.CatLanguage, lang, "", filePath, 1, nil),
		}
	}

	return nil
}