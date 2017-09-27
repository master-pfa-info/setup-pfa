package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	usr *user.User
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("master-pfa: ")

	var err error

	usr, err = user.Current()
	if err != nil {
		log.Fatalf("could not get current user infos: %v", err)
	}

	goroot, gopath, err := installGo("1.9")
	if err != nil {
		log.Fatalf("could not install Go-1.9: %v", err)
	}

	log.Printf("goroot=%q\n", goroot)
	log.Printf("gopath=%q\n", gopath)
	srcdir := filepath.Join(gopath, "src")
	err = os.MkdirAll(srcdir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("srcdir=%q\n", srcdir)
	for _, pkg := range pkgs {
		clone(pkg, srcdir)
	}

	gocmd, err := exec.LookPath("go")
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(gocmd, "get", "-v", "github.com/master-pfa-info/mcpi")
	log.Printf("installing mcpi: %v", cmd.Args)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func installGo(v string) (string, string, error) {
	log.Printf("downloading go-%v...", v)
	burl := "https://golang.org/dl/go" + v + ".linux-amd64.tar.gz"
	resp, err := http.Get(burl)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	goroot := filepath.Join(usr.HomeDir, "M_"+usr.Username, "go-"+v)

	err = os.MkdirAll(goroot, 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("installing go-%v...", v)
	cmd := exec.Command("tar", "zxf", "-")
	cmd.Dir = goroot
	cmd.Stdin = resp.Body
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	goroot = filepath.Join(goroot, "go")
	gopath := filepath.Join(usr.HomeDir, "M_"+usr.Username, "go")
	os.MkdirAll(gopath, 0755)

	os.Setenv("GOROOT", goroot)
	os.Setenv("GOPATH", gopath)
	os.Setenv("PATH",
		strings.Join([]string{
			filepath.Join(goroot, "bin"),
			filepath.Join(gopath, "bin"),
			os.Getenv("PATH"),
		},
			":",
		),
	)

	fname := filepath.Join(usr.HomeDir, ".bashrc")
	err = appendFile(
		fname,
		[]byte(fmt.Sprintf(`
### AUTOMATICALLY added by setup-pfa
export GOROOT=%q
export GOPATH=%q
export PATH=${GOROOT}/bin:${GOPATH}/bin:${PATH}
### AUTOMATICALLY added by setup-pfa [DONE]
`,
			goroot,
			gopath,
		)),
	)
	if err != nil {
		log.Fatalf("could not modify bash_profile: %v", err)
	}

	return goroot, gopath, nil
}

func appendFile(fname string, data []byte) error {
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Seek(0, 2)
	if err != nil {
		log.Fatalf("could not seek to the end: %v", err)
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func clone(pkg pkgType, srcdir string) {
	_, err := os.Stat(filepath.Join(srcdir, pkg.Path))
	if err == nil {
		return
	}
	buf := new(bytes.Buffer)
	cmd := exec.Command("git", "clone", "--depth=5", "https://"+pkg.Repo, pkg.Path)
	cmd.Dir = srcdir
	cmd.Stdout = buf
	cmd.Stderr = buf
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		log.Printf("%v\n", string(buf.Bytes()))
		log.Fatalf("error running %v: %v", cmd.Args, err)
	}
}

type pkgType struct {
	Path string
	Repo string
}

var (
	pkgs = []pkgType{
		{"bitbucket.org/zombiezen/gopdf", "github.com/master-pfa-info/gopdf"},
		{"github.com/golang/freetype", "github.com/golang/freetype"},
		{"github.com/gonuts/binary", "github.com/gonuts/binary"},
		{"github.com/llgcode/draw2d", "github.com/llgcode/draw2d"},
		{"github.com/ajstarks/svgo", "github.com/ajstarks/svgo"},
		{"go-hep.org/x/hep", "github.com/go-hep/hep"},
		{"golang.org/x/exp", "github.com/golang/exp"},
		{"golang.org/x/image", "github.com/golang/image"},
		{"golang.org/x/mobile", "github.com/golang/mobile"},
		{"golang.org/x/net", "github.com/golang/net"},
		{"gonum.org/v1/plot", "github.com/gonum/plot"},
		{"gonum.org/v1/gonum", "github.com/gonum/gonum"},
	}
)
