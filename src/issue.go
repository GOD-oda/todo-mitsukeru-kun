package main

import (
	"encoding/json"
	"fmt"
)

type Issue struct {
	Title  IssueTitle   `json:"title"`
	Body   IssueBody    `json:"body"`
	Labels []IssueLabel `json:"labels"`
}

func (i Issue) toJson() string {
	issueMap := map[string]interface{}{
		"title": i.Title.Value,
		"body":  i.Body.Value,
	}

	var labelValues []string
	if len(i.Labels) > 0 {
		labelValues = make([]string, len(i.Labels))
		for idx, label := range i.Labels {
			labelValues[idx] = label.Value
		}

		issueMap["labels"] = labelValues
	}

	b, err := json.Marshal(issueMap)
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(b)
}

type IssueTitle struct {
	Value string
}

type IssueBody struct {
	Value string
}

type IssueLabel struct {
	Value string
}

func (i *IssueBody) add(value string) {
	i.Value += value
}
