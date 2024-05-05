package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mascii/lr-copy/cpplan"
)

// main
func main() {
	log.SetFlags(0) // disable timestamps

	srcDirPath := flag.String("src", "", "Source directory path")
	dstBaseDirPath := flag.String("dst", "", "Destination directory base path")
	overwrite := flag.Bool("overwrite", false, "Overwrite existing files")
	separate := flag.Bool("separate", true, "Separate directory excepting JPEG by file type (e.g. ORF, ARW, etc.)")
	fallback := flag.Bool("fallback", false, "If the date of shooting is unknown, fall back to the file's modification time")

	flag.Parse()

	if *srcDirPath == "" || *dstBaseDirPath == "" {
		log.Print("Please provide source and destination directory paths.\n\n")
		flag.Usage()
		return
	}

	files, err := os.ReadDir(*srcDirPath)
	if err != nil {
		log.Fatal(err)
	}
	files = slices.DeleteFunc(files, func(f os.DirEntry) bool {
		return f.IsDir() || f.Name()[0] == '.'
	})

	cfg := cpplan.NewGenerateCopyPlanConfig(*srcDirPath, *dstBaseDirPath, *separate, *fallback)
	plan := cpplan.GenerateCopyPlan(files, cfg)
	if len(plan) == 0 {
		fmt.Println("No files to copy.")
		return
	}

	printPlan(plan)
	printDivider()
	if !confirmContinuation() {
		return
	}

	stats := struct {
		copied  uint
		skipped uint
		failed  uint
	}{}

	for _, m := range plan {
		skipped, err := copyFile(m.SrcFilePath, m.DstFilePath, *overwrite)
		if err != nil {
			log.Printf("Failed to copy %s to %s (%v)\n", m.SrcFilePath, m.DstFilePath, err)
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
	if stats.failed > 0 {
		os.Exit(1)
	}
}

func printPlan(plan []*cpplan.FilePathMapping) {
	for _, m := range plan {
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
	text = strings.Replace(text, "\r\n", "", -1) // for Windows
	text = strings.Replace(text, "\n", "", -1)

	if strings.ToLower(text) != "y" {
		fmt.Println("Process aborted by the user.")
		return false
	}

	return true
}

func copyFile(from, to string, overwrite bool) (skipped bool, err error) {
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
	if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
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
