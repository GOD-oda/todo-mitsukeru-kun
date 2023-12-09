package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
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
	print "hoge"
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

func createIssue(filePath string, comments []Comment) {
	if comments == nil || len(comments) < 1 {
		return
	}

	token := os.Getenv("INPUT_GITHUB_TOKEN")
	repoName := os.Getenv("GITHUB_REPOSITORY")
	issueTitle := fmt.Sprintf("[todo-mitsukeru-kun] %s", filePath)
	issueBody := "<details>\\n<summary>Todo Comments</summary>\\n"
	for _, comment := range comments {
		issueBody += fmt.Sprintf("%d: %s\\n\\n", comment.LineNumber, strings.ReplaceAll(comment.Body, "\t", ""))
	}
	issueBody += "</details>\\n"

	url := fmt.Sprintf("https://api.github.com/repos/%s/issues", repoName)
	jsonData := fmt.Sprintf(`{"title": "%s", "body": "%s"}`, issueTitle, issueBody)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
}

func visitFile(fp string, fi os.DirEntry, err error) error {
	fmt.Println("")

	if err != nil {
		fmt.Println(err)
		return nil
	}

	if fi.IsDir() {
		return nil
	}

	ext := strings.ToLower(filepath.Ext(fp))
	todoPrefix := commentFormat[ext]
	if todoPrefix == "" {
		return nil
	}

	comments, err := processFile(fp, todoPrefix)
	if err != nil {
		fmt.Println("Error processing file:", err)
	} else {
		commentCount := len(comments)
		if commentCount < 1 {
			return nil
		}
	}

	createIssue(fp, comments)

	return nil
}

type Params struct {
	GithubToken string
	TargetDir   string
}

func GetParams() Params {
	githubToken := os.Getenv("INPUT_GITHUB_TOKEN")
	if githubToken == "" {
		fmt.Println("INPUT_GITHUB_TOKEN not found. Set INPUT_GITHUB_TOKEN as environment variable.")
		os.Exit(1)
	}

	targetDir := os.Getenv("INPUT_TARGET_DIR")
	if targetDir == "" {
		fmt.Println("INPUT_TARGET_DIR not found. Set INPUT_TARGET_DIR as environment variable.")
		os.Exit(1)
	}

	return Params{GithubToken: githubToken, TargetDir: targetDir}
}

func main() {
	params := GetParams()
	err := filepath.WalkDir(params.TargetDir, visitFile)
	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}
