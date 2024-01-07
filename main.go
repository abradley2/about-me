package main

import (
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"os"
	"path"
	"text/template"
)

type JobExperience struct {
	Name      string `json:"name"`
	Title     string `json:"title"`
	DateRange string `json:"dateRange"`
}

type Skill struct {
	Name string `json:"name"`
}

type TemplateData struct {
	Name           string          `json:"name"`
	Skills         []string        `json:"skills"`
	JobExperiences []JobExperience `json:"jobExperiences"`
	ContactInfo    []string        `json:"contactInfo"`
	Blurb          string          `json:"blurb"`
}

func main() {
	var exitCode int
	defer os.Exit(exitCode)

	logger := slog.Default()

	cwd, err := os.Getwd()

	if err != nil {
		logger.Error("Failed to get working directory")
		exitCode = 1
		return
	}

	databuf, err := ioutil.ReadFile(path.Join(cwd, "resume.json"))

	if err != nil {
		logger.Error("Failed to open resume.json file", "error", err)
		exitCode = 1
		return
	}

	var td TemplateData
	err = json.Unmarshal(databuf, &td)

	if err != nil {
		logger.Error("Invalid JSON for template, failed to unmarshal", "error", err)
		exitCode = 1
		return
	}

	tbuff, err := ioutil.ReadFile(path.Join(cwd, "resume.html"))

	if err != nil {
		logger.Error("Failed to open resume.html file", "error", err)
		exitCode = 1
		return
	}

	tpl, err := template.New("resume").Parse(string(tbuff))

	f, err := os.OpenFile(path.Join(cwd, "public", "index.html"), os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer f.Close()

	if err != nil {
		logger.Error("Failed to open/create destination file", "error", err)
		exitCode = 1
		return
	}

	err = tpl.Execute(f, td)

	if err != nil {
		logger.Error("Failed to execute template against data", "error", err)
		exitCode = 1
		return
	}

	logger.Info("Done!")
}
