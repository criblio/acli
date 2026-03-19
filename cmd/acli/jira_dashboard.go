package acli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ============================================================================
// Dashboards
// ============================================================================

var jiraDashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Short:   "Manage dashboards",
	Aliases: []string{"dash"},
	RunE:    helpRunE,
}

var jiraDashboardListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List dashboards",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		startAt, _ := cmd.Flags().GetInt("start-at")
		maxResults, _ := cmd.Flags().GetInt("max-results")
		name, _ := cmd.Flags().GetString("name")
		all, _ := cmd.Flags().GetBool("all")
		jsonFlag := isJSONOutput(cmd)

		if name != "" {
			result, err := client.SearchDashboards(name, startAt, maxResults)
			if err != nil {
				return err
			}
			if all {
				for !result.IsLast && len(result.Values) < result.Total {
					next, err := client.SearchDashboards(name, startAt+len(result.Values), maxResults)
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
			_, _ = fmt.Fprintf(w, "ID\tNAME\tOWNER\n")
			for _, d := range result.Values {
				owner := ""
				if d.Owner != nil {
					owner = d.Owner.DisplayName
				}
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", d.ID, d.Name, owner)
			}
			_ = w.Flush()
			printPaginationHint(cmd, len(result.Values), result.Total)
			return nil
		}

		result, err := client.GetDashboards(startAt, maxResults)
		if err != nil {
			return err
		}
		if all {
			for len(result.Dashboards) < result.Total {
				next, err := client.GetDashboards(startAt+len(result.Dashboards), maxResults)
				if err != nil {
					return err
				}
				if len(next.Dashboards) == 0 {
					break
				}
				result.Dashboards = append(result.Dashboards, next.Dashboards...)
			}
		}
		if jsonFlag {
			return outputJSON(result)
		}
		w := newTabWriter()
		_, _ = fmt.Fprintf(w, "ID\tNAME\tOWNER\n")
		for _, d := range result.Dashboards {
			owner := ""
			if d.Owner != nil {
				owner = d.Owner.DisplayName
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", d.ID, d.Name, owner)
		}
		_ = w.Flush()
		printPaginationHint(cmd, len(result.Dashboards), result.Total)
		return nil
	},
}

var jiraDashboardGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a dashboard by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		d, err := client.GetDashboard(args[0])
		if err != nil {
			return err
		}
		jsonFlag := isJSONOutput(cmd)
		if jsonFlag {
			return outputJSON(d)
		}
		w := newTabWriter()
		_, _ = fmt.Fprintf(w, "ID\t%s\n", d.ID)
		_, _ = fmt.Fprintf(w, "Name\t%s\n", d.Name)
		_, _ = fmt.Fprintf(w, "Description\t%s\n", d.Description)
		if d.Owner != nil {
			_, _ = fmt.Fprintf(w, "Owner\t%s\n", d.Owner.DisplayName)
		}
		_ = w.Flush()
		return nil
	},
}

var jiraDashboardCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		body := map[string]interface{}{
			"name": name,
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		d, err := client.CreateDashboard(body)
		if err != nil {
			return err
		}
		fmt.Printf("Dashboard created: %s (ID: %s)\n", d.Name, d.ID)
		return nil
	},
}

var jiraDashboardUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a dashboard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		body := map[string]interface{}{}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			body["name"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		d, err := client.UpdateDashboard(args[0], body)
		if err != nil {
			return err
		}
		fmt.Printf("Dashboard updated: %s (ID: %s)\n", d.Name, d.ID)
		return nil
	},
}

var jiraDashboardDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a dashboard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.DeleteDashboard(args[0]); err != nil {
			return err
		}
		fmt.Printf("Dashboard %s deleted\n", args[0])
		return nil
	},
}

var jiraDashboardCopyCmd = &cobra.Command{
	Use:   "copy <id>",
	Short: "Copy a dashboard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		body := map[string]interface{}{
			"name": name,
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		d, err := client.CopyDashboard(args[0], body)
		if err != nil {
			return err
		}
		fmt.Printf("Dashboard copied: %s (ID: %s)\n", d.Name, d.ID)
		return nil
	},
}

// ============================================================================
// Dashboard Gadgets
// ============================================================================

var jiraDashboardGadgetCmd = &cobra.Command{
	Use:   "gadget",
	Short: "Manage dashboard gadgets",
	RunE:  helpRunE,
}

var jiraDashboardGadgetListCmd = &cobra.Command{
	Use:     "list <dashboard-id>",
	Aliases: []string{"ls"},
	Short:   "List gadgets on a dashboard",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		gadgets, err := client.GetDashboardGadgets(args[0])
		if err != nil {
			return err
		}
		jsonFlag := isJSONOutput(cmd)
		if jsonFlag {
			return outputJSON(gadgets)
		}
		w := newTabWriter()
		_, _ = fmt.Fprintln(w, "ID\tTITLE\tMODULE KEY\tURI")
		for _, g := range gadgets.Gadgets {
			_, _ = fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", g.ID, g.Title, g.ModuleKey, g.URI)
		}
		return w.Flush()
	},
}

