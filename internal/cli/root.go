/*
Package cli provides a command-line interface (CLI) for the delete-gh-workflow-runs tool, utilizing cobra-cli.
*/
package cli

import (
	"context"
	"fmt"
	"log"
	"os"

	gh "github.com/google/go-github/v66/github"
	"github.com/spf13/cobra"

	"github.com/christosgalano/delete-gh-workflow-runs/internal/github"
)

var (
	owner    string
	repo     string
	token    string
	workflow string
)

const maxWorkers = 10

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Version: "v1.0.0",
	Use:     "delete-gh-workflow-runs",
	Short:   "delete-gh-workflow-runs is a command-line tool to delete GitHub Actions workflow runs.",
	Long: `delete-gh-workflow-runs is a command-line tool to delete GitHub Actions workflow runs.

It allows you to delete all workflow runs or only the runs of a specific workflow.
Only 'completed' workflow runs are considered for deletion.

In order to run it, the provided token must have the following permissions:
- Read access to metadata
- Read and Write access to actions`,
	//revive:disable:unused-parameter
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if err := deleteWorkflowRuns(ctx, owner, repo, token, workflow); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

// init initializes the root command.
func init() {
	// Repository owner - required
	rootCmd.Flags().StringVarP(&owner, "owner", "o", "", "repository owner")
	if err := rootCmd.MarkFlagRequired("owner"); err != nil {
		log.Fatalf("Flag 'owner' could not be marked as required.")
	}

	// Repository name - required
	rootCmd.Flags().StringVarP(&repo, "repo", "r", "", "repository name")
	if err := rootCmd.MarkFlagRequired("repo"); err != nil {
		log.Fatalf("Flag 'repo' could not be marked as required.")
	}

	// Token - required
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "api token to get and delete workflow runs")
	if err := rootCmd.MarkFlagRequired("token"); err != nil {
		log.Fatalf("Flag 'token' could not be marked as required.")
	}

	// Workflow name - required
	rootCmd.Flags().StringVarP(&workflow, "workflow", "w", "all", "workflow to delete runs from or 'all' to delete all runs")
}

// deleteWorkflowRuns deletes GitHub workflow runs for a specified repository and workflow.
// It creates a GitHub client using the provided token, retrieves the workflows based on the
// provided workflow name, fetches the run IDs for each workflow, and deletes the workflow runs.
//
// Parameters:
//   - ctx: The context for the request.
//   - owner: The owner of the repository.
//   - repo: The name of the repository.
//   - token: The GitHub authentication token.
//   - workflow: The name of the workflow to delete runs for. If "all", deletes runs for all workflows.
//
// Returns:
//   - error: An error if any step in the process fails, otherwise nil.
func deleteWorkflowRuns(ctx context.Context, owner, repo, token, workflow string) error {
	// Create GitHub client
	client := gh.NewClient(nil).WithAuthToken(token)

	// Create repository
	repository := github.Repository{
		Owner:      owner,
		Repository: repo,
	}

	// Get repository workflows.
	// If the workflow name is "all", it returns all workflows.
	// If a specific workflow name is provided, it returns only that workflow.
	workflows, err := github.GetWorkflowsByName(ctx, client, repository, workflow)
	if err != nil {
		return err
	}

	// Get the run ids for every workflow
	if err := github.GetAllWorkflowRunIDs(ctx, client, repository, workflows); err != nil {
		return err
	}

	fmt.Println("Starting to delete workflow runs...")

	// Delete workflow runs
	if err := github.DeleteWorkflowRuns(ctx, client, repository, workflows, maxWorkers); err != nil {
		return err
	}

	fmt.Println("Workflow runs deleted successfully.")

	return nil
}
