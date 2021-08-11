package main

import (
	"os/exec"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type recipe struct {
	target string
	deps []string
	scripts []string
}

func (r recipe) Fprint(w io.Writer) {
	fmt.Fprintf(w, "%s:", r.target)
	for _, d := range r.deps {
		fmt.Fprintf(w, " %s", d)
	}
	fmt.Fprintf(w, "\n")
	for _, s := range r.scripts {
		fmt.Fprintf(w, "\t%s\n", s)
	}
}

func FprintRecipes(w io.Writer, rs []recipe) {
	for _, r := range rs {
		r.Fprint(w)
	}
}

func main() {
	all := recipe{target: "all"}
	var recipes []recipe
	for i:=0; i<5; i++ {
		name := fmt.Sprintf("a%d", i)
		all.deps = append(all.deps, name)
		new_rec := recipe{target:name}
		new_rec.scripts = append(new_rec.scripts, fmt.Sprintf("touch %s", name))
		recipes = append(recipes, new_rec)
	}
	tmpfile, err := ioutil.TempFile("", "Makefile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name())
	all.Fprint(tmpfile)
	FprintRecipes(tmpfile, recipes)
	tmpfile.Close()
	exec.Command("make", "-f", tmpfile.Name()).Run()
}
