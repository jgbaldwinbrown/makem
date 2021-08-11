package make_make

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
	m.All.Deps = append(m.All.Deps, r.Target)
}

func (m *MakeData) Fprint(w io.Writer) {
	m.All.Target = "all"
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
	Target string
	Deps []string
	Scripts []string
}

func (r Recipe) Fprint(w io.Writer) {
	fmt.Fprintf(w, "%s:", r.Target)
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
