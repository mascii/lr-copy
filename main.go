package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"github.com/mascii/lr-copy/cpplan"
)

// main
func main() {
	srcDirPath := flag.String("src", "", "Source directory path")
	dstBaseDirPath := flag.String("dst", "", "Destination directory base path")
	overwrite := flag.Bool("overwrite", false, "Overwrite existing files")
	separate := flag.Bool("separate", true, "Separate directory excepting JPEG by file type (e.g. ORF, ARW, etc.)")

	flag.Parse()

	if *srcDirPath == "" || *dstBaseDirPath == "" {
		fmt.Print("Please provide source and destination directory paths.\n\n")
		flag.Usage()
		return
	}

	files, err := os.ReadDir(*srcDirPath)
	if err != nil {
		log.Fatal(err)
	}

	plan := cpplan.GenerateCopyPlan(files, *srcDirPath, *dstBaseDirPath, *separate)
	if plan.HasNoFilesToCopy() {
		fmt.Println("No files to copy.")
		return
	}

	printPlan(files, plan)
	printDivider()
	if !confirmContinuation() {
		return
	}

	stats := struct {
		copied  uint
		skipped uint
		failed  uint
	}{}

	for _, file := range files {
		m, ok := plan.FindFilePathMapping(file)
		if !ok {
			continue
		}

		skipped, err := copy(m.SrcFilePath, m.DstFilePath, *overwrite)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy %s to %s (%v)\n", m.SrcFilePath, m.DstFilePath, err)
			stats.failed++
		} else if skipped {
			fmt.Printf("Skipped %s to %s\n", m.SrcFilePath, m.DstFilePath)
			stats.skipped++
		} else {
			fmt.Printf("Copied %s to %s\n", m.SrcFilePath, m.DstFilePath)
			stats.copied++
		}
	}

	printDivider()
	fmt.Printf("Copied: %d, Skipped: %d, Failed: %d\n", stats.copied, stats.skipped, stats.failed)
}

func printPlan(files []fs.DirEntry, plan cpplan.Plan) {
	for _, file := range files {
		m, ok := plan.FindFilePathMapping(file)
		if !ok {
			continue
		}
		fmt.Printf("%s => %s\n", m.SrcFilePath, m.DstFilePath)
	}
}

func printDivider() {
	fmt.Println("-----------------------------------")
}

func confirmContinuation() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Continue with the process? (y/N): ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	if strings.ToLower(text) != "y" {
		fmt.Println("Process aborted by the user.")
		return false
	}

	return true
}

func copy(from, to string, overwrite bool) (skipped bool, err error) {
	src, err := os.Open(from)
	if err != nil {
		return false, err
	}
	defer src.Close()

	srcInfo, err := src.Stat()
	if err != nil {
		return false, err
	}

	if !overwrite {
		if _, err := os.Stat(to); !os.IsNotExist(err) {
			return true, nil
		}
	}

	// ディレクトリを作成しておく処理
	if err := os.MkdirAll(path.Dir(to), 0755); err != nil {
		return false, err
	}

	dst, err := os.Create(to)
	if err != nil {
		return false, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return false, err
	}

	// Set the timestamp of the destination file to be the same as the source file
	err = os.Chtimes(to, srcInfo.ModTime(), srcInfo.ModTime())
	if err != nil {
		return false, err
	}

	return false, nil
}
