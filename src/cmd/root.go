package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:       "ube [path] [flags]",
	Short:     "Ube is a code statistics tool for your terminal.",
	Example:   "  $ ube /path/to/directory \n  $ ube /path/to/file.go",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"path"},
	Version:   "2.1.0",
	Run:       func(cmd *cobra.Command, args []string) {},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
