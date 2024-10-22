package github

// Workflow represents a GitHub Actions workflow with ID, Name, and a slice of run IDs.
type Workflow struct {
	ID   int64
	Name string
	Runs []int64
}

// Repository represents a GitHub repository with owner and repository name.
type Repository struct {
	Owner      string
	Repository string
}
