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
