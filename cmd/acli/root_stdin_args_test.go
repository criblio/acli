package acli

import (
	"strings"
	"testing"
)

func TestResolveCLIArgsPassthrough(t *testing.T) {
	in := strings.NewReader(`["jira","issue","list"]`)
	args, err := resolveCLIArgs([]string{"jira", "issue", "list"}, in)
	if err != nil {
		t.Fatalf("resolveCLIArgs returned error: %v", err)
	}

	if len(args) != 3 || args[0] != "jira" || args[1] != "issue" || args[2] != "list" {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestResolveCLIArgsFromStdin(t *testing.T) {
	in := strings.NewReader(`["jira","issue","list"]`)
	args, err := resolveCLIArgs([]string{"-"}, in)
	if err != nil {
		t.Fatalf("resolveCLIArgs returned error: %v", err)
	}

	if len(args) != 3 || args[0] != "jira" || args[1] != "issue" || args[2] != "list" {
		t.Fatalf("unexpected stdin args: %#v", args)
	}
}

func TestReadArgsFromStdinEmptyInput(t *testing.T) {
	_, err := readArgsFromStdin(strings.NewReader(" \n\t "))
	if err == nil {
		t.Fatal("expected an error for empty stdin input")
	}
}

func TestReadArgsFromStdinInvalidJSON(t *testing.T) {
	_, err := readArgsFromStdin(strings.NewReader(`{"arg":"jira"}`))
	if err == nil {
		t.Fatal("expected an error for non-array JSON input")
	}
}

func TestReadArgsFromStdinNonStringArray(t *testing.T) {
	_, err := readArgsFromStdin(strings.NewReader(`[1,2,3]`))
	if err == nil {
		t.Fatal("expected an error for non-string array input")
	}
}
