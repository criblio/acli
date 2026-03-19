package acli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// isJSONOutput returns true if the user requested JSON output via --output json or --json.
func isJSONOutput(cmd *cobra.Command) bool {
	// Check global --output flag first
	if output, _ := cmd.Flags().GetString("output"); output == "json" {
		return true
	}
	// Backward compat: check per-command --json flag if it exists
	if f := cmd.Flags().Lookup("json"); f != nil {
		if v, _ := cmd.Flags().GetBool("json"); v {
			return true
		}
	}
	return false
}

// outputJSON outputs v as indented JSON to stdout. Used for data responses.
func outputJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// OutputResult represents a structured result for mutation operations (create, update, delete, etc).
// Agents can parse this to confirm actions and extract identifiers.
type OutputResult struct {
	Status  string      `json:"status"`
	Action  string      `json:"action"`
	ID      string      `json:"id,omitempty"`
	Key     string      `json:"key,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// outputResult outputs a structured result for mutation operations.
// In text mode, it prints the message. In JSON mode, it outputs a structured envelope.
func outputResult(cmd *cobra.Command, action, key, message string, data interface{}) error {
	if isJSONOutput(cmd) {
		return outputJSON(OutputResult{
			Status:  "ok",
			Action:  action,
			Key:     key,
			Message: message,
			Data:    data,
		})
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), message)
	return nil
}
