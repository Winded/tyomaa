package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"text/template"

	"github.com/winded/go-webgen"
)

var (
	debug bool
)

type TemplateData struct {
	ApiEndpoint string
}

func init() {
	flag.BoolVar(&debug, "debug", false, "Build in debug mode")
}

func buildJs(output string) error {
	if err := os.MkdirAll(path.Dir(output), 0755); err != nil {
		return err
	}

	params := []string{"build", "github.com/winded/tyomaa/frontend/js", "-o", output}
	if !debug {
		params = append(params, "-m")
	}

	cmd := exec.Command(os.ExpandEnv("$GOPATH/bin/gopherjs"), params...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func env(key, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if found {
		return value
	} else {
		return defaultValue
	}
}

func main() {
	flag.Parse()

	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	rootDir = path.Join(rootDir, "..")

	if err := buildJs(path.Join(rootDir, "static", "js", "bundle.js")); err != nil {
		panic(err)
	}

	outputDir := env("OUTPUT_DIR", path.Join(rootDir, "output"))
	templateDir := path.Join(rootDir, "templates")

	wgen := webgen.NewGenerator(webgen.GeneratorConfig{
		OutputDir:   outputDir,
		StaticDir:   path.Join(rootDir, "static"),
		TemplateDir: templateDir,

		CompressStaticFiles: debug,
		URLPrefix:           env("PATH_PREFIX", ""),
		StaticOutputPrefix:  "static",
	})

	wgen.Funcs(template.FuncMap{
		"view": func(p string) (string, error) {
			data, err := ioutil.ReadFile(path.Join(templateDir, "views", p))
			return string(data), err
		},
	})

	wgen.Add("/index.html", "index", &TemplateData{
		ApiEndpoint: env("API_ENDPOINT", "localhost"),
	})

	if err := wgen.Generate(); err != nil {
		panic(err)
	}

	if debug {
		http.Handle("/", http.FileServer(http.Dir(outputDir)))
		fmt.Println("Serving files from " + outputDir)
		http.ListenAndServe(":80", nil)
	}
}
