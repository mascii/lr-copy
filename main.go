package main

import (
	"bufio"
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
	if len(os.Args) < 3 {
		fmt.Println("Please provide source and destination directory paths.")
		return
	}

	srcDirPath := os.Args[1]
	dstDirPath := os.Args[2]

	files, err := os.ReadDir(srcDirPath)
	if err != nil {
		log.Fatal(err)
	}

	plan, err := cpplan.GenerateCopyPlan(files, srcDirPath, dstDirPath)
	if err != nil {
		panic(err)
	}
	printPlan(files, plan)

	fmt.Println("-----------------------------------")

	if !confirmContinuation() {
		return
	}

	// ディレクトリを作成しておく処理
	for _, dm := range plan {
		if err := os.MkdirAll(dm.DstDir, 0755); err != nil {
			panic(err)
		}
	}

	for _, file := range files {
		m, ok := plan.FindFilePathMapping(file)
		if !ok {
			continue
		}

		if err := copy(m.SrcFilePath, m.DstFilePath); err != nil {
			fmt.Printf("Failed to copy %s to %s (%v)\n", m.SrcFilePath, m.DstFilePath, err)
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

func copy(from, to string) error {
	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
