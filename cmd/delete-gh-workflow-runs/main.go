/*
MIT License

Copyright (c) 2024 Christos Galanopoulos

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/*
delete-gh-workflow-runs is a command-line tool to delete GitHub Actions workflow runs.

It allows you to delete all workflow runs or only the runs of a specific workflow.

In order to run it, the provided token must have the following permissions:
- Read access to metadata
- Read and Write access to actions

Example usage:

Delete the runs of a specific workflow:

	delete-gh-workflow-runs --owner owner --repo repo --token token --workflow workflow

Delete all workflow runs of a repository:

	delete-gh-workflow-runs --owner owner --repo repo --token token

For full usage details, run `delete-gh-workflow-runs --help`.
*/

package main

import (
	"os"

	"github.com/christosgalano/delete-gh-workflow-runs/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
