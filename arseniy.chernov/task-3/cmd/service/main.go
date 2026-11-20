package main

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Input file: %s\n", config.InputFile)
	fmt.Printf("Output file: %s\n", config.OutputFile)
}
