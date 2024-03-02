package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/mascii/lr-copy/cpplan"
)

// main
func main() {
	srcDirPath := flag.String("src", "", "Source directory path")
	dstDirPath := flag.String("dst", "", "Destination directory path")
	overwrite := flag.Bool("overwrite", false, "Overwrite existing files")

	flag.Parse()

	if *srcDirPath == "" || *dstDirPath == "" {
		fmt.Print("Please provide source and destination directory paths.\n\n")
		flag.Usage()
		return
	}

	files, err := os.ReadDir(*srcDirPath)
	if err != nil {
		log.Fatal(err)
	}

	plan, err := cpplan.GenerateCopyPlan(files, *srcDirPath, *dstDirPath)
	if err != nil {
		panic(err)
	}
	printPlan(files, plan)

	fmt.Println("-----------------------------------")

	if !confirmContinuation() {
		return
	}

	// ディレクトリを作成しておく処理
	dstDirs := plan.GetDstDirs()
	for dstDir := range dstDirs {
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			panic(err)
		}
	}

	for _, file := range files {
		m, ok := plan.FindFilePathMapping(file)
		if !ok {
			continue
		}

		skipped, err := copy(m.SrcFilePath, m.DstFilePath, *overwrite)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy %s to %s (%v)\n", m.SrcFilePath, m.DstFilePath, err)
		} else if skipped {
			fmt.Printf("Skipped %s to %s\n", m.SrcFilePath, m.DstFilePath)
		} else {
			fmt.Printf("Copied %s to %s\n", m.SrcFilePath, m.DstFilePath)
		}
	}
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
