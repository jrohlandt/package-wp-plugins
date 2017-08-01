package main

import (
	"fmt"
	"log"
	"io"
	"os"
	"archive/zip"
	"path/filepath"
	"strings"
)

var pluginName string // commandline arg 1 e.g. webinarignition
var pluginVersion string // commandline arg 2 e.g. 1.9.89
var sourcePath string
var targetArchive string

var excludeDirectories = []string{".git"}
var excludeFiles = []string{"dev_readme.md", ".gitignore", ".gitattributes"}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("----------------------------------------")
		fmt.Println(" Package Plugin ")
		fmt.Println(" ERROR: You are missing some arguments.")
		fmt.Println(" Usage: webinarignition 1.9.89")
		fmt.Println("----------------------------------------")		
		return
	}

	pluginVersion = os.Args[2]

	switch pluginName := os.Args[1]; pluginName {
		case "webinarignition":
			sourcePath = "/home/jeandre/code/wp_test/wp-content/plugins/webinarignition"
			err := os.Mkdir("/home/jeandre/backups/projects/webinarignition/versions/" + pluginVersion, 0755 )
			if err != nil {
				if os.IsExist(err) {
					fmt.Println("ERROR: Directory already exists!")
				}
				fmt.Println(err)
				return			
			}
			targetArchive = "/home/jeandre/backups/projects/webinarignition/versions/" + pluginVersion + "/webinarignition.zip"
		default:
			fmt.Println("Invalid command. Usage: webinarignition 1.9.89")
	}

	err := zipit(sourcePath, targetArchive)
	if err != nil {
		fmt.Println("-----------------------ERROR-----------------------")	
		log.Fatal(err)
	} else {
		fmt.Println("-----------------------Copy Finish-----------------------")	
	}
}

func directoryShouldBeExcluded(fileinfo os.FileInfo) bool {
	for _, exDirName := range excludeDirectories {
		if (exDirName == fileinfo.Name() && fileinfo.IsDir()) {
			return true
		}
	}
	return false
}

func fileShouldBeExcluded(fileinfo os.FileInfo) bool {

	for _, exFileName := range excludeFiles {
		if (exFileName == fileinfo.Name() && ! fileinfo.IsDir()) {
			return true
		}
	}
	return false 	
}

// zipit http://blog.ralch.com/tutorial/golang-working-with-zip/
func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// fmt.Println(info.Name(), directoryShouldBeExcluded(info.Name()))
		if directoryShouldBeExcluded(info) {
			return filepath.SkipDir
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		if fileShouldBeExcluded(info) {
			return nil
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}