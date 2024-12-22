package utils

import (
	"fmt"
	"grimoire/app/config"
	"grimoire/app/log"
	"html/template"
	"os"
)

var cfg = config.Get()

func LoadAndGenerateHTML(gitInfoPath string) error {
	gitInfo, err := LoadGitInfo(gitInfoPath)
	if err != nil {
		logger.Error("failed to load git info",
			log.String("error", err.Error()),
			log.String("path", gitInfoPath),
		)
		return err
	}

	data := GitInfo{
		CommitHash: gitInfo.CommitHash,
		Branch:     gitInfo.Branch,
		CommitDate: gitInfo.CommitDate,
		Authors:    gitInfo.Authors,
		Version:    gitInfo.Version,
		Commits:    gitInfo.Commits,
	}

	err = GenerateHTML(data)
	if err != nil {
		logger.Error("failed to generate HTML",
			log.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func GenerateHTML(data GitInfo) error {
	tmplContent, err := os.ReadFile(cfg.Service.TemplateWelcome)
	if err != nil {
		logger.Error("error reading HTML template",
			log.String("err", err.Error()),
		)
		return err
	}

	tmplParsed, err := template.New("welcome").Parse(string(tmplContent))
	if err != nil {
		logger.Error("error parsing template",
			log.String("err", err.Error()),
		)
		return err
	}

	shortCommitHash := ""
	if len(data.CommitHash) >= 20 {
		shortCommitHash = data.CommitHash[:20]
	} else {
		shortCommitHash = data.CommitHash
	}

	var authorsList string
	for _, author := range data.Authors {
		authorsList += fmt.Sprintf("%d commits - %s\n", author.Commits, author.Name)
	}

	templateData := struct {
		ShortCommitHash string
		Branch          string
		CommitDate      string
		AuthorsList     string
		Version         string
		Commits         []CommitDay
	}{
		ShortCommitHash: shortCommitHash,
		Branch:          data.Branch,
		CommitDate:      data.CommitDate,
		AuthorsList:     authorsList,
		Version:         data.Version,
		Commits:         data.Commits,
	}

	err = os.MkdirAll("./static", os.ModePerm)
	if err != nil {
		logger.Error("error creating static directory",
			log.String("err", err.Error()),
		)
		return err
	}

	outputFile := cfg.Service.TemplateStatic
	file, err := os.Create(outputFile)
	if err != nil {
		logger.Error("error creating HTML file",
			log.String("err", err.Error()),
		)
		return err
	}
	defer file.Close()

	err = tmplParsed.Execute(file, templateData)
	if err != nil {
		logger.Error("error executing template",
			log.String("err", err.Error()),
		)
		return err
	}
	logger.Debug("Static HTML file generated successfully in",
		log.String("directory", outputFile),
	)
	return nil
}