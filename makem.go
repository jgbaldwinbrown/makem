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
	Recipes []Recipe
}

func (m *MakeData) Add(r Recipe) {
	m.Recipes = append(m.Recipes, r)
	m.All.Deps = append(m.All.Deps, r.Targets...)
}

func (m *MakeData) Fprint(w io.Writer) {
	m.All.AddTarget("all")
	m.All.Fprint(w)
	FprintRecipes(w, m.Recipes)
}

func (m *MakeData) Exec() (err error) {
	tmpfile, err := ioutil.TempFile("", "Makefile")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())
	m.Fprint(tmpfile)
	tmpfile.Close()
	exec.Command("make", "-f", tmpfile.Name()).Run()
	return nil
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
