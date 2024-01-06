package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("INPUT_GITHUB_TOKEN", "github token")
	os.Setenv("INPUT_TARGET_DIR", "target dir")
	expected := Params{GithubToken: "github token", TargetDir: "target dir"}
	actual := getEnv()
	if expected != actual {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

func TestCommentMakeLine(t *testing.T) {
	os.Setenv("GITHUB_REPOSITORY", "owner/repo")
	c := Comment{Body: "TODO: TODO sample", LineNumber: 1, FilePath: "src/main.go"}
	expected := "[1: TODO: TODO sample\\n\\n](https://github.com/owner/repo/blob/main/src/main.go#L1)"

	actual := c.makeLine()
	if expected != actual {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}
