# delete-gh-workflow-runs

[![Go Report Card](https://goreportcard.com/badge/github.com/christosgalano/delete-gh-workflow-runs)](https://goreportcard.com/report/github.com/christosgalano/delete-gh-workflow-runs)
[![Go Reference](https://pkg.go.dev/badge/github.com/christosgalano/delete-gh-workflow-runs.svg)](https://pkg.go.dev/github.com/christosgalano/delete-gh-workflow-runs)
[![Github Downloads](https://img.shields.io/github/downloads/christosgalano/delete-gh-workflow-runs/total.svg)](https://github.com/christosgalano/delete-gh-workflow-runs/releases)

## Table of contents

- [Description](#description)
- [Installation](#installation)
- [Requirements](#requirements)
- [Usage](#usage)
- [GitHub Action](#github-action)
- [Contributing](#contributing)
- [License](#license)

## Description

**delete-gh-workflow-runs** is a command-line tool that deletes GitHub Actions workflow runs based on the provided input.

## Installation

### Homebrew

```bash
brew tap christosgalano/christosgalano
brew install delete-gh-workflow-runs
```

### Go

```bash
go install github.com/christosgalano/delete-gh-workflow-runs/cmd/delete-gh-workflow-runs@latest
```

### Binary

Download the latest binary from the [releases page](https://github.com/christosgalano/delete-gh-workflow-runs/releases/latest).

## Requirements

To run delete-gh-workflow-runs, you must have a GitHub token with the `repo` scope and `workflow` permissions.

![Permissions](assets/images/permissions.png)

## Usage

delete-gh-workflow-runs is a command-line tool that deletes GitHub Actions workflow runs based on the provided input.

**Arguments:**

- `--owner` - The owner of the repository.
- `--repo` - The name of the repository.
- `--workflow` - The name of the workflow or "all" to delete all workflow runs; default is "all".
- `--token` - The GitHub token.

> NOTE: Only 'completed' workflow runs are considered for deletion.

### Example usage

Delete the runs of a specific workflow:

```bash
delete-gh-workflow-runs --owner {owner} --repo {repo} --workflow {workflow} --token {token}
```

Delete all workflow runs of a repository:

```bash
delete-gh-workflow-runs --owner {owner} --repo {repo} --token {token}
```

## GitHub Action

You can use this tool as a GitHub Action. Here is an example:

### Syntax

```yaml
uses: christosgalano/delete-gh-workflow-runs@v1.0.0
with:
  owner: ${{ github.repository_owner }}   # The owner of the repository
  repo: ...                               # The name of the repository
  workflow: "all"                         # The name of the workflow or "all" to delete all workflow runs
  token: ${{ github.TOKEN }}              # The GitHub token
```

### Examples

Delete the runs of a specific workflow:

```yaml
- name: Delete workflow runs
  uses: christosgalano/delete-gh-workflow-runs@v1.0.0
  with:
    owner: ${{ github.repository_owner }}
    repo: ${{ inputs.repository }}
    workflow: ${{ inputs.workflow }}
    token: ${{ github.TOKEN }}
```

Delete all workflow runs of a repository:

```yaml
- name: Delete all workflow runs
  uses: christosgalano/delete-gh-workflow-runs@v1.0.0
  with:
    owner: ${{ github.repository_owner }}
    repo: ${{ inputs.repository }}
    token: ${{ github.TOKEN }}
```

## Contributing

Information about contributing to this project can be found [here](CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](LICENSE).
