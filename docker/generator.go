package docker

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/1704mori/docker-init/language"
	"github.com/1704mori/docker-init/templates"
)

var DockerfileTemplates = map[string]string{
	"Go":     templates.GO,
	"Node":   templates.NODE,
	"Python": templates.PYTHON,
}

func GenerateFiles(project *language.Project) {
	// dockerfileTemplate, err := getDockerfileTemplate(project.Language)
	dockerfileTemplate, found := DockerfileTemplates[project.Language]
	if !found {
		fmt.Println("Unsupported language")
		os.Exit(1)
	}

	dockerfileContet, err := executeTemplate(dockerfileTemplate, project)
	if err != nil {
		fmt.Println("Error executing Dockerfile template", err)
		os.Exit(1)
	}

	writeToFile("Dockerfile", dockerfileContet)
}

func getDockerfileTemplate(language string) (string, error) {
	templateFile := fmt.Sprintf("templates/%s.Dockerfile", language)
	templateData, err := os.ReadFile(templateFile)
	if err != nil {
		return "", err
	}

	return string(templateData), nil
}

func executeTemplate(templateStr string, data interface{}) (string, error) {
	tmpl, err := template.New("dockerfile").Parse(templateStr)
	if err != nil {
		return "", err
	}

	buffer := &strings.Builder{}
	if err := tmpl.Execute(buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func writeToFile(filename, content string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("Error writing to file file %s: %v\n", filename, err)
		os.Exit(1)
	}
}

func readFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	return string(content)
}
