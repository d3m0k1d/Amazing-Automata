package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	ciMode   bool
	cdMode   bool
	dryRun   bool
	appendM  bool
	filename string
	projroot string
)

var rootCmd = &cobra.Command{
	Use:   "github.com/d3m0k1d/Amazing-Automata [flags] <filename>.yml",
	Short: "A universal CI/CD pipeline generator for GitHub Actions",
	Long: `Amazing-Automata is a universal CI/CD pipeline generator for GitHub Actions
that lets DevOps engineers create and customize workflows in seconds.`,
	Args: cobra.ExactArgs(0),
	RunE: func(c *cobra.Command, args []string) error {
		return YamlGenerator(filename, projroot, ciMode, cdMode, dryRun, appendM)
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&ciMode, "ci", "c", false, "write a new CI pipeline")
	rootCmd.Flags().BoolVarP(&cdMode, "cd", "d", false, "write a new CD pipeline")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "stdout pipeline in shell")
	// rootCmd.Flags().BoolVarP(&appendM, "append", "a", false, "append changes to existing pipeline")
	rootCmd.Flags().StringVarP(&filename, "output-file", "o", "workflow.yml", "file to output generated workflow into")
	rootCmd.Flags().StringVarP(&projroot, "root", "r", ".", "path to the root of the project to walk")

	rootCmd.Flags().SetInterspersed(false)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
