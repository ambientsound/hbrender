package main

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/ghodss/yaml"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"os"
)

type Config struct {
	TemplateFile  string
	VariablesFile string
}

var (
	cfg Config
)

type TemplateVariables map[string]interface{}

func templatedFile(data []byte, ctx TemplateVariables) ([]byte, error) {
	template, err := raymond.Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse template file: %s", err)
	}

	output, err := template.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("execute template: %s", err)
	}

	return []byte(output), nil
}

func templateVariablesFromFile(path string) (TemplateVariables, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: open file: %s", path, err)
	}

	vars := TemplateVariables{}
	err = yaml.Unmarshal(file, &vars)

	return vars, err
}

func init() {
	flag.StringVar(&cfg.TemplateFile, "template", "", "Handlebars template file")
	flag.StringVar(&cfg.VariablesFile, "vars", "/dev/null", "JSON or YAML file with variables")
	flag.Parse()
}

func run() error {
	if len(cfg.TemplateFile) == 0 {
		flag.Usage()
		return fmt.Errorf("template file must be specified")
	}

	template, err := ioutil.ReadFile(cfg.TemplateFile)
	if err != nil {
		return err
	}

	vars, err := templateVariablesFromFile(cfg.VariablesFile)
	if err != nil {
		return err
	}

	rendered, err := templatedFile(template, vars)
	if err != nil {
		return err
	}

	fmt.Print(string(rendered))

	return nil
}

func main() {
	err := run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fatal: %s\n", err.Error())
	}
}
