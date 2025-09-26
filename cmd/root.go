package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ciMode  bool
	cdMode  bool
	dryRun  bool
	appendM bool
)

var rootCmd = &cobra.Command{
	Use:   "amazing-automata [flags] <filename>.yml",
	Short: "A universal CI/CD pipeline generator for GitHub Actions",
	Long: `Amazing-Automata is a universal CI/CD pipeline generator for GitHub Actions
that lets DevOps engineers create and customize workflows in seconds.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]
		fmt.Println("File:", filename)
		if ciMode {

		}
		if cdMode {
			fmt.Println("– CD mode enabled")
		}
		if dryRun {
			fmt.Println("– Dry run mode enabled")
		}
		if appendM {
			fmt.Println("– Append mode enabled")
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&ciMode, "ci", "c", false, "write a new CI pipeline")
	rootCmd.Flags().BoolVarP(&cdMode, "cd", "d", false, "write a new CD pipeline")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "stdout pipeline in shell")
	rootCmd.Flags().BoolVarP(&appendM, "append", "a", false, "append changes to existing pipeline")

	rootCmd.Flags().SetInterspersed(false)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
