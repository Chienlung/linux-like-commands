package main
import (
	"os"
	"fmt"
	"path"
	"sort"
)

func main() {
	cmdArgs := os.Args
	if len(cmdArgs) == 1 {
		cmdArgs = append(cmdArgs, ".")
	}

	showAllFiles(cmdArgs[1:])
}

func showAllFiles(filenames []string) {
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			continue
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			continue
		}
		indent := ""

		if fi.IsDir() {
			fmt.Print(indent + fi.Name(), "\n")
			showDirFile(file, indent, false)
		} else {
			showRegularFile(file, indent, false)
		}
	}
}

func showDirFile(file *os.File, indent string, isLastDir bool) {
	dirnames, err := file.Readdirnames(0)
	if err != nil {
		return
	}
	sort.Strings(dirnames)
	var isLastFile bool
	for i, dirname := range dirnames {
		if i + 1 == len(dirnames) {
			isLastFile = true
			isLastDir = true
		} else {
			isLastFile = false
			isLastDir = false
		}
		file, err := os.Open(path.Join(file.Name(),dirname))
		if err != nil {
			continue
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			continue
		}
		if fi.IsDir() {
			_, lastPart := path.Split(file.Name())
			if isLastDir {
				fmt.Print(indent + "└── ", lastPart, "\n")
			} else {
				fmt.Print(indent + "├── ", lastPart, "\n")
			}
			if isLastDir {
				showDirFile(file, "    " + indent, isLastDir)
			} else {
				showDirFile(file, indent + "│   ", isLastDir)
			}
		} else {
			showRegularFile(file, indent, isLastFile)
		}
	}
}

func showRegularFile(file *os.File, indent string, isLastFile bool) {
	//fmt.Print("|")
	_, lastPart := path.Split(file.Name())
	if isLastFile {
		fmt.Print(indent + "└── ", lastPart, "\n")
	} else {
		fmt.Print(indent + "├── ", lastPart, "\n")
	}
}