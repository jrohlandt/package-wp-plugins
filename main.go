package main

// source https://gist.github.com/magicianzrh/d29d08f4fdd8b3c2f8ee
import (
	// "fmt"
	"io"
	// "log"
	"os"
	"archive/zip"
	"path/filepath"
	"strings"
)

// var pluginName string
// var pluginVersion string

var excludeDirectories = []string{".git"}
var excludeFiles = []string{"dev_readme.md", ".gitignore", ".gitattributes"}

func main() {

	// pluginName = os.Args[1]
	// pluginVersion = os.Args[2]
	source_path := "/home/jeandre/code/wp_test/wp-content/plugins/webinarignition"
	// target_path := "/home/jeandre/testcopy/wi"

	zipit(source_path, "/home/jeandre/testcopy/wi.zip")
	return
	// err := copy_folder(source_path, target_path)
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Print("copy finish")
	// }

}

// func copy_folder(source string, dest string) (err error) {

// 	sourceinfo, err := os.Stat(source)
// 	if err != nil {
// 		return err
// 	}

// 	// fmt.Println(sourceinfo.Name())
// 	if (directoryShouldBeExcluded(sourceinfo.Name())) {
// 		return
// 	}

// 	// create destination directory
// 	err = os.MkdirAll(dest, sourceinfo.Mode())
// 	if err != nil {
// 		return err
// 	}

// 	directory, _ := os.Open(source)

// 	objects, err := directory.Readdir(-1)


// 	for _, obj := range objects {
		
// 		sourcefilepointer := source + "/" + obj.Name()

// 		destinationfilepointer := dest + "/" + obj.Name()

// 		if obj.IsDir() {
// 			err = copy_folder(sourcefilepointer, destinationfilepointer)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 		} else {
// 			if (fileShouldBeExcluded(obj.Name())) {
// 				continue
// 			}
// 			err = copy_file(sourcefilepointer, destinationfilepointer)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 		}

// 	}
// 	return
// }

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

func copy_file(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
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
		return nil
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