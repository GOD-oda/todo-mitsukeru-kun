package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var commentFormat = map[string]string{
	".go":  "// TODO",
	".php": "// TODO",
	".rb":  "# TODO",
}

type Comment struct {
	Body       string
	LineNumber int
}

func processFile(filePath string, todoPrefix string) ([]Comment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fmt.Println("Processing file:", filePath)

	var commentLines []Comment
	lineNumber := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), todoPrefix) {
			commentLines = append(commentLines, Comment{Body: line, LineNumber: lineNumber})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return commentLines, nil
}

func visitFile(fp string, fi os.DirEntry, err error) error {
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if fi.IsDir() {
		return nil
	}

	ext := strings.ToLower(filepath.Ext(fp))
	todoPrefix := commentFormat[ext]
	comments, err := processFile(fp, todoPrefix)
	if err != nil {
		fmt.Println("Error processing file:", err)
	} else {
		fmt.Printf("CommentCount: %d\n", len(comments))
		for _, comment := range comments {
			fmt.Printf("%d: %s\n", comment.LineNumber, strings.TrimSpace(comment.Body))
		}
	}
	fmt.Println("--------------")

	return nil
}

type Params struct {
	GithubToken string
}

func GetParams() Params {
	token := os.Getenv("INPUT_GITHUB_TOKEN")
	fmt.Println(token)

	return Params{GithubToken: token}
}

func main() {
	params := GetParams()
	if params.GithubToken == "" {
		fmt.Println("Github token not found. Set the github_token environment variable.")
		os.Exit(1)
	}

	// actionのwithで指定するから多分いらない
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory>")
		os.Exit(1)
	}

	targetDirectory := os.Args[1]

	err := filepath.WalkDir(targetDirectory, visitFile)
	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}
