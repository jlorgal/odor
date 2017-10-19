# Odor

A more efficient approach of niji platform.

## Usage

Clone this repository inside directory `$GOPATH/src/github.com/jlorgal/odor`.

Then you can use make to build the project.

```
Usage: make <command>
Commands:
  help:            Show this help information
  dep:             Ensure dependencies with dep tool
  build:           Build the application
  test-acceptance: Pass component tests
  release:         Create a new release (tag and release notes)
  run:             Launch the service with docker-compose (for testing purposes)
  clean:           Clean the project
  pipeline-pull:   Launch pipeline to handle a pull request
  pipeline-dev:    Launch pipeline to handle the merge of a pull request
  pipeline:        Launch the pipeline for the selected environment
```
