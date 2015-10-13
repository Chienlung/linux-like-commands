package main
import (
	"os"
	"fmt"
	"path"
	"sort"
	"strings"
)

var (
	directoriesCnt = 0
	filesCnt = 0
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

		fmt.Printf("\n%d directories, %d files\n", directoriesCnt, filesCnt)
	}
}

func showDirFile(file *os.File, indent string, isLastDir bool) {
	dirnames, err := file.Readdirnames(0)
	if err != nil {
		return
	}
	dirnames = sortDir(dirnames)
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
			directoriesCnt++
			_, lastPart := path.Split(file.Name())
			if isLastDir {
				fmt.Print(indent + "└── ", lastPart, "\n")
			} else {
				fmt.Print(indent + "├── ", lastPart, "\n")
			}
			if isLastDir {
				showDirFile(file, indent + "    ", isLastDir)
			} else {
				showDirFile(file, indent + "│   ", isLastDir)
			}
		} else {
			filesCnt++
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

type pair struct {
	origi string
	toLower string
}
type pairs []*pair

func (a pairs) Len() int {
	return len(a)
}

func (a pairs) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a pairs) Less(i, j int) bool {
	if strings.Compare(a[i].toLower, a[j].toLower) < 0 {
		return true
	} else {
		return false
	}
}

func sortDir(dirnames []string) []string {
	length := len(dirnames)
	ps := make(pairs, 0, length)
	for _, dirname := range dirnames {
		ps = append(ps, &pair{dirname, strings.ToLower(strings.TrimLeft(dirname, "."))})
	}
	sort.Sort(ps)

	var result []string
	for _, p := range ps {
		result = append(result, p.origi)
	}
	return result
}