package makem

import (
	"os/exec"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type MakeData struct {
	Preamble string
	Recipes []Recipe
}

type ExecInternal struct {
	Parallel bool
	Cores int
	AllCores bool
}

type ExecOption func(e *ExecInternal)

func (m *MakeData) Add(rs ...Recipe) {
	m.Recipes = append(m.Recipes, rs...)
}

func (m *MakeData) Fprint(w io.Writer) {
	all := Recipe{}
	all.AddTargets("all")

	for _, recipe := range m.Recipes {
		all.AddDeps(recipe.Targets...)
	}

	all.Fprint(w)
	fmt.Fprint(w, m.Preamble)
	FprintRecipes(w, m.Recipes...)
}

func UseCores(corenum int) ExecOption {
	return func(e *ExecInternal) {
		e.Parallel = true
		e.Cores = corenum
		e.AllCores = false
	}
}

func UseAllCores() ExecOption {
	return func(e *ExecInternal) {
		e.AllCores = true
		e.Parallel = false
	}
}

func (m *MakeData) Exec(options ...ExecOption) (err error) {
	settings := ExecInternal{}
	for _, option := range options {
		option(&settings)
	}
	tmpfile, err := ioutil.TempFile("", "Makefile")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())
	m.Fprint(tmpfile)
	tmpfile.Close()

	jobs_string := ""
	jobs_count := ""
	if settings.Parallel {
		jobs_string = fmt.Sprintf("-j")
		jobs_count = fmt.Sprintf("%v", settings.Cores)
	}
	if settings.AllCores {
		jobs_string = fmt.Sprintf("-j")
	}

	command := exec.Command("make", jobs_string, "-f", tmpfile.Name())
	if settings.Parallel {
		command = exec.Command("make", jobs_string, jobs_count, "-f", tmpfile.Name())
	}
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Run()

	return nil
}

func (m *MakeData) AppendPreamble(s string) {
	m.Preamble = m.Preamble + s
}

func (m *MakeData) SetPreamble(s string) {
	m.Preamble = s
}

type Recipe struct {
	Targets []string
	Deps []string
	Scripts []string
}

func (r *Recipe) AddTargets(ts ...string) {
	r.Targets = append(r.Targets, ts...)
}

func (r *Recipe) AddDeps(ts ...string) {
	r.Deps = append(r.Deps, ts...)
}

func (r *Recipe) AddScripts(ts ...string) {
	r.Scripts = append(r.Scripts, ts...)
}

func (r Recipe) Fprint(w io.Writer) {
	fmt.Fprintf(w, "%s", r.Targets[0])
	for _, t := range r.Targets[1:] {
		fmt.Fprintf(w, " %s", t)
	}
	fmt.Fprintf(w, ":")
	for _, d := range r.Deps {
		fmt.Fprintf(w, " %s", d)
	}
	fmt.Fprintf(w, "\n")
	for _, s := range r.Scripts {
		fmt.Fprintf(w, "\t%s\n", s)
	}
	fmt.Fprintf(w, "\n")
}

func FprintRecipes(w io.Writer, rs ...Recipe) {
	for _, r := range rs {
		r.Fprint(w)
	}
}
