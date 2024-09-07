package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Load the YAML file into a map for dynamic handling
func loadYAML(file string) (map[string]interface{}, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var result map[string]interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %w", err)
	}
	return result, nil
}

// Process template file and write to corresponding YAML file
func processTemplate(tplFile string, params map[string]interface{}) error {
	// Parse the template file
	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		return fmt.Errorf("error parsing template file: %w", err)
	}

	// Determine the output file name by replacing .tpl with .yaml
	outputFile := strings.Replace(tplFile, ".tpl", ".yaml", 1)

	// Create or overwrite the output file
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer out.Close()

	// Execute the template with the loaded parameters
	err = tpl.Execute(out, params)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	fmt.Printf("Generated: %s\n", outputFile)
	return nil
}

func crawlDirectory(rootDir string) error {
	// Walk through the `applications` directory and look for version directories
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Check if the current path contains a `compose/` directory
		if d.IsDir() && strings.Contains(path, "compose") {
			// Load `parameters.yaml` from the parent directory of the `compose/`
			parentDir := filepath.Dir(path)
			paramsFile := filepath.Join(parentDir, "parameters.yaml")
			params, err := loadYAML(paramsFile)
			if err != nil {
				log.Printf("Error loading parameters.yaml in %s: %v", path, err)
				return nil
			}

			// Now search for `.tpl` files in the `compose/` subdirectory
			err = filepath.WalkDir(path, func(filePath string, fileInfo fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				// Process `.tpl` files found in the compose/ directory
				if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), ".tpl") {
					err := processTemplate(filePath, params)
					if err != nil {
						log.Printf("Error processing template %s: %v", filePath, err)
					}
				}
				return nil
			})
			if err != nil {
				log.Printf("Error walking directory %s: %v", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	return nil
}

func main() {
	// Start crawling from the `applications` directory
	rootDir := "applications"

	// Crawl the directory and process templates
	err := crawlDirectory(rootDir)
	if err != nil {
		log.Fatalf("Error crawling directory: %v", err)
	}
}