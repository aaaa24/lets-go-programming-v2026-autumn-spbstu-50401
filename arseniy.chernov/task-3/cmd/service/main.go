package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v4"
	"golang.org/x/text/encoding/charmap"
)

type Float float64

func (f *Float) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string

	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	s = strings.Replace(strings.TrimSpace(s), ",", ".", -1)

	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*f = Float(f64)
	return nil
}

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
	NumCode   int      `xml:"NumCode"`
	CharCode  string   `xml:"CharCode"`
	Nominal   int      `xml:"Nominal"`
	Name      string   `xml:"Name"`
	Value     Float    `xml:"Value"`
	VunitRate Float    `xml:"VunitRate"`
}

type ValuteOutput struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
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

func readXml(path string) ValCurs {
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

func writeJson(path string, data []ValuteOutput) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		panic(err)
	}
}

func createValuteOutput(valCurs *ValCurs) []ValuteOutput {
	var result []ValuteOutput

	for _, v := range valCurs.ValuteItems {
		item := ValuteOutput{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    float64(v.Value),
		}
		result = append(result, item)
	}

	return result
}

func main() {
	config := getConfig()

	valCurs := readXml(config.InputFile)

	sort.Slice(valCurs.ValuteItems, func(i, j int) bool {
		return valCurs.ValuteItems[i].Value > valCurs.ValuteItems[j].Value
	})

	valuteOutput := createValuteOutput(&valCurs)
	writeJson(config.OutputFile, valuteOutput)
}
