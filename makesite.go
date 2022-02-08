package main

import (
	"errors"
	"flag"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

type Page struct {
	Title   string
	Content string
	name    string // private, the name of the file.
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func walkDir(dirPath string) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			err = convertToHTML(f.Name())
			if err != nil {
				panic(err)
			}
		} else {
			walkDir(dirPath + "/" + f.Name())
		}
	}
}

func convertToHTML(filePath string) error {
	page := Page{}
	page.name = filePath
	page.Title = "Hello World"

	content, err := os.ReadFile(filePath)

	if err == nil { // The docstring tells us this is good.
		page.Content = string(content)
	} else {
		return err
	}

	newname := strings.TrimSpace(strings.Split(filePath, ".")[0] + ".html")
	println("Converting", filePath, "to", newname)
	f, err := os.Create("out/" + newname)
	if err != nil {
		return err
	}
	t := template.Must(template.New("page.tmpl").ParseGlob("*.tmpl"))
	err = t.Execute(f, page)
	if err != nil {
		return err
	}

	println("Built page:", newname)

	return nil
}

func main() {
	help := flag.Bool("help", false, "Show help")
	file := flag.String("file", "", "the file to read")
	dir := flag.String("dir", "", "the directory of files to read")

	flag.Parse()

	if *help { // If the user needs help, show it then exit.
		flag.PrintDefaults()
		return // This is the same as exit(0)
	}

	if *file == "" && *dir == "" {
		// TODO: Switch away from panic
		flag.Usage()
		panic("No file or directory specified")
	}

	if *file != "" { // We should be only operating on one file.
		if !pathExists(*file) {
			// TODO: Switch away from panic
			panic("File does not exist")
		}

		err := convertToHTML(*file)
		if err != nil {
			panic(err)
		}

	}

	if *dir != "" { // We should be operating on the given directory
		// Learn how to walk a directory, and operate on each file.
		// We need to make sure the directory exists (IsDir)
		fileDir, _ := os.Stat(*dir)
		if !pathExists(*dir) {
			// GTFO
			panic("Directory does not exist")
		}

		if !fileDir.IsDir() {
			// TODO: Switch away from panic
			panic("Provided directory does not exist")
		}

		walkDir(*dir)
	}

	/*
		println("file:", *file)
		println("dir:", *dir)
	*/

}
