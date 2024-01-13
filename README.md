# Labeler - Declarative GitHub Label Syncing CLI Tool

[![GitHub Release](https://img.shields.io/github/v/release/shanduur/labeler)](https://github.com/shanduur/labeler/releases/latest)
[![GitHub License](https://img.shields.io/github/license/shanduur/labeler)](LICENSE)

Labeler is a command-line interface (CLI) tool designed to facilitate the declarative synchronization of GitHub labels across repositories. This tool provides an efficient and convenient way to manage labels in a GitHub repository by enabling users to define label configurations in a YAML file and then apply or retrieve these configurations using the provided commands.

## Table of Contents

- [Labeler - Declarative GitHub Label Syncing CLI Tool](#labeler---declarative-github-label-syncing-cli-tool)
  - [Table of Contents](#table-of-contents)
  - [Usage](#usage)
    - [Uploading Labels](#uploading-labels)
    - [Downloading Labels](#downloading-labels)
  - [Configuration File](#configuration-file)
  - [Examples](#examples)
  - [License](#license)
  - [Changelog](#changelog)

## Usage

### Uploading Labels

To sync labels from a file to a GitHub repository, use the following command:

```bash
labeler upload [options] <path-to-labels-file>
```

Example:

```bash
labeler upload --repository <owner>/<repo> examples/cli/labels.yaml
```

### Downloading Labels

To download and print GitHub repository labels to a file, use the following command:

```bash
labeler download [options] <path-to-output-file>
```

Example:

```bash
labeler download --repository <owner>/<repo> examples/cli/labels.yaml
```

**Note:** Replace `<owner>/<repo>` with the GitHub owner and repository name you intend to work with.

## Configuration File

The configuration file for Labeler is a YAML file that defines the labels and their attributes to be synchronized with a GitHub repository. Below is an example of a label configuration file:

```yaml
- name: dependency
  color: '#dadada'
  description: Dependency updates
- name: feature-request
  color: '#61dafb'
  description: Request for a new feature or enhancement
- name: documentation
  color: '#007bc7'
  description: Issues related to documentation
- name: bug
  color: '#d73a4a'
  description: Something isn't working as expected
```

In this example, each label is defined with the following attributes:
- **name:** The name of the label.
- **color:** The color associated with the label, specified in hexadecimal format.
- **description:** A brief description of the label's purpose or meaning.

Ensure that the configuration file adheres to the YAML syntax, and you can customize it according to your project's specific label requirements.

## Examples

Example label configurations are provided in the `examples` directory for different scenarios, including CLI, Kubernetes, and Moby projects.

## License

Labeler is licensed under the [MIT License](LICENSE).

## Changelog

Check out the [Changelog](CHANGELOG.md) to see the release history and updates.

For detailed information on commands and options, run `labeler help`.
