package acli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/chinmaymk/acli/internal/config"
	"github.com/chinmaymk/acli/internal/jira"
	"github.com/spf13/cobra"
)

func getProfile(cmd *cobra.Command) (config.Profile, error) {
	profileName, _ := cmd.Flags().GetString("profile")
	cfg, err := config.Load()
	if err != nil {
		return config.Profile{}, fmt.Errorf("loading config: %w", err)
	}
	return cfg.GetProfile(profileName)
}

func getJiraClient(cmd *cobra.Command) (*jira.Client, error) {
	profile, err := getProfile(cmd)
	if err != nil {
		return nil, err
	}
	return jira.NewClient(profile)
}

func getBitbucketClient(cmd *cobra.Command) (*bitbucket.Client, error) {
	profile, err := getProfile(cmd)
	if err != nil {
		return nil, err
	}
	return bitbucket.NewClient(profile)
}

func newTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
}

func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// truncate truncates a string to maxLen characters, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// printIssueRow prints a single issue row to a tabwriter (shared by board/sprint/search).
func printIssueRow(w *tabwriter.Writer, issue jira.IssueDetailed) {
	issueType := ""
	if issue.Fields.IssueType != nil {
		issueType = issue.Fields.IssueType.Name
	}
	status := ""
	if issue.Fields.Status != nil {
		status = issue.Fields.Status.Name
	}
	priority := ""
	if issue.Fields.Priority != nil {
		priority = issue.Fields.Priority.Name
	}
	assignee := ""
	if issue.Fields.Assignee != nil {
		assignee = issue.Fields.Assignee.DisplayName
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
		issue.Key, issueType, status, priority, assignee, issue.Fields.Summary)
}
