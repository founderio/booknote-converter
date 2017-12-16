package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

var inputFilename string
var outputFilename string
var year string
var month string
var day string

var rYear = regexp.MustCompile("^[0-9]{4}\\s*:$")
var rDate = regexp.MustCompile("^[0-9]{1,2}[,.][0-9]{1,2}[,.]?\\s*:?$")
var rTransaction = regexp.MustCompile("^(?P<value>[-+]?[0-9]+([,.][0-9]+)?)\\s*(?P<payee>[^ 0-9,].+)$")

func main() {

	flag.StringVar(&inputFilename, "input", "input.txt", "Input file for noted transactions")
	flag.StringVar(&outputFilename, "output", "output.csv", "Output file for generated csv")
	flag.StringVar(&year, "year", "", "The initial year to use, leave blank if the year is specified as 'yyyy:' in the file")
	flag.Parse()

	if inputFilename == "" {
		log.Fatalln("No input file specified.")
	}
	if outputFilename == "" {
		log.Fatalln("No output file specified.")
	}

	inFile, err := os.OpenFile(inputFilename, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalln("Error opening input file", inputFilename, "\n", err.Error())
	}
	defer func() {
		err = inFile.Close()
		if err != nil {
			log.Println("Error closing input file:", err.Error())
		}
	}()

	outFile, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalln("Error creating/opening output file:", err.Error())
	}
	defer func() {
		err = outFile.Close()
		if err != nil {
			log.Println("Error closing output file:", err.Error())
		}
	}()

	csvWriter := csv.NewWriter(outFile)
	csvWriter.Comma = ';'

	var writeError error
	write := func(s ...string) {
		err := csvWriter.Write(s)
		if err != nil && writeError == nil {
			writeError = err
		}
	}

	// Write header
	write("date", "payment", "info", "payee", "memo", "amount", "category", "tags")

	lineNumber := 0

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Is it a new year?
		if rYear.MatchString(line) {
			year = line[:4]
			continue
		}

		// We need a year number, complain if we don't have one
		if year == "" {
			log.Fatal("No year given, and first non-empty line is not a year number: ", lineNumber, " ", line)
		}

		// Is it a new date?
		if rDate.MatchString(line) {
			line = strings.Replace(line, ",", ".", -1)
			split := strings.Split(line, ".")
			day = split[0]
			month = split[1]
			continue
		}

		// We need a date, complain if we don't have one
		if day == "" || month == "" {
			log.Fatal("No date line found before the first transaction: ", lineNumber, " ", line)
		}

		// Split a normal line - we expect a transaction of some kind
		values := rTransaction.FindStringSubmatch(line)

		if len(values) < 2 {
			log.Fatal("Failed to correctly parse line: ", lineNumber, " ", line)
		}

		// Transform results into an adressable map for convenience
		result := make(map[string]string)
		for i, name := range rTransaction.SubexpNames() {
			if i != 0 {
				result[name] = values[i]
			}
		}

		tValue := strings.Replace(result["value"], ",", ".", -1)
		//tValue = strings.TrimPrefix(tValue, "+")
		tPayee := result["payee"]

		if tValue == "" || tPayee == "" {
			log.Fatal("Failed to correctly parse line: ", lineNumber, " ", line)
		}

		fValue, err := decimal.NewFromString(tValue)
		if err != nil {
			log.Fatal("Failed to correctly parse value in line: ", lineNumber, " ", line)
		}
		if !strings.HasPrefix(tValue, "+") {
			fValue = fValue.Neg()
		}

		tDate := fmt.Sprintf("%s-%s-%s", lpad(year, "0", 4), lpad(month, "0", 2), lpad(day, "0", 2))

		// "date", "payment", "info", "payee", "memo", "amount", "category", "tags"
		write(tDate, "3", "", tPayee, "", fValue.StringFixed(2), "", "")
		// payment = 3 = Cash
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading input file:", err.Error())
	}

	// Write any buffered data to the underlying writer (standard output).
	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		log.Fatal("Error writing to output file:", err.Error())
	}
}

func lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}
