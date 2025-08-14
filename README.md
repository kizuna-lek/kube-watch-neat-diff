# kube-watch-neat-diff

A tool to watch Kubernetes resources and display clean, human-readable diffs when changes occur.

![GitHub](https://img.shields.io/github/license/yourusername/kube-watch-neat-diff)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/yourusername/kube-watch-neat-diff)

## Overview

`kube-watch-neat-diff` monitors Kubernetes resources in real-time and presents changes in a clean, colorized diff format. It uses `kubectl` to watch resources, `kubectl-neat` to clean up the JSON output, and a diff library to show changes between resource versions.

This tool is particularly useful for:
- Debugging configuration changes
- Monitoring deployments and other resources
- Understanding exactly what changes occur over time

## Features

- Real-time monitoring of Kubernetes resources
- Clean, human-readable diff output
- Colorized output for better visibility
- Option to diff against the first version or previous version
- Automatic cleanup of Kubernetes JSON output using kubectl-neat

## Prerequisites

- Go 1.16 or higher
- Kubernetes cluster access configured with kubectl
- Proper kubeconfig setup

## Installation

### From source

```bash
# Clone the repository
git clone https://github.com/yourusername/kube-watch-neat-diff.git
cd kube-watch-neat-diff

# Build the binary
go build -o kube-watch-neat-diff .

# Optionally move to a directory in your PATH
sudo mv kube-watch-neat-diff /usr/local/bin/
```

## Usage

```bash
kube-watch-neat-diff [resource-type] [resource-name] [flags]
```

### Examples

```bash
# Watch a deployment named "my-app"
kube-watch-neat-diff deployment my-app

# Watch a service named "nginx-service"
kube-watch-neat-diff service nginx-service

# Watch with diff against first version instead of previous
kube-watch-neat-diff deployment my-app --diff-with-first

# Watch without color output
kube-watch-neat-diff deployment my-app --no-color
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--diff-with-first` | `-f` | Diff with first version instead of previous version |
| `--no-color` | | Disable colored output |
| `--help` | `-h` | Show help message |
| `--version` | | Show version information |

## How it works

1. Takes Kubernetes resource type and name as arguments
2. Runs `kubectl get -w` to watch the resource
3. For each update:
   - Cleans the JSON with kubectl-neat
   - Compares with the previous version using diff
   - Formats and displays the differences
4. Supports options like diff-with-first and no-color

## Output Format

The tool displays changes in a structured format:
- Created fields are shown in green with a "+" prefix
- Updated fields are shown in yellow with both old and new values
- Deleted fields are shown in red with a "-" prefix

## Dependencies

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) - Kubernetes command-line tool
- [kubectl-neat](https://github.com/itaysk/kubectl-neat) - Cleans up Kubernetes YAML/JSON output
- [kingpin](https://github.com/alecthomas/kingpin) - Command line argument parsing
- [diff](https://github.com/r3labs/diff) - JSON diffing library

## Building

```bash
# Build the project
go build -o kube-watch-neat-diff .

# Install dependencies
go mod tidy
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Acknowledgments

- Thanks to the Kubernetes community for the great tooling
- Inspired by the need for better observability of Kubernetes resource changes