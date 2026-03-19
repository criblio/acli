package acli

import (
	"fmt"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/spf13/cobra"
)

// addAllFlag adds the --all flag to a command for fetching all pages.
func addAllFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("all", false, "Fetch all pages of results (overrides --max-results/--start-at)")
}

// addBBPaginationFlags adds --page and --pagelen flags for Bitbucket commands.
func addBBPaginationFlags(cmd *cobra.Command) {
	cmd.Flags().Int("page", 0, "Page number (default: server-side default)")
	cmd.Flags().Int("pagelen", 0, "Number of results per page (max 100)")
	cmd.Flags().Bool("all", false, "Fetch all pages of results")
}

// getBBPaginationOpts reads the common Bitbucket pagination flags from a command.
func getBBPaginationOpts(cmd *cobra.Command) *bitbucket.PaginationOptions {
	page, _ := cmd.Flags().GetInt("page")
	pagelen, _ := cmd.Flags().GetInt("pagelen")
	all, _ := cmd.Flags().GetBool("all")
	return &bitbucket.PaginationOptions{
		Page:    page,
		PageLen: pagelen,
		All:     all,
	}
}

// printPaginationHint prints a hint about fetching more results if there are more pages.
func printPaginationHint(cmd *cobra.Command, shown, total int) {
	if total <= 0 || shown >= total {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nShowing %d results\n", shown)
		return
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nShowing %d of %d results (use --all to fetch all)\n", shown, total)
}
