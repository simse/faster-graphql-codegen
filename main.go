package main

import (
    "errors"
    "fmt"
    "github.com/briandowns/spinner"
    "github.com/simse/faster-graphql-codegen/internal"
    "os"
    "path/filepath"
    "time"
)

func main() {
	timeStart := time.Now()

	projects := findProjects()

	executionContext := internal.ExecutionContext{}
	executionContext.SetProjects(projects)

	loadSchemas(&executionContext)
	execute(&executionContext, timeStart)
}

func findProjects() []internal.Project {
	// get input folder
	searchFolder := "."
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) > 0 {
		searchFolder = argsWithoutProg[0]
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " Finding projects using codegen"

	s.Start()
	projects, err := internal.FindProjects(searchFolder, filepath.WalkDir)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.FinalMSG = "✗ Input folder does not exist\n"
			s.Stop()
		} else {
			s.FinalMSG = "✗ Unknown error while searching for projects: " + err.Error() + "\n"
			s.Stop()
		}

		os.Exit(1)
	} else {
		s.FinalMSG = fmt.Sprintf("✓ Found %d projects\n", len(projects))
		s.Stop()
	}

	return projects
}

func loadSchemas(e *internal.ExecutionContext) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " Loading graphql schemas"

	s.Start()
	schemasLoaded := e.LoadSchemas()

	s.FinalMSG = fmt.Sprintf("✓ Loaded %d unique schemas\n", schemasLoaded)
	s.Stop()
}

func execute(e *internal.ExecutionContext, timeStart time.Time) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " Executing codegen tasks"

	s.Start()
	e.Execute()

	s.FinalMSG = fmt.Sprintf("✓ Codegen completed in %s\n", time.Since(timeStart).String())
	s.Stop()
}