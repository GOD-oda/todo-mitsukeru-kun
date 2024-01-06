package main

import "fmt"

type Issue struct {
	Title IssueTitle `json:"title"`
	Body  IssueBody  `json:"body"`
}

func (i Issue) toJson() string {
	return fmt.Sprintf(`{"title": "%s", "body": "%s"}`, i.Title.Value, i.Body.Value)
}

type IssueTitle struct {
	Value string
}

type IssueBody struct {
	Value string
}

func (i *IssueBody) add(value string) {
	i.Value += value
}
