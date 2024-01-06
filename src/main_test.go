package main

import (
	"os"
	"testing"
)

func TestMakeIssueBody(t *testing.T) {
	c := Comment{Body: "TODO: TODO sample", LineNumber: 1}
	expected := "1: TODO: TODO sample\\n\\n"
	actual := c.makeLine()
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

func TestIssueBodyAdd(t *testing.T) {
	issueBody := IssueBody{}
	issueBody.add("body")
	expected := "body"
	actual := issueBody.Value
	if expected != actual {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

func TestIssueToJson(t *testing.T) {
	issue := Issue{Title: IssueTitle{Value: "title"}, Body: IssueBody{Value: "body"}}
	expected := `{"title": "title", "body": "body"}`
	actual := issue.toJson()
	if expected != actual {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}
