package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type AppMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Image       string `json:"image"`
	URL         string `json:"url"`
}

func main() {
	// Create the build directory if it doesn't exist
	outputDir := "build"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Load templates
	indexTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Failed to parse index template: %v", err)
	}

	appTemplate, err := template.ParseFiles("templates/application.html")
	if err != nil {
		log.Fatalf("Failed to parse application template: %v", err)
	}

	// Collect app metadata from the applications directory
	apps, err := collectApplications("applications")
	if err != nil {
		log.Fatalf("Failed to collect applications: %v", err)
	}

	// Generate index page
	if err := generateIndexPage(indexTemplate, apps, outputDir); err != nil {
		log.Fatalf("Failed to generate index page: %v", err)
	}

	// Generate individual application pages
	for _, app := range apps {
		if err := generateAppPage(appTemplate, app, outputDir); err != nil {
			log.Fatalf("Failed to generate page for %s: %v", app.Name, err)
		}
	}

	log.Println("Build completed successfully!")
}

func collectApplications(appDir string) ([]AppMetadata, error) {
	var apps []AppMetadata

	files, err := ioutil.ReadDir(appDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		metadataPath := filepath.Join(appDir, file.Name(), "metadata.json")
		data, err := ioutil.ReadFile(metadataPath)
		if err != nil {
			log.Printf("Skipping %s: failed to read metadata: %v", file.Name(), err)
			continue
		}

		var app AppMetadata
		if err := json.Unmarshal(data, &app); err != nil {
			log.Printf("Skipping %s: failed to parse metadata: %v", file.Name(), err)
			continue
		}

		apps = append(apps, app)
	}

	return apps, nil
}

func generateIndexPage(tmpl *template.Template, apps []AppMetadata, outputDir string) error {
	gridContent := ""
	for _, app := range apps {
		gridContent += `
		<div class="grid-item">
			<img src="` + app.Image + `" alt="` + app.Name + `" />
			<h2>` + app.Name + `</h2>
			<p>` + app.Description + `</p>
			<a href="` + app.Name + `.html">Read More</a>
		</div>
		`
	}

	file, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, map[string]interface{}{
		"Content": template.HTML(gridContent),
	})
}

func generateAppPage(tmpl *template.Template, app AppMetadata, outputDir string) error {
	file, err := os.Create(filepath.Join(outputDir, app.Name+".html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, app)
}