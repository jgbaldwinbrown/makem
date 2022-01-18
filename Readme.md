# Makem

Short for "Make 'em", this small library assists in the generation of makefiles
for running large numbers of concurrent jobs. It is a lightweight alternative
to programs such as [snakemake](https://github.com/snakemake/snakemake). Its
advantage is in building parallel pipelines using the Go language, rather than
attempting complex substitutions using a more traditional make-like tool. The
tool is very explicit, and can be overly verbose for very small makefiles. It
excels, however, when dealing with very large files. It is great for pipelines
that will use hundreds or thousands of inputs and outputs, as recipes can be
written programmatically rather than one-by-one. Being written in an easy-to-use,
high-level, general-purpose language allows for arbitrarily complex scripting.

all: 

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

This complete example creates a makefile, then runs it using all available computing cores. If the . The makefile specifies
that 7 new files, a0, a1, a2, a3, a4, b0, b1, and b2, should be created using "touch" and "cat", with the "b" files depending on the "a" files.

```go
package main

import (
	"fmt"
	"github.com/jgbaldwinbrown/makem"
)

func main() {
	makefile := new(makem.MakeData)

	/* Add targets in a for loop to generate multiple similar recipes: */

	for i:=0; i<5; i++ {
		name := fmt.Sprintf("a%d", i)
		new_rec := makem.Recipe{}
		new_rec.AddTargets(name)
		new_rec.AddScripts(fmt.Sprintf("touch %s", name))
		makefile.Add(new_rec)
	}

	/* Multiple targets and dependencies are possible, as are multiple scripts per recipe: */

	new_rec := makem.Recipe{}
	new_rec.AddTargets("b0", "b1")
	new_rec.AddDeps("a0", "a1")
	new_rec.AddScripts("cat a0 > b0", "cat a1 > b1")
	makefile.Add(new_rec)

	/* Alternative literal syntax: /*

	makefile.Add(makem.Recipe{
		Targets: []string{"b2"},
		Deps: []string{"a2"},
		Scripts: []string{cat a2 > b2"},
	})

	/* Run the makefile */

	makefile.Exec(makem.UseAllCores())
}
```

If the makefile were printed using the line `makefile.Fprint(os.Stdout)`, it would
produce the following output:

```make
all: a0 a1 a2 a3 a4 b0 b1 b2

a0:
	touch a0

a1:
	touch a1

a2:
	touch a2

a3:
	touch a3

a4:
	touch a4

b0 b1: a0 a1
	cat a0 > b0
	cat a1 > b1
b2: a2
	cat a2 > b2

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

#### MakeData.Add

```go
func (m *MakeData) Add(rs ...Recipe)
```

This function adds a set of recipes to the makefile.

#### MakeData.Fprint

```go
func (m *MakeData) Fprint(w io.Writer)
```

This method prints the full makefile to the specified io.Writer.

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

#### Recipe.AddTargets

```go
func (r *Recipe) AddTargets(ts ...string)
```

Add the specified targets to the recipe.

#### Recipe.AddDeps

```go
func (r *Recipe) AddDeps(ts ...string)
```

Add the specified dependencies to the recipe.

#### Recipe.AddScripts

```go
func (r *Recipe) AddScripts(ts ...string)
```

Add the specified scripts to the recipe.

#### Recipe.Fprint

```go
func (r Recipe) Fprint(w io.Writer)
```

Print the recipe to a specified `io.Writer`.

### FprintRecipes

```go
FprintRecipes(w io.Writer, rs ...Recipe)
```

Print all specified recipes to an `io.Writer`.

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

#### UseCores

```go
func UseCores(corenum int) ExecOption
```

This function is used as an option when running "MakeData.Exec". It sets Exec
to use a fixed number of cores (`corenum`) when running. Example:

```go
m := MakeData{}
/* ... */
done := m.Exec(UseCores(8))
```

#### UseAllCores

```go
func UseAllCores() ExecOption
```

This function is used as an option when running "MakeData.Exec". It sets Exec
to use all available cores (`corenum`) when running. Example:

```go
m := MakeData{}
/* ... */
done := m.Exec(UseAllCores())
```
