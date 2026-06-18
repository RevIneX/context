package detectors

import (
	"context/internal/data"
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
	".zsh": "Shell", ".ps1": "PowerShell", ".sql": "SQL", ".r": "R",
	".lua": "Lua", ".dart": "Dart", ".ex": "Elixir", ".exs": "Elixir",
	".erl": "Erlang", ".hrl": "Erlang", ".hs": "Haskell", ".lhs": "Haskell",
	".clj": "Clojure", ".cljs": "ClojureScript", ".elm": "Elm", ".vue": "Vue",
	".svelte": "Svelte", ".astro": "Astro", ".tf": "HCL", ".hcl": "HCL",
	".yaml": "YAML", ".yml": "YAML", ".toml": "TOML", ".json": "JSON",
	".xml": "XML", ".ini": "INI", ".cfg": "INI", ".conf": "INI",
	".md": "Markdown", ".txt": "Text", ".css": "CSS", ".scss": "SCSS",
	".sass": "Sass", ".less": "Less", ".html": "HTML", ".htm": "HTML",

	"dockerfile":         "Docker",
	"docker-compose.yml": "Docker Compose",
	"makefile":           "Make",
	".env":               "Env",
	".gitlab-ci.yml":     "GitLab CI",
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