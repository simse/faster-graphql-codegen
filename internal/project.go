package internal

import (
	"github.com/simse/faster-graphql-codegen/internal/plugins"
	"github.com/vektah/gqlparser/v2/ast"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

type Project struct {
	RootDir    string
	ConfigFile string
	Schemas    []string
	config     Config
}

type ProjectLoadError struct {
	Error    error
	FilePath string
}

type FindProjectsResult struct {
	Projects          []Project
	ProjectLoadErrors []ProjectLoadError
}

func (f *FindProjectsResult) TotalProjectsFound() int {
	return len(f.Projects) + len(f.ProjectLoadErrors)
}

func FindProjects(rootDir string, walkDir func(string, fs.WalkDirFunc) error) (FindProjectsResult, error) {
	result := FindProjectsResult{
		ProjectLoadErrors: make([]ProjectLoadError, 0),
		Projects:          make([]Project, 0),
	}

	// check if path exists
	if _, err := os.Stat(rootDir); err != nil {
		return result, err
	}

	folderSearchErr := walkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && d.Name() == "node_modules" {
			return fs.SkipDir
		}

		if strings.HasSuffix(path, "codegen.ts") || strings.HasSuffix(path, "codegen.yml") {
			absolutePath, _ := filepath.Abs(filepath.Dir(path))
			project := Project{
				RootDir:    absolutePath,
				ConfigFile: d.Name(),
			}

			// prime project
			config, configLoadError := project.GetConfig()
			if configLoadError != nil {
				result.ProjectLoadErrors = append(result.ProjectLoadErrors, ProjectLoadError{
					Error:    configLoadError,
					FilePath: absolutePath,
				})
			} else {
				project.Schemas = config.Schemas
				result.Projects = append(result.Projects, project)
			}
		}

		return nil
	})
	if folderSearchErr != nil {
		return result, folderSearchErr
	}

	return result, nil
}

/*
SchemaKey generates a string get is unique to a combination of schema documents
*/
func (p *Project) SchemaKey() string {
	projectConfig, _ := p.GetConfig()

	sortedSchemas := slices.Clone(projectConfig.Schemas)
	slices.Sort(sortedSchemas)

	return strings.Join(sortedSchemas, ",")
}

type ExecutionContext struct {
	Projects      []Project
	LoadedSchemas map[string]*ast.Schema
}

func (e *ExecutionContext) SetProjects(projects []Project) {
	e.Projects = projects
}

func (e *ExecutionContext) AddLoadedSchema(key string, schema *ast.Schema) {
	if e.LoadedSchemas == nil {
		e.LoadedSchemas = make(map[string]*ast.Schema)
	}

	e.LoadedSchemas[key] = schema
}

func (e *ExecutionContext) GetSchema(key string) *ast.Schema {
	return e.LoadedSchemas[key]
}

/*
LoadSchemas will find every project with a unique list of schemas and load those to cache.
*/
func (e *ExecutionContext) LoadSchemas() int {
	// find unique schemas
	var uniqueSchemas []string
	var projectsToLoad []Project
	for _, project := range e.Projects {
		schemaKey := project.SchemaKey()

		if !slices.Contains(uniqueSchemas, schemaKey) {
			uniqueSchemas = append(uniqueSchemas, schemaKey)
			projectsToLoad = append(projectsToLoad, project)
		}
	}

	// load each schema in parrallel
	var wg sync.WaitGroup
	for _, project := range projectsToLoad {
		wg.Add(1)

		go func() {
			defer wg.Done()

			var schemas []string
			for _, schema := range project.Schemas {
				schemas = append(schemas, path.Join(project.RootDir, schema))
			}

			loadedSchema, err := LoadSchema(schemas...)
			if err != nil {
				panic(err)
			}

			e.AddLoadedSchema(project.SchemaKey(), loadedSchema)
		}()
	}

	wg.Wait()

	return len(uniqueSchemas)
}

func (e *ExecutionContext) Execute() {
	var wg sync.WaitGroup

	for _, project := range e.Projects {
		// get schema from cache
		schema := e.GetSchema(project.SchemaKey())

		// execute all generation tasks
		config, _ := project.GetConfig()
		for destination, destinationConfig := range config.Generates {
			wg.Add(1)

			go func() {
				defer wg.Done()

				// ensure output dir exists
				destinationFile := path.Join(project.RootDir, destination)
				dirCreationErr := EnsureDir(destinationFile)
				if dirCreationErr != nil {
					panic(dirCreationErr)
				}

				// create output string in memory
				output := strings.Builder{}

				e.ExecuteDestinationTasks(destinationConfig, &output, schema, project)

				// create output file
				outputFile, openErr := os.Create(destinationFile)
				if openErr != nil {
					panic(openErr)
				}

				// write output file
				_, writeErr := outputFile.WriteString(output.String())
				if writeErr != nil {
					panic(writeErr)
				}

				// close output files
				err := outputFile.Close()
				if err != nil {
					panic(err)
				}
			}()
		}
	}

	wg.Wait()
}

func (e *ExecutionContext) ExecuteDestinationTasks(
	destinationConfig Generates,
	output *strings.Builder,
	schema *ast.Schema,
	project Project,
) error {
	projectConfig, err := project.GetConfig()
	if err != nil {
		return err
	}

	task := plugins.PluginTask{
		Schema: schema,
		Output: output,
		Config: projectConfig,
	}

	// execute plugins
	for _, plugin := range destinationConfig.Plugins {
		if pluginErr := plugins.VerifyPlugin(plugin, destinationConfig); pluginErr != nil {
			slog.Error("unknown plugin", "plugin", plugin)
		}

		// slog.Info(plugin)

		if plugin == "typescript" {
			task.Typescript()
		}

		if plugin == "introspection" {
			task.Introspect()
		}
	}

	return nil
}

func EnsureDir(filePath string) error {
	dir := filepath.Dir(filePath)

	// Check if the directory already exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create the directory and any necessary parents
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
