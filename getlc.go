// This is code to get all the MARC records LC has.  It is not good example
// code for Go, it's just what I crammed together quickly to get the job done.
// Don't use this to learn Go.

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

const (
	Server        = "https://chroniclingamerica.loc.gov/"
	TitlesPattern = "search/titles/results/?format=json&page=%d"
	MARCPattern   = "lccn/%s/marc.xml"
)

type Title struct {
	LCCN string `json:"lccn"`
}

type SearchResult struct {
	EndIndex   int     `json:"endIndex"`
	TotalItems int     `json:"totalItems"`
	Titles     []Title `json:"items"`
}

func setLastPageRead(p int) {
	var f, err = os.Create(".lastpage")
	if err != nil {
		log.Fatalf("Unable to write to .lastpage: %s", err)
	}
	defer f.Close()

	err = binary.Write(f, binary.LittleEndian, int32(p))
	if err != nil {
		log.Fatalf("Unable to serialize to .lastpage: %s", err)
	}
}

func getLastPageRead() int {
	var f, err = os.Open(".lastpage")
	if err != nil {
		log.Printf("Unable to read .lastpage; defaulting to page 0 (%s)", err)
		return 0
	}
	defer f.Close()

	var p int32
	err = binary.Read(f, binary.LittleEndian, &p)
	if err != nil {
		log.Printf("Unable to read bytes from .lastpage; defaulting to page 0 (%s)", err)
		return 0
	}

	return int(p)
}

func getSearchPage(p int) SearchResult {
	var url = Server + fmt.Sprintf(TitlesPattern, p+1)
	log.Printf("GET %s", url)
	var response, err = http.Get(url)
	if err != nil {
		log.Fatalf("Error searching page %d: %s", p, err)
	}
	defer response.Body.Close()

	setLastPageRead(p)
	var buf = &bytes.Buffer{}
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		log.Fatalf("Error reading from chroniclingamerica.loc.gov: %s", err)
	}

	var r SearchResult
	ioutil.WriteFile(".searchdebug", buf.Bytes(), 0644)
	err = json.Unmarshal(buf.Bytes(), &r)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}

	return r
}

func getMARC(lccn string) {
	log.Printf("INFO - LCCN: %#v", lccn)

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

	var dir = path.Join("marc", lccn)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		log.Printf("WARN - Error creating directory: %s", err)
		return
	}

	var f *os.File
	f, err = os.Create(path.Join(dir, "marc.xml"))
	if err != nil {
		log.Printf("WARN - Error opening file for writing marc.xml: %s", err)
		return
	}
	defer f.Close()
	f.Write(xml.Bytes())
}

func main() {
	// Start searching titles from last page pulled or page 1
	var p = getLastPageRead()
	var sr = getSearchPage(p)

	log.Printf("INFO - processing %d titles", sr.TotalItems)

	for sr.TotalItems > sr.EndIndex {
		for _, title := range sr.Titles {
			// Let's not DOS Chronicling America
			time.Sleep(time.Millisecond * 500)
			getMARC(title.LCCN)
		}
		p++
		sr = getSearchPage(p)
	}
}
