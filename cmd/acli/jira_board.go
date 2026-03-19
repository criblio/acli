package acli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var jiraBoardCmd = &cobra.Command{
	Use:     "board",
	Aliases: []string{"b"},
	Short:   "Manage boards (Agile API)",
	RunE:    helpRunE,
}

var jiraBoardListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List boards",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		project, _ := defaultProject(cmd)
		boardType, _ := cmd.Flags().GetString("type")
		name, _ := cmd.Flags().GetString("name")
		all, _ := cmd.Flags().GetBool("all")
		jsonFlag := isJSONOutput(cmd)

		result, err := client.GetBoards(startAt, maxResults, project, boardType, name)
		if err != nil {
			return err
		}
		if all {
			for !result.IsLast && len(result.Values) < result.Total {
				next, err := client.GetBoards(startAt+len(result.Values), maxResults, project, boardType, name)
				if err != nil {
					return err
				}
				if len(next.Values) == 0 {
					break
				}
				result.Values = append(result.Values, next.Values...)
				result.IsLast = next.IsLast
			}
		}
		if jsonFlag {
			return outputJSON(result)
		}
		w := newTabWriter()
		_, _ = fmt.Fprintln(w, "ID\tNAME\tTYPE\tPROJECT")
		for _, b := range result.Values {
			project := ""
			if b.Location != nil {
				project = b.Location.ProjectKey
			}
			_, _ = fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", b.ID, b.Name, b.Type, project)
		}
		_ = w.Flush()
		printPaginationHint(cmd, len(result.Values), result.Total)
		return nil
	},
}

var jiraBoardGetCmd = &cobra.Command{
	Use:   "get <board-id>",
	Short: "Get board details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid board ID: %w", err)
		}
		jsonFlag := isJSONOutput(cmd)

		board, err := client.GetBoard(id)
		if err != nil {
			return err
		}
		if jsonFlag {
			return outputJSON(board)
		}
		fmt.Printf("ID:      %d\n", board.ID)
		fmt.Printf("Name:    %s\n", board.Name)
		fmt.Printf("Type:    %s\n", board.Type)
		if board.Location != nil {
			fmt.Printf("Project: %s (%s)\n", board.Location.ProjectName, board.Location.ProjectKey)
		}
		return nil
	},
}

var jiraBoardConfigCmd = &cobra.Command{
	Use:     "config <board-id>",
	Aliases: []string{"configuration"},
	Short:   "Get board configuration",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid board ID: %w", err)
		}
		config, err := client.GetBoardConfiguration(id)
		if err != nil {
			return err
		}
		return outputJSON(config)
	},
}

var jiraBoardIssuesCmd = &cobra.Command{
	Use:   "issues <board-id>",
	Short: "List issues on a board",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid board ID: %w", err)
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		jql, _ := cmd.Flags().GetString("jql")
		all, _ := cmd.Flags().GetBool("all")
		jsonFlag := isJSONOutput(cmd)

		result, err := client.GetBoardIssues(id, startAt, maxResults, jql)
		if err != nil {
			return err
		}
		if all {
			for len(result.Issues) < result.Total {
				next, err := client.GetBoardIssues(id, startAt+len(result.Issues), maxResults, jql)
				if err != nil {
					return err
				}
				if len(next.Issues) == 0 {
					break
				}
				result.Issues = append(result.Issues, next.Issues...)
			}
		}
		if jsonFlag {
			return outputJSON(result)
		}
		w := newTabWriter()
		_, _ = fmt.Fprintln(w, "KEY\tTYPE\tSTATUS\tPRIORITY\tASSIGNEE\tSUMMARY")
		for _, issue := range result.Issues {
			printIssueRow(w, issue)
		}
		_ = w.Flush()
		printPaginationHint(cmd, len(result.Issues), result.Total)
		return nil
	},
}

var jiraBoardBacklogCmd = &cobra.Command{
	Use:   "backlog <board-id>",
	Short: "List backlog issues for a board",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid board ID: %w", err)
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		jql, _ := cmd.Flags().GetString("jql")
		all, _ := cmd.Flags().GetBool("all")
		jsonFlag := isJSONOutput(cmd)

		result, err := client.GetBoardBacklog(id, startAt, maxResults, jql)
		if err != nil {
			return err
		}
		if all {
			for len(result.Issues) < result.Total {
				next, err := client.GetBoardBacklog(id, startAt+len(result.Issues), maxResults, jql)
				if err != nil {
					return err
				}
				if len(next.Issues) == 0 {
					break
				}
				result.Issues = append(result.Issues, next.Issues...)
			}
		}
		if jsonFlag {
			return outputJSON(result)
		}
		w := newTabWriter()
		_, _ = fmt.Fprintln(w, "KEY\tTYPE\tSTATUS\tPRIORITY\tASSIGNEE\tSUMMARY")
		for _, issue := range result.Issues {
			printIssueRow(w, issue)
		}
		_ = w.Flush()
		printPaginationHint(cmd, len(result.Issues), result.Total)
		return nil
	},
}

