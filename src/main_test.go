package main

import (
	"os"
	"testing"
)

func TestMakeIssueBody(t *testing.T) {
	c := Comment{Body: "TODO: TODO sample", LineNumber: 1}
	expected := "1: TODO: TODO sample\\n\\n"
	actual := c.makeIssueBody()
	if expected != actual {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("INPUT_GITHUB_TOKEN", "github token")
	os.Setenv("INPUT_TARGET_DIR", "target dir")
	expected := Params{GithubToken: "github token", TargetDir: "target dir"}
	actual := getEnv()
	if expected != actual {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}
