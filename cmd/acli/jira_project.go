package acli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var jiraProjectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"p"},
	Short:   "Manage projects",
	RunE:    helpRunE,
}

func init() {
	// --- project list ---
	projectListCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List or search projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}

			jsonOutput := isJSONOutput(cmd)
			query, _ := cmd.Flags().GetString("query")
			maxResults, _ := cmd.Flags().GetInt("max-results")
			startAt, _ := cmd.Flags().GetInt("start-at")
			all, _ := cmd.Flags().GetBool("all")

			result, err := client.SearchProjects(query, startAt, maxResults, "")
			if err != nil {
				return err
			}

			if all {
				for !result.IsLast && len(result.Values) < result.Total {
					next, err := client.SearchProjects(query, startAt+len(result.Values), maxResults, "")
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

			if jsonOutput {
				return outputJSON(result)
			}

			w := newTabWriter()
			_, _ = fmt.Fprintln(w, "KEY\tNAME\tTYPE\tLEAD")
			for _, p := range result.Values {
				lead := ""
				if p.Lead != nil {
					lead = p.Lead.DisplayName
				}
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Key, p.Name, p.ProjectTypeKey, lead)
			}
			_ = w.Flush()
			printPaginationHint(cmd, len(result.Values), result.Total)
			return nil
		},
	}
	projectListCmd.Flags().String("query", "", "Search query to filter projects")
	projectListCmd.Flags().Int("max-results", 50, "Maximum number of results per page")
	projectListCmd.Flags().Int("start-at", 0, "Index of the first result")
	addAllFlag(projectListCmd)
	projectListCmd.Flags().Bool("json", false, "Output as JSON")
	jiraProjectCmd.AddCommand(projectListCmd)

	// --- project get ---
	projectGetCmd := &cobra.Command{
		Use:   "get <project-key>",
		Short: "Get project details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}

			jsonOutput := isJSONOutput(cmd)

			project, err := client.GetProject(args[0], "")
			if err != nil {
				return err
			}

			if jsonOutput {
				return outputJSON(project)
			}

			lead := "N/A"
			if project.Lead != nil {
				lead = project.Lead.DisplayName
			}
			category := "N/A"
			if project.ProjectCategory != nil {
				category = project.ProjectCategory.Name
			}

			fmt.Printf("Key:          %s\n", project.Key)
			fmt.Printf("Name:         %s\n", project.Name)
			fmt.Printf("ID:           %s\n", project.ID)
			fmt.Printf("Type:         %s\n", project.ProjectTypeKey)
			fmt.Printf("Lead:         %s\n", lead)
			fmt.Printf("Description:  %s\n", project.Description)
			fmt.Printf("Category:     %s\n", category)
			fmt.Printf("Style:        %s\n", project.Style)
			fmt.Printf("Simplified:   %v\n", project.Simplified)
			fmt.Printf("Archived:     %v\n", project.Archived)
			fmt.Printf("URL:          %s\n", project.URL)
			return nil
		},
	}
	projectGetCmd.Flags().Bool("json", false, "Output as JSON")
	jiraProjectCmd.AddCommand(projectGetCmd)

	// --- project create ---
	projectCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}

			key, _ := cmd.Flags().GetString("key")
			name, _ := cmd.Flags().GetString("name")
			projectType, _ := cmd.Flags().GetString("type")
			lead, _ := cmd.Flags().GetString("lead")
			description, _ := cmd.Flags().GetString("description")

			body := map[string]interface{}{
				"key":            key,
				"name":           name,
				"projectTypeKey": projectType,
			}
			if lead != "" {
				body["leadAccountId"] = lead
			}
			if description != "" {
				body["description"] = description
			}

			project, err := client.CreateProject(body)
			if err != nil {
				return err
			}

			return outputResult(cmd, "created", project.Key, fmt.Sprintf("Created project %s (ID: %s)", project.Key, project.ID), project)
		},
	}
	projectCreateCmd.Flags().String("key", "", "Project key (required)")
	projectCreateCmd.Flags().String("name", "", "Project name (required)")
	projectCreateCmd.Flags().String("type", "software", "Project type key")
	projectCreateCmd.Flags().String("lead", "", "Lead account ID")
	projectCreateCmd.Flags().String("description", "", "Project description")
	_ = projectCreateCmd.MarkFlagRequired("key")
	_ = projectCreateCmd.MarkFlagRequired("name")
	jiraProjectCmd.AddCommand(projectCreateCmd)

	// --- project update ---
	projectUpdateCmd := &cobra.Command{
		Use:   "update <project-key>",
		Short: "Update a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}

			body := map[string]interface{}{}
			if cmd.Flags().Changed("name") {
				name, _ := cmd.Flags().GetString("name")
				body["name"] = name
			}
			if cmd.Flags().Changed("description") {
				description, _ := cmd.Flags().GetString("description")
				body["description"] = description
			}
			if cmd.Flags().Changed("lead") {
				lead, _ := cmd.Flags().GetString("lead")
				body["leadAccountId"] = lead
			}

			if len(body) == 0 {
				return fmt.Errorf("at least one flag (--name, --description, --lead) must be provided")
			}

			project, err := client.UpdateProject(args[0], body)
			if err != nil {
				return err
			}

			return outputResult(cmd, "updated", project.Key, fmt.Sprintf("Updated project %s", project.Key), project)
		},
	}
	projectUpdateCmd.Flags().String("name", "", "Project name")
	projectUpdateCmd.Flags().String("description", "", "Project description")
	projectUpdateCmd.Flags().String("lead", "", "Lead account ID")
	jiraProjectCmd.AddCommand(projectUpdateCmd)

	// --- project delete ---
	jiraProjectCmd.AddCommand(&cobra.Command{
		Use:   "delete <project-key>",
		Short: "Delete a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}
			if err = client.DeleteProject(args[0]); err != nil {
				return err
			}
			return outputResult(cmd, "deleted", args[0], fmt.Sprintf("Deleted project %s", args[0]), nil)
		},
	})

	// --- project components ---
	jiraProjectCmd.AddCommand(&cobra.Command{
		Use:   "components <project-key>",
		Short: "List project components",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}
			components, err := client.GetProjectComponents(args[0])
			if err != nil {
				return err
			}
			w := newTabWriter()
			_, _ = fmt.Fprintln(w, "ID\tNAME\tLEAD\tASSIGNEE TYPE")
			for _, c := range components {
				lead := ""
				if c.Lead != nil {
					lead = c.Lead.DisplayName
				}
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", c.ID, c.Name, lead, c.AssigneeType)
			}
			return w.Flush()
		},
	})

	// --- project versions ---
	jiraProjectCmd.AddCommand(&cobra.Command{
		Use:   "versions <project-key>",
		Short: "List project versions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}
			versions, err := client.GetProjectVersions(args[0])
			if err != nil {
				return err
			}
			w := newTabWriter()
			_, _ = fmt.Fprintln(w, "ID\tNAME\tSTATUS\tRELEASE DATE")
			for _, v := range versions {
				status := "Unreleased"
				if v.Released {
					status = "Released"
				} else if v.Archived {
					status = "Archived"
				}
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", v.ID, v.Name, status, v.ReleaseDate)
			}
			return w.Flush()
		},
	})

	// --- project statuses ---
	jiraProjectCmd.AddCommand(&cobra.Command{
		Use:   "statuses <project-key>",
		Short: "List project statuses by issue type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}
			issueTypes, err := client.GetProjectStatuses(args[0])
			if err != nil {
				return err
			}
			w := newTabWriter()
			_, _ = fmt.Fprintln(w, "ISSUE TYPE\tSTATUS NAME\tCATEGORY")
			for _, it := range issueTypes {
				for _, s := range it.Statuses {
					_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", it.Name, s.Name, s.StatusCategory.Name)
				}
			}
			return w.Flush()
		},
	})

	// --- project roles ---
	jiraProjectCmd.AddCommand(&cobra.Command{
		Use:   "roles <project-key>",
		Short: "List project roles",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}
			roles, err := client.GetProjectRoles(args[0])
			if err != nil {
				return err
			}
			w := newTabWriter()
			_, _ = fmt.Fprintln(w, "ROLE NAME\tURL")
			for name, url := range roles {
				_, _ = fmt.Fprintf(w, "%s\t%s\n", name, url)
			}
			return w.Flush()
		},
	})

	// --- project archive ---
	jiraProjectCmd.AddCommand(&cobra.Command{
		Use:   "archive <project-key>",
		Short: "Archive a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}
			if err = client.ArchiveProject(args[0]); err != nil {
				return err
			}
			return outputResult(cmd, "archived", args[0], fmt.Sprintf("Archived project %s", args[0]), nil)
		},
	})

	// --- project restore ---
	jiraProjectCmd.AddCommand(&cobra.Command{
		Use:   "restore <project-key>",
		Short: "Restore an archived project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}
			if err = client.RestoreProject(args[0]); err != nil {
				return err
			}
			return outputResult(cmd, "restored", args[0], fmt.Sprintf("Restored project %s", args[0]), nil)
		},
	})

	// --- project features ---
	jiraProjectCmd.AddCommand(&cobra.Command{
		Use:   "features <project-key>",
		Short: "List project features",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getJiraClient(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetProjectFeatures(args[0])
			if err != nil {
				return err
			}
			w := newTabWriter()
			_, _ = fmt.Fprintln(w, "FEATURE\tSTATE")
			for _, f := range resp.Features {
				name := f.LocalisedName
				if name == "" {
					name = f.Feature
				}
				_, _ = fmt.Fprintf(w, "%s\t%s\n", name, f.State)
			}
			return w.Flush()
		},
	})
}
