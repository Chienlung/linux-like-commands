package main
import (
	"os"
	"log"
	"fmt"
	"path"
	"time"
)

const (
	G = 1024 * 1024 * 1024
	M = 1024 * 1024
	K = 1024
)
var filenames = os.Args

func main() {
	if len(filenames) == 1 {
		filenames = append(filenames, ".")
	}

	getAllFiles(filenames)
}

func getAllFiles(filenames []string) {
	for _, filename := range filenames[1:] {
		fmt.Println(filename + ":")
		file, err := os.Open(filename)
		if err != nil {
			log.Println(filename, "is a illegal filename.")
			continue
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			continue
		}
		if fi.IsDir() {
			names, err := file.Readdirnames(0)
			if err != nil {
				return
			}

			for _, name := range names {
				file, err := os.Open(path.Join(filename, name))
				if err != nil {
					log.Println(err)
					continue
				}
				defer file.Close()
				fi, err := file.Stat()
				if err != nil {
					continue
				}
				fmt.Println(fmt.Sprintf("%s\t%s\t%v\t%s", fi.Mode().String(), formatSize(fi.Size()), formatTime(fi.ModTime()), fi.Name()))
			}
		} else {
			fmt.Println(fmt.Sprintf("%s\t%s\t%v\t%s", fi.Mode().String(), formatSize(fi.Size()), formatTime(fi.ModTime()), fi.Name()))
		}
	}
}

func formatSize(size int64) string {
	if size / G > 0 {
		return fmt.Sprintf("%5.1fG", float64(size) / G)
	} else if size / M > 0 {
		return fmt.Sprintf("%5.1fM", float64(size) * 1.0 / M)
	} else if size / K > 0 {
		return fmt.Sprintf("%5.1fK", float64(size) * 1.0 / K)
	} else {
		return fmt.Sprintf("%6d", size)
	}
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%4d/%02d/%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}