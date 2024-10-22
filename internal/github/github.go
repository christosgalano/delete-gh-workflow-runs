package github

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/go-github/v66/github"
)

// Constants for GitHub API.
const (
	elementPerPage = 100
)

// GetWorkflows retrieves all workflows for a given repository using the provided GitHub client.
// It returns a slice of Workflow objects containing the ID and Name of each workflow.
//
// Parameters:
//   - ctx: The context for the request.
//   - client: The GitHub client used to make API requests.
//   - repo: The repository information including owner and repository name.
//
// Returns:
//   - []Workflow: A slice of Workflow objects containing the ID and Name of each workflow.
//   - error: An error if the request fails, otherwise nil.
func GetWorkflows(ctx context.Context, client *github.Client, repo Repository) ([]Workflow, error) {
	// Get all repository workflows
	repoWorkflows, _, err := client.Actions.ListWorkflows(ctx, repo.Owner, repo.Repository, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflows: %w", err)
	}

	// Extract workflow information
	var workflows []Workflow
	for _, wf := range repoWorkflows.Workflows {
		workflows = append(workflows, Workflow{
			ID:   wf.GetID(),
			Name: wf.GetName(),
		})
	}

	return workflows, nil
}

// GetWorkflowsByName retrieves workflows from a GitHub repository by their name.
// If the specified workflow name is "all", it returns all workflows.
//
// Parameters:
//   - ctx: The context for the request.
//   - client: The GitHub client to use for making API requests.
//   - repo: The repository from which to retrieve workflows.
//   - workflow: The name of the workflow to filter by. If "all", returns all workflows.
//
// Returns:
//   - A slice of Workflow objects.
//   - An error if there was an issue retrieving the workflows.
func GetWorkflowsByName(ctx context.Context, client *github.Client, repo Repository, workflow string) ([]Workflow, error) {
	repoWorkflows, err := GetWorkflows(ctx, client, repo)
	if err != nil {
		return nil, err
	}

	// Filter workflows by name
	var workflows []Workflow
	for _, wf := range repoWorkflows {
		if workflow == "all" || wf.Name == workflow {
			workflows = append(workflows, wf)
		}
	}

	return workflows, nil
}

// GetWorkflowRunIDs retrieves the IDs of all workflow runs for a specified workflow and repository.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation and timeouts.
//   - client: The GitHub client used to interact with the GitHub API.
//   - workflow: The workflow for which to retrieve run IDs.
//   - repo: The repository containing the workflow.
//
// Returns:
//   - A slice of int64 containing the IDs of the workflow runs.
//   - An error if the request to the GitHub API fails or if there are issues processing the response.
func GetWorkflowRunIDs(ctx context.Context, client *github.Client, workflow Workflow, repo Repository) ([]int64, error) {
	var allRunIDs []int64
	opts := &github.ListWorkflowRunsOptions{Status: "completed", ListOptions: github.ListOptions{PerPage: elementPerPage}}
	for {
		// Get a page of workflow runs
		runs, resp, err := client.Actions.ListWorkflowRunsByID(ctx, repo.Owner, repo.Repository, workflow.ID, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get workflow runs for %s: %w", workflow.Name, err)
		}

		// Collect run IDs for completed runs
		for _, run := range runs.WorkflowRuns {
			status := run.GetStatus()
			if status == "completed" {
				allRunIDs = append(allRunIDs, run.GetID())
			}
		}

		// Check if there are more pages
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRunIDs, nil
}

// It performs the retrieval concurrently for each workflow and collects any errors encountered.
//
// Parameters:
//   - ctx: The context to control cancellation and timeouts.
//   - client: The GitHub client used to make API requests.
//   - repo: The repository containing the workflows.
//   - workflows: A slice of Workflow objects for which to retrieve run IDs.
//
// Returns:
//   - error: An error if any of the workflow run ID retrievals fail, otherwise nil.
func GetAllWorkflowRunIDs(ctx context.Context, client *github.Client, repo Repository, workflows []Workflow) error {
	var wg sync.WaitGroup
	errorChan := make(chan error)

	for i := range workflows {
		wg.Add(1)
		go func(wf *Workflow) {
			defer wg.Done()
			runIDs, err := GetWorkflowRunIDs(ctx, client, *wf, repo)
			if err != nil {
				errorChan <- err
				return
			}
			wf.Runs = runIDs
		}(&workflows[i])
	}

	// Close the error channel after all requests are done
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Collect errors
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to get run IDs for some workflows: %v", errors)
	}

	return nil
}

// DeleteWorkflowRun deletes a specified GitHub Actions workflow run.
//
// Parameters:
//   - ctx: The context for the request.
//   - client: The GitHub client to use for making the request.
//   - runID: The ID of the workflow run to delete.
//   - repo: The repository information containing the owner and repository name.
//
// Returns:
//   - error: An error if the deletion fails, otherwise nil.
func DeleteWorkflowRun(ctx context.Context, client *github.Client, runID int64, repo Repository, workflow Workflow) error {
	// Delete workflow run
	_, err := client.Actions.DeleteWorkflowRun(ctx, repo.Owner, repo.Repository, runID)
	if err != nil {
		return fmt.Errorf("error deleting workflow run %d for workflow %s: %w", runID, workflow.Name, err)
	}

	fmt.Printf("Successfully deleted workflow run %d for workflow %q\n", runID, workflow.Name)
	return nil
}

// DeleteWorkflowRuns deletes multiple workflow runs concurrently for a given repository.
// It uses a worker pool to limit the number of concurrent deletions.
//
// Parameters:
//   - ctx: The context for the API requests.
//   - client: A GitHub client to interact with the GitHub API.
//   - repo: The repository containing the workflow runs to be deleted.
//   - workflows: A slice of Workflow objects, each containing the IDs of the runs to be deleted.
//   - maxWorkers: The maximum number of concurrent deletions.
//
// Returns:
//   - error: An error if any of the deletions fail, or nil if all deletions succeed.
func DeleteWorkflowRuns(ctx context.Context, client *github.Client, repo Repository, workflows []Workflow, maxWorkers int) error {
	var wg sync.WaitGroup
	errorChan := make(chan error)

	// Worker pool to limit concurrent deletion
	workerPool := make(chan struct{}, maxWorkers)

	for _, w := range workflows {
		for _, runID := range w.Runs {
			wg.Add(1)
			workerPool <- struct{}{} // block if pool is full
			go func(runID int64, workflow Workflow) {
				defer wg.Done()
				defer func() { <-workerPool }() // release worker slot

				err := DeleteWorkflowRun(ctx, client, runID, repo, workflow)
				if err != nil {
					errorChan <- err
				}
			}(runID, w)
		}
	}

	// Close the error channel after all deletions are done
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Collect errors
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to delete some workflow runs: %v", errors)
	}

	return nil
}
