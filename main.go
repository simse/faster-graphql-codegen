package main

import (
    "github.com/lmittmann/tint"
    "github.com/simse/faster-graphql-codegen/internal"
    "log/slog"
    "os"
    "path/filepath"
    "time"
)

func main() {
	// set up coloured logging
	w := os.Stderr

	// set global logger with custom options
	slog.SetDefault(slog.New(
	    tint.NewHandler(w, &tint.Options{
	        Level:      slog.LevelDebug,
	        TimeFormat: time.TimeOnly,
	    }),
	))

	/*timeStart := time.Now()
	schema, err := internal.LoadSchema("examples/github/github.graphql")
	if err != nil {
		panic(err)
	}
	loadSchemaTime := time.Since(timeStart)

	convertTimeStart := time.Now()
	output := strings.Builder{}
	internal.ConvertSchema(schema, &output)
	convertSchemaTime := time.Since(convertTimeStart)

	writeTypesTimeStart := time.Now()
	outputFile, openErr := os.Create("output/github/baseTypes.ts")
    if openErr != nil {
        panic(err)
    }

	_, writeErr := outputFile.WriteString(output.String())
	if writeErr != nil {
		panic(err)
	}
	writeTypesTime := time.Since(writeTypesTimeStart)

	outputFile.Close()
	slog.Info("finished codegen", "duration", time.Since(timeStart), "load_duration", loadSchemaTime, "convert_duration", convertSchemaTime, "write_duration", writeTypesTime)*/
	timeStart := time.Now()
	projects := internal.FindProjects("./examples/projects/monorepo", filepath.WalkDir)
	slog.Info("search for packages using codegen complete", "found_projects", len(projects))

	/*for _, project := range projects {
		internal.ExecuteProject(project)
	}*/
	executionContext := internal.ExecutionContext{}
	executionContext.SetProjects(projects)
	executionContext.LoadSchemas()
	executionContext.Execute()
	slog.Info("finished all codegen tasks", "duration", time.Since(timeStart))
}