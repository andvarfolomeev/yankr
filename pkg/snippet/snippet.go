package snippet

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
)

func List(snippetsDir string) ([]string, error) {
	files, err := os.ReadDir(snippetsDir)
	if err != nil {
		return nil, err
	}

	var snippets []string
	for _, file := range files {
		if !file.IsDir() {
			snippets = append(snippets, file.Name())
		}
	}

	return snippets, nil
}

func ExtractParams(snippetText string) map[string]string {
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches := re.FindAllStringSubmatch(snippetText, -1)

	paramsMap := make(map[string]string)
	for _, match := range matches {
		if len(match) > 1 {
			paramName := match[1]
			paramsMap[paramName] = ""
		}
	}

	return paramsMap
}

func Process(snippetPath string, cliParams map[string]string) error {
	content, err := os.ReadFile(snippetPath)
	if err != nil {
		return err
	}

	snippetText := string(content)
	paramsMap := ExtractParams(snippetText)

	for param, value := range cliParams {
		paramsMap[param] = value
	}

	missingParams := false
	for _, value := range paramsMap {
		if value == "" {
			missingParams = true
			break
		}
	}

	if missingParams {
		scanner := bufio.NewScanner(os.Stdin)
		for param, value := range paramsMap {
			if value == "" {
				fmt.Printf("Enter value for parameter '%s': ", param)
				scanner.Scan()
				paramsMap[param] = scanner.Text()
			} else {
				fmt.Printf("Using provided value for '%s': %s\n", param, value)
			}
		}
	}

	for param, value := range paramsMap {
		placeholder := "{{" + param + "}}"
		snippetText = strings.ReplaceAll(snippetText, placeholder, value)
	}

	if err := clipboard.WriteAll(snippetText); err != nil {
		return err
	}

	fmt.Println("Snippet copied to clipboard!")
	return nil
}

func Create(snippetsDir, name string) error {
	snippetPath := filepath.Join(snippetsDir, name)

	if _, err := os.Stat(snippetPath); err == nil {
		return fmt.Errorf("snippet '%s' already exists", name)
	}

	if err := os.WriteFile(snippetPath, []byte{}, 0644); err != nil {
		return err
	}

	if err := editSnippetContent(snippetPath); err != nil {
		return err
	}

	fmt.Printf("Snippet '%s' created successfully!\n", name)
	return nil
}

func editSnippetContent(snippetPath string) error {
	editorCmd := os.Getenv("EDITOR")
	if editorCmd != "" {
		return openInEditor(editorCmd, snippetPath)
	}

	return captureFromStdin(snippetPath)
}

func openInEditor(editorCmd string, snippetPath string) error {
	cmd := strings.Fields(editorCmd)
	executable := cmd[0]
	args := append(cmd[1:], snippetPath)

	editor := exec.Command(executable, args...)
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr

	if err := editor.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	return nil
}

func captureFromStdin(snippetPath string) error {
	fmt.Println("Enter snippet content (press Ctrl+D to finish):")
	scanner := bufio.NewScanner(os.Stdin)
	var content strings.Builder
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}

	if err := os.WriteFile(snippetPath, []byte(content.String()), 0644); err != nil {
		return err
	}

	return nil
}

func GetParams(snippetPath string) (map[string]bool, error) {
	content, err := os.ReadFile(snippetPath)
	if err != nil {
		return nil, err
	}

	snippetText := string(content)

	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches := re.FindAllStringSubmatch(snippetText, -1)

	params := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			params[match[1]] = true
		}
	}

	return params, nil
}
