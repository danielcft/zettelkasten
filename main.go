package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func add(editor string, zetDir string) {
	id := time.Now().Format("0601021504")
	filepath := getFilePath(zetDir, id)
	createFile(filepath, "# "+id+" TITLE\n\ntags: #\n\n")
	edit(editor, filepath)
}

func createFile(filepath string, text string) {
	f, _ := os.Create(filepath)
	f.WriteString(text)
	f.Close()
}

func edit(editor string, filepath string) {
	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func delete(zetDir string, id string) {
	err := os.Remove(getFilePath(zetDir, id))
	if err != nil {
		fmt.Println(err)
	}
}

func view(zetDir string, id string) {
	year := string(id[0]) + string(id[1])
	printFile(zetDir+"/"+year+"/"+id+".md", -1)
}

func grep(zetDir string, grepRegexp string) {
	matches := []string{}

	filepath.Walk(zetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if info.IsDir() {
			return nil
		}

		// -E extended regexp, -q silent (use status code for match)
		cmd := exec.Command("grep", "-Eq", grepRegexp, path)
		cmd.Start()
		e := cmd.Wait()
		// e is nil if status code is 0, meaning there is match
		if e == nil {
			matches = append(matches, path)
		}
		return err
	})

	for _, match := range matches {
		fmt.Println(strings.Repeat("=", 80))
		printFile(match, -1)
	}

}

func printFile(filePath string, linesToPrint int) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() && linesToPrint != 0 {
		fmt.Println(scanner.Text())
		linesToPrint--
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func list(zetDir string) {
	filepath.Walk(zetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if info.IsDir() {
			return err
		}

		printFile(path, 1)

		return err
	})

}

func help() {
	help := ` 
usage: zet [CMD]
  add        Add a zet
  view       View zet
  edit       Edit a zet
  delete     Delete a zet
  grep       Grep for keywords
  list       List zets
`
	fmt.Printf("%s", help)
}

func getZetDir() string {
	zetDir := os.Getenv("ZET_DIR")
	if zetDir == "" {
		fmt.Println("ZET_DIR not defined")
		os.Exit(1)
	}
	return zetDir
}

func getEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "vi"
	}
	return editor
}

func getFilePath(zetDir string, id string) string {
	year := string(id[0]) + string(id[1])
	return zetDir + "/" + year + "/" + id + ".md"
}

func showErrorMissingParameter(command string, parameter string) {
	fmt.Printf("%s: missing parameter '%s'\n", command, parameter)
	os.Exit(1)
}

func main() {
	if len(os.Args) <= 1 {
		help()
		os.Exit(1)
	}
	zetDir := getZetDir()
	editor := getEditor()
	command := os.Args[1]
	switch {
	case command == "add" || command == "a":
		add(editor, zetDir)
	case command == "view" || command == "v":
		if len(os.Args) < 3 {
			showErrorMissingParameter("view", "id")
		}
		id := os.Args[2]
		view(zetDir, id)
	case command == "edit" || command == "e":
		if len(os.Args) < 3 {
			showErrorMissingParameter("edit", "id")
		}
		id := os.Args[2]
		edit(editor, getFilePath(zetDir, id))
	case command == "delete" || command == "d":
		if len(os.Args) < 3 {
			showErrorMissingParameter("delete", "id")
		}
		id := os.Args[2]
		delete(zetDir, id)
	case command == "list" || command == "ls":
		list(zetDir)
	case command == "grep" || command == "g":
		if len(os.Args) < 3 {
			showErrorMissingParameter("grep", "regexp")
		}
		regexp := os.Args[2]
		grep(zetDir, regexp)
	default:
		fmt.Printf("Wrong CMD '%s'\n", command)
		os.Exit(1)
	}
}
