# Makem

Short for "Make 'em", this small library assists in the generation of makefiles
for running large numbers of concurrent jobs. It is a lightweight alternative
to programs such as [snakemake](https://github.com/snakemake/snakemake).

## Installation

Simply run:

```sh
go get github.com/jgbaldwinbrown/makem
```

Then, in your file, import the library as follows:

```go
import (
	"github.com/jgbaldwinbrown/makem"
)
```

## A simple example

This complete example creates a makefile, then runs it. The makefile specifies
that 7 new files, a0, a1, a2, a3, a4, b0, and b1, should be created using "touch".

```go
package main

import (
	"fmt"
	"local/jgbaldwinbrown/makem"
	"os"
)

func main() {
	makefile := new(makem.MakeData)

	for i:=0; i<5; i++ {
		name := fmt.Sprintf("a%d", i)
		new_rec := makem.Recipe{}
		new_rec.AddTarget(name)
		new_rec.Scripts = append(new_rec.Scripts, fmt.Sprintf("touch %s", name))
		makefile.Add(new_rec)
	}

	new_rec := makem.Recipe{}
	new_rec.AddTargets([]string{"b0", "b1"})
	new_rec.AddScripts([]string{
		"touch b0",
		"touch b1",
	})
	makefile.Add(new_rec)

	makefile.Exec(makem.UseAllCores())
}
```

## Documentation

### MakeData

```go
type MakeData struct {
	All Recipe
	Preamble string
	Recipes []Recipe
}
```

This is the main type in the library. It holds all of the recipes in the
makefile, plus the special "All" recipe which depends on all other recipes,
plus any preamble added by the user.

#### MakeData.Fprint

```go
func (m *MakeData) Fprint(w io.Writer)
```

This method allows printing of the full makefile to the specified io.Writer.

#### MakeData.Exec

```go
func (m *MakeData) Exec(options ...ExecOption) (err error)
```

This method runs the makefile with the specified options. It return `nil` on
success and an error on failure.

#### MakeData.AppendPreamble

```go
func (m *MakeData) AppendPreamble(s string)
```

This appends to the existing preamble. Note that, for full flexibility, no
trailing newline is appended by default.

#### MakeData.SetPreamble

```go
func (m *MakeData) SetPreamble(s string)
```

This sets the preamble. Note that, for full flexibility, no
trailing newline is appended by default. This is a good place to put build options and
make-specific variables.

### Recipe

```go
type Recipe struct {
	Targets []string
	Deps []string
	Scripts []string
}
```

A recipe consists of a target file to create, all dependencies of the target,
and all scripts used to generate the target.  Tabs will automatically be
prepended to script lines. Multiple targets and multiple dependencies are
supported, and should be represented as multiple items in a `[]string`. Recipes
are intended to be generated with the helper methods below, but can be
generated with literals as well:

```go
myrecipe := Recipe {
	Targets: []string{"a.txt", "b.txt"},
	Deps: []string{"c.txt"},
	Scripts: []string{"cat c.txt > a.txt", "cat c.txt > b.txt"},
}
```

#### Recipe.AddTarget



### ExecInternal

```go
type ExecInternal struct {
	Parallel bool
	Cores int
	AllCores bool
}
```

This is mainly for internal use. It holds all of the options for executing the
makefile. These should not be set manually, but rather should be set using an
ExecOption function.

### ExecOption

```go
type ExecOption func(e *ExecInternal)
```

This type is the only argument type for the Exec 
