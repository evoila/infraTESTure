# parser
--
    import "github.com/evoila/infraTESTure/parser"


## Usage

#### func  GetAnnotations

```go
func GetAnnotations(path string) []string
```
Return a list of all annotations in a given go file 

@param path Path to the test project 

@return []string List of annotations from functions within the test project

#### func  GetFunctionNames

```go
func GetFunctionNames(test config.Test, path string) ([]string, error)
```
Return a list of all functions whose annotations match with the test names in the config 

@param test Initialized test struct from github.com/evoila/infraTESTure/config 

@param path Path to the test project

@return []string Name list of available functions in the test project
