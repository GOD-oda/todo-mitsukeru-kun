package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"gododa/todo-mitsukeru-kun/github"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

func (c Comment) makeIssueBody() string {
	return fmt.Sprintf("%d: %s\\n\\n", c.LineNumber, strings.TrimSpace(c.Body))
}

func processFile(filePath string, todoPrefix string) ([]Comment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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

type CachedItems struct {
	Items []map[string]interface{}
}

var (
	mu          sync.Mutex
	cachedItems CachedItems
	fetched     bool
)

func getIssues() (CachedItems, error) {
	mu.Lock()
	if fetched {
		mu.Unlock()
		return cachedItems, nil
	}
	mu.Unlock()

	token := os.Getenv("INPUT_GITHUB_TOKEN")
	repoName := os.Getenv("GITHUB_REPOSITORY")

	url := fmt.Sprintf("repos/%s/issues?creator=app/github-actions", repoName)
	res, err := github.Get(url, "", github.MakeHeader(token))
	if err != nil {
		return CachedItems{}, err
	}
	defer res.Body.Close()

	var items []map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&items)
	if err != nil {
		return CachedItems{}, err
	}

	mu.Lock()
	cachedItems = CachedItems{Items: items}
	fetched = true
	mu.Unlock()

	return cachedItems, err
}

func saveIssue(filePath string, comments []Comment) {
	if comments == nil || len(comments) < 1 {
		return
	}

	// TODO: use getEnv()
	token := os.Getenv("INPUT_GITHUB_TOKEN")
	repoName := os.Getenv("GITHUB_REPOSITORY")
	issueTitle := fmt.Sprintf("[todo-mitsukeru-kun] %s", filePath)
	issueBody := "<details>\\n<summary>Todo Comments</summary>\\n\\n\\n"
	for _, comment := range comments {
		issueBody += comment.makeIssueBody()
	}
	issueBody += "</details>\\n"

	url := fmt.Sprintf("https://api.github.com/repos/%s/issues", repoName)
	jsonData := fmt.Sprintf(`{"title": "%s", "body": "%s"}`, issueTitle, issueBody)

	_, err := getIssues()
	if err != nil {
		fmt.Println("Error getting issues:", err)
		return
	}

	var issueId float64
	for _, issue := range cachedItems.Items {
		if issue["title"] == issueTitle {
			issueId = issue["number"].(float64)
			break
		}
	}

	var httpMethod string
	if issueId != 0 {
		url = fmt.Sprintf("%s/%d", url, int(issueId))
		httpMethod = "PATCH"
	} else {
		httpMethod = "POST"
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewBufferString(jsonData))
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

	saveIssue(fp, comments)

	return nil
}

type Params struct {
	GithubToken string
	TargetDir   string
}

func getEnv() Params {
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
	params := getEnv()
	err := filepath.WalkDir(params.TargetDir, visitFile)
	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}