var jiraBoardSprintsCmd = &cobra.Command{
	Use:   "sprints <board-id>",
	Short: "List sprints for a board",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid board ID: %w", err)
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		state, _ := cmd.Flags().GetString("state")
		all, _ := cmd.Flags().GetBool("all")
		jsonFlag := isJSONOutput(cmd)

		result, err := client.GetBoardSprints(id, startAt, maxResults, state)
		if err != nil {
			return err
		}
		if all {
			for !result.IsLast && len(result.Values) < result.Total {
				next, err := client.GetBoardSprints(id, startAt+len(result.Values), maxResults, state)
				if err != nil {
					return err
				}
				if len(next.Values) == 0 {
					break
				}
				result.Values = append(result.Values, next.Values...)
				result.IsLast = next.IsLast
			}
		}
		if jsonFlag {
			return outputJSON(result)
		}
		w := newTabWriter()
		_, _ = fmt.Fprintln(w, "ID\tNAME\tSTATE\tSTART DATE\tEND DATE")
		for _, s := range result.Values {
			_, _ = fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", s.ID, s.Name, s.State, s.StartDate, s.EndDate)
		}
		_ = w.Flush()
		printPaginationHint(cmd, len(result.Values), result.Total)
		return nil
	},
}

var jiraBoardEpicsCmd = &cobra.Command{
	Use:   "epics <board-id>",
	Short: "List epics for a board",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid board ID: %w", err)
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		all, _ := cmd.Flags().GetBool("all")
		jsonFlag := isJSONOutput(cmd)

		result, err := client.GetBoardEpics(id, startAt, maxResults)
		if err != nil {
			return err
		}
		if all {
			for !result.IsLast && len(result.Values) < result.Total {
				next, err := client.GetBoardEpics(id, startAt+len(result.Values), maxResults)
				if err != nil {
					return err
				}
				if len(next.Values) == 0 {
					break
				}
				result.Values = append(result.Values, next.Values...)
				result.IsLast = next.IsLast
			}
		}
		if jsonFlag {
			return outputJSON(result)
		}
		w := newTabWriter()
		_, _ = fmt.Fprintln(w, "ID\tKEY\tNAME\tDONE\tSUMMARY")
		for _, e := range result.Values {
			_, _ = fmt.Fprintf(w, "%d\t%s\t%s\t%v\t%s\n", e.ID, e.Key, e.Name, e.Done, e.Summary)
		}
		_ = w.Flush()
		printPaginationHint(cmd, len(result.Values), result.Total)
		return nil
	},
}

func init() {
	// board list
	jiraBoardListCmd.Flags().String("project", "", "Filter by project key or ID (uses profile default if not set)")
	jiraBoardListCmd.Flags().String("type", "", "Filter by board type (scrum, kanban)")
	jiraBoardListCmd.Flags().String("name", "", "Filter by board name")
	jiraBoardListCmd.Flags().Int("start-at", 0, "Start index")
	jiraBoardListCmd.Flags().Int("max-results", 50, "Max results per page")
	addAllFlag(jiraBoardListCmd)
	jiraBoardListCmd.Flags().Bool("json", false, "Output as JSON")
	jiraBoardCmd.AddCommand(jiraBoardListCmd)

	// board get
	jiraBoardGetCmd.Flags().Bool("json", false, "Output as JSON")
	jiraBoardCmd.AddCommand(jiraBoardGetCmd)

	// board config
	jiraBoardCmd.AddCommand(jiraBoardConfigCmd)

	// board issues
	jiraBoardIssuesCmd.Flags().Int("start-at", 0, "Start index")
	jiraBoardIssuesCmd.Flags().Int("max-results", 50, "Max results per page")
	jiraBoardIssuesCmd.Flags().String("jql", "", "JQL filter")
	addAllFlag(jiraBoardIssuesCmd)
	jiraBoardIssuesCmd.Flags().Bool("json", false, "Output as JSON")
	jiraBoardCmd.AddCommand(jiraBoardIssuesCmd)

	// board backlog
	jiraBoardBacklogCmd.Flags().Int("start-at", 0, "Start index")
	jiraBoardBacklogCmd.Flags().Int("max-results", 50, "Max results per page")
	jiraBoardBacklogCmd.Flags().String("jql", "", "JQL filter")
	addAllFlag(jiraBoardBacklogCmd)
	jiraBoardBacklogCmd.Flags().Bool("json", false, "Output as JSON")
	jiraBoardCmd.AddCommand(jiraBoardBacklogCmd)

	// board sprints
	jiraBoardSprintsCmd.Flags().Int("start-at", 0, "Start index")
	jiraBoardSprintsCmd.Flags().Int("max-results", 50, "Max results per page")
	jiraBoardSprintsCmd.Flags().String("state", "", "Filter by state (active, closed, future)")
	addAllFlag(jiraBoardSprintsCmd)
	jiraBoardSprintsCmd.Flags().Bool("json", false, "Output as JSON")
	jiraBoardCmd.AddCommand(jiraBoardSprintsCmd)

	// board epics
	jiraBoardEpicsCmd.Flags().Int("start-at", 0, "Start index")
	jiraBoardEpicsCmd.Flags().Int("max-results", 50, "Max results per page")
	addAllFlag(jiraBoardEpicsCmd)
	jiraBoardEpicsCmd.Flags().Bool("json", false, "Output as JSON")
	jiraBoardCmd.AddCommand(jiraBoardEpicsCmd)
}
