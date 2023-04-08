package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cbroglie/mustache"
)

type File struct {
	XMLName xml.Name `xml:"rdf"`
	Rdf     []Obj    `xml:"metadata"`
}

type Obj struct {
	XMLName    xml.Name `xml:"metadata"`
	Identifier string   `xml:"identifier"`
	Title      string   `xml:"title"`
	Date       string   `xml:"date"`

	Image string `xml:Image`
	Desc  string `xml:"description"`

	Type        string `xml:"type"`
	Format      string `xml:"format"`
	PhysicalMed string `xml:"PhysicalMedium"`
	Location    string `xml:"Location"`

	Subject  string `xml:"subject"`
	Coverage string `xml:"coverage"`
	Spatial  string `xml:"spatial"`
	Temporal string `xml:"temporal"`

	Source      string `xml:"source"`
	Creator     string `xml:"creator"`
	Contributor string `xml:"contributor"`
	Publisher   string `xml:"publisher"`

	Rights   string `xml:"rights"`
	Language string `xml:"language"`

	Relation string `xml:"relation"`
}

var fileArray []string

func main() {

	//copying images into website folder
	oldDir := "./content/media"
	newDir := "./theme/style/media"

	cmd := exec.Command("cp", "--recursive", oldDir, newDir)
	cmd.Run()

	filepath.Walk("./content", VisitFiles)

	for i := 0; i < len(fileArray); i++ {
		//fmt.Println(fileArray[i])
		// opening the file
		xmlFile, err := os.Open(string(fileArray[i]))
		//xmlFile, err := os.Open("AC00002.rdf")
		// if there is an err, it's handled here
		if err != nil {
			fmt.Println("err", err)
		}
		// defer so we can parse it
		defer xmlFile.Close()

		fileBytes, _ := ioutil.ReadAll(xmlFile)
		var item File
		xml.Unmarshal(fileBytes, &item)

		for m := 0; m < len(item.Rdf); m++ {
			//fmt.Println("ITEM: " + item.Rdf[m].Title)
			webpageName := "./theme/pages/" + item.Rdf[m].Identifier + ".html"
			stachio(item.Rdf[m], webpageName)
		}
	}
}

func stachio(entry Obj, pageName string) {
	//template, _ := mustache.ParseFile("item.html.mustache")
	//rendered, _ := mustache.RenderFile("item.html.mustache", entry)
	rendered, _ := mustache.RenderFileInLayout("./theme/item.html.mustache", "./theme/layout.html.mustache", entry)
	ioutil.WriteFile(pageName, []byte(rendered), 0644)
}

func VisitFiles(path string, info os.FileInfo, err error) error {
	// looking through content folder for each item's file
	if err != nil {
		fmt.Println(err) // can't walk here,
		return nil       // but continue walking elsewhere
	}
	if info.IsDir() {
		return nil // not a file.  ignore.
	}

	matched, err := filepath.Match("*.rdf", info.Name())
	if err != nil {
		fmt.Println(err) // malformed pattern
		return err       // this is fatal.
	}
	if matched {
		fileArray = append(fileArray, path)
	}
	return nil
}
