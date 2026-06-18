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
	"dockerfile": {}, "docker-compose.yml": {}, "makefile": {}, ".env": {},
	".gitlab-ci.yml": {},
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