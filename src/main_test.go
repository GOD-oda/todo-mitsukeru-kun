package main

import (
	"os"
	"reflect"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("INPUT_GITHUB_TOKEN", "github token")
	os.Setenv("INPUT_TARGET_DIR", "target dir")
	os.Setenv("INPUT_ISSUE_LABELS", "")
	expected := Params{GithubToken: "github token", TargetDir: "target dir", issueLabels: []IssueLabel{}}
	actual := getEnv()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

func TestCommentMakeLine(t *testing.T) {
	os.Setenv("GITHUB_REPOSITORY", "owner/repo")
	c := Comment{Body: "TODO: TODO sample", LineNumber: 1, FilePath: "src/main.go"}
	expected := "[1: TODO: TODO sample](https://github.com/owner/repo/blob/main/src/main.go#L1)\\n\\n"

	actual := c.makeLine()
	if expected != actual {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}
