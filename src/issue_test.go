package main

import "testing"

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
