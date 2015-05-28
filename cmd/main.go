package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mattetti/paca"
)

var (
	source    = flag.String("source", "", "source path to the ide case path to convert")
	overwrite = flag.Bool("overwrite", false, "overwrite the generated helper files")
	dest      = flag.String("dest", "seltest", "path to the directory where to dump the converted cases")
	pkgName   = flag.String("pkg", "seltest", "name of the Go package to use when generating cases")
)

func main() {
	flag.Parse()
	if *source == "" {
		fmt.Println("pass the source path to the ide case file")
		os.Exit(1)
	}
	ideCase, err := paca.IDEConverter(*source)
	if err != nil {
		fmt.Printf("failed to process %s - %v\n", *source, err)
		os.Exit(1)
	}
	_ = os.Mkdir(*dest, os.ModePerm)

	helperFilename := fmt.Sprintf("%s/helper_test.go", *dest)
	if _, err := os.Stat(helperFilename); os.IsNotExist(err) || *overwrite {
		f, err := os.Create(helperFilename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(paca.HelperFileContent(*pkgName))
	}

	f, err := os.Create(fmt.Sprintf("%s/%s_test.go", *dest, paca.Camelize(ideCase.Title)))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(ideCase.TestCode(*pkgName))

	fmt.Printf("test case %s converted, check the %s folder\n", *dest, ideCase.Title)
}
