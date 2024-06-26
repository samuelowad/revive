# Revive

Revive is a tool for enabling hot reload in applications just with a command definition. It watches for file changes in the current directory and its subdirectories, and automatically restarts the application when any Go file is modified.
like nodemon written in JS.

## Installation

Downloadable binary is not available yet, but you can clone the repository and build locally using the following command:

```sh
## locally
go build -o revive
```
### or

```sh
## globally
go install

```

## Configuration
Create a configuration file named `revive.yaml` or `revive.json` in the root of your project with the following structure:

```yaml
command: go run main.go
ignoreDirectories:
  - vendor
  - tmp
monitorFileExt:
  - .go
ignoreFileNameEndsWith:
  - _test.go

```
```json
{
  "command": "node test-data/test.js",
  "ignoreDirectories": ["vendor", "tmp"],
  "monitorFileExt": [".go", ".html", ".css", ".js", ".json"],
  "ignoreFileNameEndsWith": ["_test.go"]
}

```

## Usage
Execute the binary to start the application and enable hot reload:

```sh
## if complied locally
./revive
```

```sh
## if installed globally
revive
```

