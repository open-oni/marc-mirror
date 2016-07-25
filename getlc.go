package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	Server = "http://chroniclingamerica.loc.gov/"
	NewspaperJSON = "newspapers.json"
	MARCPattern = "lccn/%s/marc.xml"
)

type Newspaper struct {
	LCCN string `json:"lccn"`
}

type Resp struct {
	Newspapers []Newspaper `json:"newspapers"`
}

func main() {
	// Get list of newspapers
	var response, err = http.Get(Server + NewspaperJSON)
	if err != nil {
		log.Fatalf("Error fetching newspapers.json: %s", err)
	}

	var buf = &bytes.Buffer{}
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		log.Fatalf("Error reading from chroniclingamerica.loc.gov: %s", err)
	}

	var r Resp
	json.Unmarshal(buf.Bytes(), &r)

	log.Printf("INFO - processing %d papers", len(r.Newspapers))

	for _, newspaper := range r.Newspapers {
		var lccn = newspaper.LCCN
		getMARC(lccn)
	}
}

func getMARC(lccn string) {
	log.Printf("INFO - LCCN: %#v", lccn)
	os.Mkdir(lccn, 0755)

	var response, err = http.Get(Server + fmt.Sprintf(MARCPattern, lccn))
	if err != nil {
		log.Printf("WARN - Couldn't fetch newspaper: %s", err)
		return
	}
	var xml = &bytes.Buffer{}
	_, err = xml.ReadFrom(response.Body)
	if err != nil {
		log.Printf("WARN - Error reading MARC XML: %s", err)
		return
	}

	var f *os.File
	f, err = os.Create(path.Join(lccn, "marc.xml"))
	if err != nil {
		log.Printf("WARN - Error opening file for writing marc.xml: %s", err)
		return
	}
	defer f.Close()
	f.Write(xml.Bytes())
}
