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
	TargetDir   string
}

func GetParams() Params {
	githubToken := os.Getenv("INPUT_GITHUB_TOKEN")
	targetDir := os.Getenv("INPUT_TARGET_DIR")

	return Params{GithubToken: githubToken, TargetDir: targetDir}
}

func main() {
	fmt.Println("=============== EnvVars ===============")
	envVars := os.Environ()
	for _, envVar := range envVars {
		fmt.Println(envVar)
	}
	fmt.Println("=======================================")

	fmt.Println("=============== Getwd ================")
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	fmt.Println("Current Directory:", currentDir)
	fmt.Println("=======================================")

	params := GetParams()
	if params.GithubToken == "" {
		fmt.Println("Github token not found. Set the github_token environment variable.")
		os.Exit(1)
	}

	err = filepath.WalkDir(params.TargetDir, visitFile)
	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}
