# Lucide for Go

<p>
  <a href="https://www.figma.com/community/plugin/939567362549682242/Lucide-Icons"><img src="https://img.shields.io/badge/Figma-F24E1E?logo=figma&logoColor=white" alt="figma installs"></a>
</p>
<p>
  <a href="https://lucide.dev/icons/">Icons</a>
  ·
  <a href="https://lucide.dev/guide/">Guide</a>
</p>

`go-lucide` provides golang ports for [Lucide](https://github.com/lucide-icons/lucide) icons.

Lucide is an open-source icon library that provides 1000+ vector (svg) files for displaying icons and symbols in digital and non-digital projects. The library aims to make it easier for designers and developers to incorporate icons into their projects by providing several official [packages](https://lucide.dev/packages) to make it easier to use these icons in your project.

## Packages

| Package | Version | Links |
| ------- | ------- | ----- |
| **`go-templ-lucide-icons`** | [![npm](https://img.shields.io/github/v/release/bryanvaz/go-templ-lucide-icons)](https://github.com/bryanvaz/go-templ-lucide-icons/releases) | [Docs](https://pkg.go.dev/github.com/bryanvaz/go-templ-lucide-icons) · [Source](https://github.com/bryanvaz/go-templ-lucide-icons) |

### Figma

The lucide figma plugin.

Visit [Figma community page](https://www.figma.com/community/plugin/939567362549682242/Lucide-Icons) to install the plugin.

<img width="420" src="https://www.figma.com/community/plugin/939567362549682242/thumbnail" alt="Figma Lucide Cover">

## Usage

### Prerequisites

* Go 1.23
* `gh` CLI

### Sync next version

```bash
make clean
make deps
make build
make test
make commit
make publish
```

### Sync specific version

```bash
make build TARGET=v0.465.0
```

## License

Lucide is totally free for commercial use and personal use, this software is licensed under the [ISC License](https://github.com/lucide-icons/lucide/blob/main/LICENSE).

This wrapper is licensed under the MIT license.

## Sponsors

This library is currently supported on a best-effort basis. 
If you would like to sponsor this project, please reach out via Gihub.