var jiraDashboardGadgetAddCmd = &cobra.Command{
	Use:   "add <dashboard-id>",
	Short: "Add a gadget to a dashboard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		body := map[string]interface{}{}
		if v, _ := cmd.Flags().GetString("module-key"); v != "" {
			body["moduleKey"] = v
		}
		if v, _ := cmd.Flags().GetString("uri"); v != "" {
			body["uri"] = v
		}
		if v, _ := cmd.Flags().GetString("title"); v != "" {
			body["title"] = v
		}
		gadget, err := client.AddDashboardGadget(args[0], body)
		if err != nil {
			return err
		}
		fmt.Printf("Gadget added: %s (ID: %d)\n", gadget.Title, gadget.ID)
		return nil
	},
}

var jiraDashboardGadgetUpdateCmd = &cobra.Command{
	Use:   "update <dashboard-id> <gadget-id>",
	Short: "Update a gadget on a dashboard",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		body := map[string]interface{}{}
		if v, _ := cmd.Flags().GetString("title"); v != "" {
			body["title"] = v
		}
		if err := client.UpdateDashboardGadget(args[0], args[1], body); err != nil {
			return err
		}
		fmt.Printf("Gadget %s updated on dashboard %s\n", args[1], args[0])
		return nil
	},
}

var jiraDashboardGadgetRemoveCmd = &cobra.Command{
	Use:   "remove <dashboard-id> <gadget-id>",
	Short: "Remove a gadget from a dashboard",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getJiraClient(cmd)
		if err != nil {
			return err
		}
		if err := client.RemoveDashboardGadget(args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("Gadget %s removed from dashboard %s\n", args[1], args[0])
		return nil
	},
}

func init() {
	// Dashboard list
	jiraDashboardListCmd.Flags().String("name", "", "Search dashboards by name")
	jiraDashboardListCmd.Flags().Int("max-results", 50, "Maximum number of results per page")
	jiraDashboardListCmd.Flags().Int("start-at", 0, "Index of the first result")
	addAllFlag(jiraDashboardListCmd)
	jiraDashboardListCmd.Flags().Bool("json", false, "Output as JSON")
	jiraDashboardCmd.AddCommand(jiraDashboardListCmd)

	// Dashboard get
	jiraDashboardGetCmd.Flags().Bool("json", false, "Output as JSON")
	jiraDashboardCmd.AddCommand(jiraDashboardGetCmd)

	// Dashboard create
	jiraDashboardCreateCmd.Flags().String("name", "", "Dashboard name (required)")
	jiraDashboardCreateCmd.Flags().String("description", "", "Dashboard description")
	_ = jiraDashboardCreateCmd.MarkFlagRequired("name")
	jiraDashboardCmd.AddCommand(jiraDashboardCreateCmd)

	// Dashboard update
	jiraDashboardUpdateCmd.Flags().String("name", "", "Dashboard name")
	jiraDashboardUpdateCmd.Flags().String("description", "", "Dashboard description")
	jiraDashboardCmd.AddCommand(jiraDashboardUpdateCmd)

	// Dashboard delete
	jiraDashboardCmd.AddCommand(jiraDashboardDeleteCmd)

	// Dashboard copy
	jiraDashboardCopyCmd.Flags().String("name", "", "Name for the copy (required)")
	jiraDashboardCopyCmd.Flags().String("description", "", "Description for the copy")
	_ = jiraDashboardCopyCmd.MarkFlagRequired("name")
	jiraDashboardCmd.AddCommand(jiraDashboardCopyCmd)

	// Gadget subcommand
	jiraDashboardGadgetListCmd.Flags().Bool("json", false, "Output as JSON")
	jiraDashboardGadgetCmd.AddCommand(jiraDashboardGadgetListCmd)

	jiraDashboardGadgetAddCmd.Flags().String("module-key", "", "Gadget module key")
	jiraDashboardGadgetAddCmd.Flags().String("uri", "", "Gadget URI")
	jiraDashboardGadgetAddCmd.Flags().String("title", "", "Gadget title")
	jiraDashboardGadgetCmd.AddCommand(jiraDashboardGadgetAddCmd)

	jiraDashboardGadgetUpdateCmd.Flags().String("title", "", "Gadget title")
	jiraDashboardGadgetCmd.AddCommand(jiraDashboardGadgetUpdateCmd)

	jiraDashboardGadgetCmd.AddCommand(jiraDashboardGadgetRemoveCmd)

	jiraDashboardCmd.AddCommand(jiraDashboardGadgetCmd)
}
