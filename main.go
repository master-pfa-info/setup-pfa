package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("master-pfa: ")

	gopath := getGoPath()
	log.Printf("gopath=%q\n", gopath)
	srcdir := filepath.Join(gopath, "src")
	log.Printf("srcdir=%q\n", srcdir)
	for _, pkg := range pkgs {
		clone(pkg, srcdir)
	}
}

func getGoPath() string {
	p := os.Getenv("GOPATH")
	if p != "" {
		return p
	}
	raw, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(string(raw), "\n")
}

func clone(pkg pkgType, srcdir string) {
	cmd := exec.Command("git", "clone", "--depth=5", pkg.Repo, pkg.Path)
	cmd.Dir = srcdir
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

type pkgType struct {
	Path string
	Repo string
}

var (
	pkgs = []pkgType{
		{"gonum.org/v1/plot", "github.com/gonum/plot"},
		{"go-hep.org/x/hep", "github.com/go-hep/hep"},
		{"bitbucket.org/zombiezen/gopdf", "github.com/master-pfa-info/gopdf"},
		{"golang.org/x/exp", "github.com/golang/exp"},
		{"golang.org/x/mobile", "github.com/golang/mobile"},
		{"golang.org/x/image", "github.com/golang/image"},
	}
)
