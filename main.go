package main

import (
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/gookit/color"
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
	projectSearchResult, err := internal.FindProjects(searchFolder, filepath.WalkDir)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.FinalMSG = errorString("Input folder does not exist: %s", searchFolder)
			s.Stop()
			println()
		} else {
			s.FinalMSG = "✗ Unknown error while searching for projects: " + err.Error() + "\n"
			s.Stop()
		}

		os.Exit(1)
	} else {
		s.FinalMSG = successString("Found %d projects\n", projectSearchResult.TotalProjectsFound())
		s.Stop()

		if len(projectSearchResult.ProjectLoadErrors) > 0 {
			fmt.Println(errorString("%d project config files failed to load", len(projectSearchResult.ProjectLoadErrors)))

			for _, loadError := range projectSearchResult.ProjectLoadErrors {
				color.Gray.Println("\t" + loadError.FilePath)
				fmt.Println("\t↳ " + loadError.Error.Error())
			}
		}
	}

	return projectSearchResult.Projects
}

func loadSchemas(e *internal.ExecutionContext) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " Loading graphql schemas"

	s.Start()
	schemasLoaded := e.LoadSchemas()

	if schemasLoaded == 0 {
		s.FinalMSG = errorString("No schemas loaded. Did any config files load?\n")
		s.Stop()
		os.Exit(1)
	}

	s.FinalMSG = successString("Loaded %d unique schemas\n", schemasLoaded)
	s.Stop()
}

func execute(e *internal.ExecutionContext, timeStart time.Time) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " Executing codegen tasks"

	s.Start()
	e.Execute()

	s.FinalMSG = successString("Codegen completed in %s\n", time.Since(timeStart).String())
	s.Stop()
}

func greenTick() string {
	return color.Green.Sprint("✓")
}

func redCross() string {
	return color.Red.Sprintf("✗")
}

func successString(format string, arguments ...interface{}) string {
	return greenTick() + " " + fmt.Sprintf(format, arguments...)
}

func errorString(format string, arguments ...interface{}) string {
	return redCross() + " " + fmt.Sprintf(format, arguments...)
}
