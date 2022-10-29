package xml_parser

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
)

func GetSuitesFromXmlFile(name string) (suites JUnitTestSuites) {
	xmlFile, err := os.Open(name)
	if err != nil {
		log.Panicf("unable to open file: %s", err.Error())
	}

	defer func() {
		_ = xmlFile.Close()
	}()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	err = xml.Unmarshal(byteValue, &suites)
	if err != nil {
		log.Panicf("error unmarshalling xml file: %s", err.Error())
	}
	return
}
