package main

import (
	"os"
	"fmt"
	// "time"
	// "io/ioutil"
	"bufio"
	"path/filepath"
	"strings"
	"sort"
	"github.com/skratchdot/open-golang/open"
)

var s = struct {
	namesFileName string // file name of the file containing the file names
	pathesFileName string // file name of the file containing the file pathes
	renameFilesPath string // where to put the names and pathes file. if empty, it tries to figure it out
} {
	namesFileName: "_names.txt",
	pathesFileName: "_pathes",
	renameFilesPath: "",
}

func main() {
	
	/* pathes */
	namesFilePath := ""
	pathesFilePath := ""
	
	/* prepare pathes */
	/* figure out the absolut file pathes */
	pathes := make([]string, 0)
	for _, path := range os.Args[1:] {
		if absolutePath, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absolutePath); err == nil {
				/* valid absolute path */
				pathes = append(pathes, absolutePath)
			} else {
				/* no valid file */
				fmt.Println("Skipping:", path)
				fmt.Println(err)
			}
		} else {
			/* could not get absolute path */
			fmt.Println("Skipping:", path)
		}
	}
	// fmt.Println("Pathes:", len(pathes), pathes)
	
	/* check for names and pathes files */
	for _, path := range pathes {
		if filepath.Base(path) == s.namesFileName {
			namesFilePath = path
		} else if filepath.Base(path) == s.pathesFileName {
			pathesFilePath = path
		}
	}
	if namesFilePath != "" && pathesFilePath == "" {
		pathesFilePath = filepath.Dir(namesFilePath) + string(os.PathSeparator) + s.pathesFileName
	}
	if namesFilePath == "" && pathesFilePath != "" {
		namesFilePath = filepath.Dir(pathesFilePath) + string(os.PathSeparator) + s.namesFileName
	}
	
	/* if names and pathes file given, rename */
	if namesFilePath != "" && pathesFilePath != "" {
		fmt.Println("Renaming:")
		
		/* rename files */
		if err := rename(namesFilePath, pathesFilePath); err != nil {
			panic(err)
		}
		
		/* delete names and pathes files */
		if err := os.Remove(namesFilePath); err != nil {
			panic(err)
		}
		if err := os.Remove(pathesFilePath); err != nil {
			panic(err)
		}
		
	} else if len(pathes) > 0 {
		fmt.Println("Preparing renaming:")
		
		generatedNamesFilePath, _, err := prepare(pathes)
		if err != nil {
			panic(err)
		}
		
		/* open file */
		open.Run(generatedNamesFilePath)
		
	}
	
	// time.Sleep(10 * time.Second)
	
}

func prepare(pathes []string) (string, string, error) {
	
	/* sort pathes */
	sort.Strings(pathes)
	
	/* get directory to put names and pathes files in */
	renameFilesPath := s.renameFilesPath
	for _, path := range pathes {
		if (renameFilesPath != "") {
			break
		}
		renameFilesPath = filepath.Dir(path)
	}
	if (renameFilesPath == "") {
		renameFilesPath = "."
	}
	
	/* prepare file pathes */
	namesFilePath := renameFilesPath + string(os.PathSeparator) + s.namesFileName
	pathesFilePath := renameFilesPath + string(os.PathSeparator) + s.pathesFileName
	
	/* create names file */
	namesFile, err := os.Create(namesFilePath)
	if err != nil {
		return namesFilePath, pathesFilePath, err
	}
	defer namesFile.Close()
	
	/* create pathes file */
	pathesFile, err := os.Create(pathesFilePath)
	if err != nil {
		return namesFilePath, pathesFilePath, err
	}
	defer pathesFile.Close()
	
	/* go through pathes and write files */
	for _, path := range pathes {
	
		// fmt.Println(path)
		
		/* prepare values */
		info, _ := os.Stat(path)
		name := filepath.Base(path)
		extension := ""
		basename := name
		if !info.IsDir() {
			extension = filepath.Ext(name)
			basename = strings.TrimSuffix(basename, extension)
		}
		
		/* write values */
		pathesFile.WriteString(path + "\n")
		namesFile.WriteString(basename + "\n")
		
		fmt.Println(path)
		
	}
	
	return namesFilePath, pathesFilePath, nil
}

func rename(namesFilePath, pathesFilePath string) error {
	
	/* open names file */
	namesFile, err := os.Open(namesFilePath)
	if err != nil {
		return err
	}
	defer namesFile.Close()
	
	/* open pathes file */
	pathesFile, err := os.Open(pathesFilePath)
	if err != nil {
		return err
	}
	defer pathesFile.Close()
	
	/* create scanners */
	namesScanner := bufio.NewScanner(namesFile)
	pathesScanner := bufio.NewScanner(pathesFile)
	
	/* read line for line and rename files */
	for namesScanner.Scan() && pathesScanner.Scan() {
		
		/* prepare renaming */
		path := pathesScanner.Text()
		fmt.Println(path)
		
		/* get file info */
		info, err := os.Stat(path)
		if err != nil {
			fmt.Println(" --X error getting file info")
			continue
		}
		
		/* construct new name */
		basename := namesScanner.Text()
		newname := basename
		if !info.IsDir() {
			newname += filepath.Ext(path)
		}
		
		/* check if renaming required */
		if newname == filepath.Base(path) {
			fmt.Println(" == name unchanged")
			continue
		}
		
		/* construct new path */
		newpath := filepath.Dir(path) + string(os.PathSeparator) + newname
		
		/* rename */
		if err := os.Rename(path, newpath); err != nil {
			fmt.Println(" --X error renaming")
			continue
		}
		
		fmt.Println(" --> " + newname)
		
	}
	
	if err := namesScanner.Err(); err != nil {
		return err
	}
	if err := pathesScanner.Err(); err != nil {
		return err
	}
	
	return nil
}
