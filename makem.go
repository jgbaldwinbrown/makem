package makem

import (
	"os/exec"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type MakeData struct {
	All Recipe
	Preamble string
	Recipes []Recipe
}

type ExecInternal struct {
	Parallel bool
	Cores int
	AllCores bool
}

type ExecOption func(e *ExecInternal)

func (m *MakeData) Add(r Recipe) {
	m.Recipes = append(m.Recipes, r)
	m.All.Deps = append(m.All.Deps, r.Targets...)
}

func (m *MakeData) Fprint(w io.Writer) {
	m.All.AddTarget("all")
	m.All.Fprint(w)
	fmt.Fprint(w, m.Preamble)
	FprintRecipes(w, m.Recipes)
}

func UseCores(corenum int) ExecOption {
	return func(e *ExecInternal) {
		e.Parallel = true
		e.Cores = corenum
	}
}

func UseAllCores() ExecOption {
	return func(e *ExecInternal) {
		e.AllCores = true
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
	if settings.Parallel {
		jobs_string = fmt.Sprintf("-j %v", settings.Cores)
	}
	if settings.AllCores {
		jobs_string = fmt.Sprintf("-j")
	}

	exec.Command("make", jobs_string, "-f", tmpfile.Name()).Run()

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

func (r *Recipe) AddTarget(t string) {
	r.Targets = append(r.Targets, t)
}

func (r *Recipe) AddTargets(ts []string) {
	for _, t := range ts {
		r.AddTarget(t)
	}
}

func (r *Recipe) AddDep(t string) {
	r.Deps = append(r.Deps, t)
}

func (r *Recipe) AddDeps(ts []string) {
	for _, t := range ts {
		r.AddDep(t)
	}
}

func (r *Recipe) AddScript(t string) {
	r.Scripts = append(r.Scripts, t)
}

func (r *Recipe) AddScripts(ts []string) {
	for _, t := range ts {
		r.AddScript(t)
	}
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
}

func FprintRecipes(w io.Writer, rs []Recipe) {
	for _, r := range rs {
		r.Fprint(w)
	}
}
