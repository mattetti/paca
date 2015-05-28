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
	_ = os.Mkdir("seltest", os.ModePerm)

	helperFilename := "seltest/helper_test.go"
	if _, err := os.Stat(helperFilename); os.IsNotExist(err) || *overwrite {
		f, err := os.Create("seltest/helper_test.go")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(paca.HelperFileContent())
	}

	f, err := os.Create(fmt.Sprintf("seltest/%s_test.go", paca.Camelize(ideCase.Title)))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(ideCase.TestCode())

	fmt.Printf("test case %s converted, check the seltest folder\n", ideCase.Title)
}
