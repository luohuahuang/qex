package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// part of the credit to https://github.com/jstemmer/go-junit-report/blob/master/formatter/formatter.go
// JUnitTestSuites is a collection of JUnit test suites.
type JUnitTestSuites struct {
	XMLName xml.Name         `xml:"testsuites"`
	Suites  []JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite is a single JUnit test suite which may contain many
// testcases.
type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Time       string          `xml:"time,attr"`
	Name       string          `xml:"name,attr"`
	Properties []JUnitProperty `xml:"properties>property,omitempty"`
	TestCases  []JUnitTestCase `xml:"testcase"`
}

// JUnitTestCase is a single test case with its result.
type JUnitTestCase struct {
	XMLName     xml.Name          `xml:"testcase"`
	Classname   string            `xml:"classname,attr"`
	Name        string            `xml:"name,attr"`
	Time        string            `xml:"time,attr"`
	SkipMessage *JUnitSkipMessage `xml:"skipped,omitempty"`
	Failure     *JUnitFailure     `xml:"failure,omitempty"`
}

// JUnitSkipMessage contains the reason why a testcase was skipped.
type JUnitSkipMessage struct {
	Message string `xml:"message,attr"`
}

// JUnitProperty represents a key/value pair used to define properties.
type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// JUnitFailure contains data related to a failed test.
type JUnitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

func main() {

	file := flag.String("file", "", `provide file name`)

	flag.Parse()

	if *file == "" {
		flag.PrintDefaults()
		return
	}

	xmlFile, err := os.Open(*file)
	if err != nil {
		fmt.Println(err)
	}

	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var suites JUnitTestSuites
	xml.Unmarshal(byteValue, &suites)

	for i := 0; i < len(suites.Suites); i++ {
		for j := 0; j < len(suites.Suites[i].TestCases); j++ {
			if suites.Suites[i].TestCases[j].Failure != nil {
				//fmt.Println(k.Name)
				if j+1 < len(suites.Suites[i].TestCases) {
					if suites.Suites[i].TestCases[j].Name == suites.Suites[i].TestCases[j+1].Name {
						suites.Suites[i].TestCases[j].SkipMessage = &JUnitSkipMessage{Message: "RERUN: " + suites.Suites[i].TestCases[j].Failure.Contents}
						suites.Suites[i].TestCases[j].Failure = nil
					}
				}
			}
		}
	}

	// to xml
	bytes, err := xml.MarshalIndent(suites, "", "\t")
	err = ioutil.WriteFile(*file+".revised", bytes, 0777)
	if err != nil {
		log.Println(err.Error())
	}
}
