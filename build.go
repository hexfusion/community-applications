package main

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type AppMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	Category    string `json:"category"`
}

func main() {
	outputDir := "public"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// load templates
	indexTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Failed to parse index template: %v", err)
	}

	appTemplate, err := template.ParseFiles("templates/application.html")
	if err != nil {
		log.Fatalf("Failed to parse application template: %v", err)
	}

	// sort app metadata from the applications directory
	appsByCategory, categories, err := collectAndSortApplicationsByCategory("applications")
	if err != nil {
		log.Fatalf("Failed to collect applications: %v", err)
	}

	// index page
	if err := generateIndexPage(indexTemplate, appsByCategory, categories, outputDir); err != nil {
		log.Fatalf("Failed to generate index page: %v", err)
	}

	// app pages
	for _, apps := range appsByCategory {
		for _, app := range apps {
			if err := generateAppPage(appTemplate, app, outputDir); err != nil {
				log.Fatalf("Failed to generate page for %s: %v", app.Name, err)
			}
		}
	}

	log.Println("Static site generated successfully!")
}

// collectAndSortApplicationsByCategory collects apps and groups them by category, then sorts within each category.
// It also returns a list of unique categories for the sidebar.
func collectAndSortApplicationsByCategory(appDir string) (map[string][]AppMetadata, []string, error) {
	appsByCategory := make(map[string][]AppMetadata)
	uniqueCategories := make(map[string]bool)

	files, err := os.ReadDir(appDir)
	if err != nil {
		return nil, nil, err
	}

	for _, file := range files {
		metadataPath := filepath.Join(appDir, file.Name(), "metadata.json")
		data, err := os.ReadFile(metadataPath)
		if err != nil {
			log.Printf("Skipping %s: failed to read metadata: %v", file.Name(), err)
			continue
		}

		var app AppMetadata
		if err := json.Unmarshal(data, &app); err != nil {
			log.Printf("Skipping %s: failed to parse metadata: %v", file.Name(), err)
			continue
		}

		// Group by category
		appsByCategory[app.Category] = append(appsByCategory[app.Category], app)

		// Track unique categories for the sidebar
		uniqueCategories[app.Category] = true
	}

	// Sort apps within each category
	for category := range appsByCategory {
		sort.Slice(appsByCategory[category], func(i, j int) bool {
			return appsByCategory[category][i].Name < appsByCategory[category][j].Name
		})
	}

	// Convert the uniqueCategories map to a sorted slice
	categories := []string{}
	for category := range uniqueCategories {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	return appsByCategory, categories, nil
}

func generateIndexPage(tmpl *template.Template, appsByCategory map[string][]AppMetadata, categories []string, outputDir string) error {
	var gridContent string

	// grid content
	for _, apps := range appsByCategory {
		gridContent += `<div class="grid-container">` // Start grid for this category

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

		gridContent += "</div>\n"
	}

	// sidebar content
	sidebarContent := "<ul>\n"
	for _, category := range categories {
		sidebarContent += `<li><a href="#">` + category + `</a></li>` + "\n"
	}
	sidebarContent += "</ul>"

	// create index file
	file, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, map[string]interface{}{
		"Content": template.HTML(gridContent),
		"Sidebar": template.HTML(sidebarContent),
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
