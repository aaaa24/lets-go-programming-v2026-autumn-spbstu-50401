package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"go.yaml.in/yaml/v4"
	"golang.org/x/text/encoding/charmap"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	XMLName     xml.Name `xml:"ValCurs"`
	Date        string   `xml:"Date,attr"`
	Name        string   `xml:"name,attr"`
	ValuteItems []Valute `xml:"Valute"`
}

type Valute struct {
	XMLName   xml.Name `xml:"Valute"`
	ID        string   `xml:"ID,attr"`
	NumCode   string   `xml:"NumCode"`
	CharCode  string   `xml:"CharCode"`
	Nominal   int      `xml:"Nominal"`
	Name      string   `xml:"Name"`
	Value     string   `xml:"Value"`
	VunitRate string   `xml:"VunitRate"`
}

func newDecoder(input io.Reader) *xml.Decoder {
	d := xml.NewDecoder(input)
	d.CharsetReader = func(charset string, r io.Reader) (io.Reader, error) {
		if strings.EqualFold(charset, "windows-1251") {
			return charmap.Windows1251.NewDecoder().Reader(r), nil
		}
		return nil, fmt.Errorf("unsupported charset: %q", charset)
	}
	return d
}

func getXmlData(path string) ValCurs {
	xmlFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer func(xmlFile *os.File) {
		err := xmlFile.Close()
		if err != nil {
			panic(err)
		}
	}(xmlFile)

	valCurs := ValCurs{}

	decoder := newDecoder(xmlFile)
	if err := decoder.Decode(&valCurs); err != nil {
		panic(err)
	}

	return valCurs
}

func getConfig() Config {
	configData, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func main() {
	config := getConfig()

	valCurs := getXmlData(config.InputFile)

	fmt.Printf("ValCurs: %v\n", valCurs)
}
