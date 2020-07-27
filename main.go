package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/loeffel-io/ls-lint/base"
	"gopkg.in/yaml.v2"
)

func getFullPath(path string) string {
	return fmt.Sprintf("%s%s%s", root, sep, path)
}

func main() {
	var files base.FlagFiles
	flag.Var(&files, "f", "Specify alternate config files (default: .ls-lint.yml)")
	flag.Parse()

	var linter = &Linter{
		Errors:  make([]*Error, 0),
		RWMutex: new(sync.RWMutex),
	}

	if len(files) == 0 {
		if err := files.Set(".ls-lint.yml"); err != nil {
			log.Fatal(err)
		}
	}

	var configs []*Config
	for _, f := range files {
		// open config file
		file, err := os.Open(f)

		if err != nil {
			log.Fatal(err)
		}

		// read file
		configBytes, err := ioutil.ReadAll(file)

		if err != nil {
			log.Fatal(err)
		}

		// config
		var tmpConfig = &Config{
			RWMutex: new(sync.RWMutex),
		}

		// to yaml
		if err := yaml.Unmarshal(normalizeConfig(configBytes, byte(runeUnixSep), byte(runeSep)), &tmpConfig); err != nil {
			log.Fatal(err)
		}

		// close file
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}

		// add config
		configs = append(configs, tmpConfig)
	}

	// runner
	if err := linter.Run(configs[0]); err != nil {
		log.Fatal(err)
	}

	// errors
	errors := linter.getErrors()

	// no errors
	if len(errors) == 0 {
		os.Exit(0)
	}

	// with errors
	for _, err := range linter.getErrors() {
		var ruleMessages []string

		for _, rule := range err.getRules() {
			ruleMessages = append(ruleMessages, rule.GetErrorMessage())
		}

		log.Printf("%s failed for rules: %s", err.getPath(), strings.Join(ruleMessages, fmt.Sprintf(" %s ", or)))
	}

	os.Exit(1)
}
