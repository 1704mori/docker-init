package language

import (
	"os"
	"path/filepath"
)

var FILES_TO_CHECK = []string{"go.mod", "package.json", "requirements.txt", "Cargo.toml"}

func DetectLanguage(project *Project) {
	for _, file := range FILES_TO_CHECK {
		if _, err := os.Stat(file); err == nil {
			switch filepath.Ext(file) {
			case ".mod":
				project.Language = "Go"
			case ".json":
				project.Language = "Node"
			case ".txt":
				project.Language = "Python"
			case ".toml":
				project.Language = "Rust"
			default:
				project.Language = "Other"
			}
		}
	}
}
