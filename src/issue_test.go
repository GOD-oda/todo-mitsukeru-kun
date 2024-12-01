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
	issue := Issue{Title: IssueTitle{Value: "title"}, Body: IssueBody{Value: "body"}, Labels: []IssueLabel{
		{Value: "bug"},
		{Value: "urgent"},
	}}
	expected := `{"body":"body","labels":["bug","urgent"],"title":"title"}`
	actual := issue.toJson()
	if expected != actual {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

func TestIssueToJsonWithZeroIssueLabels(t *testing.T) {
	issue := Issue{Title: IssueTitle{Value: "title"}, Body: IssueBody{Value: "body"}}
	expected := `{"body":"body","title":"title"}`
	actual := issue.toJson()
	if expected != actual {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}
