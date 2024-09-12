package internal

import (
    "github.com/vektah/gqlparser/v2/ast"
    "io/fs"
    "os"
    "path"
    "path/filepath"
    "slices"
    "strings"
    "sync"
)

type Project struct {
	RootDir string
	ConfigFile string
	Schemas []string
	config Config
}

func FindProjects(rootDir string, walkDir func(string, fs.WalkDirFunc) error) []Project {
	var projects []Project

	err := walkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, "codegen.ts") || strings.HasSuffix(path, "codegen.yml") {
			project := Project{
				RootDir: filepath.Dir(path),
				ConfigFile: d.Name(),
			}

			// prime project
			config := project.GetConfig()
			project.Schemas = config.Schema

			projects = append(projects, project)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return projects
}

/*
SchemaKey generates a string get is unique to a combination of schema documents
*/
func (p *Project) SchemaKey() string {
	projectConfig := p.GetConfig()

	sortedSchemas := slices.Clone(projectConfig.Schema)
	slices.Sort(sortedSchemas)

	return strings.Join(sortedSchemas, ",")
}

type ExecutionContext struct {
	Projects []Project
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
func (e *ExecutionContext) LoadSchemas() {
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
}

func (e *ExecutionContext) Execute() {
	var wg sync.WaitGroup

	for _, project := range e.Projects {
		// get schema from cache
		schema := e.GetSchema(project.SchemaKey())

		// execute all generation tasks
		config := project.GetConfig()
		for destination, _ := range config.Generates {
			wg.Add(1)

			go func() {
				defer wg.Done()

				destinationFile := path.Join(project.RootDir, destination)

				dirCreationErr := EnsureDir(destinationFile)
				if dirCreationErr != nil {
					panic(dirCreationErr)
				}

				output := strings.Builder{}
				ConvertSchema(schema, &output)

				outputFile, openErr := os.Create(destinationFile)
			    if openErr != nil {
			        panic(openErr)
			    }

				_, writeErr := outputFile.WriteString(output.String())
				if writeErr != nil {
					panic(writeErr)
				}

		        err := outputFile.Close()
		        if err != nil {
		            panic(err)
		        }
			}()
		}
	}

	wg.Wait()
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