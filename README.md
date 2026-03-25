# GitSubset

**Browse any GitHub repository and download only the files or folders you need, without cloning the whole thing.**

GitSubset is a TUI app. No commands to memorize, no configuration files. Just paste a GitHub URL, explore the repository tree, and grab what you want.

---

## Why GitSubset?

Cloning a large repository to get a single folder is wasteful, it downloads gigabytes of history and files you'll never use. GitSubset lets you browse the repository first, pick exactly what you need, and download only that.

---

## Features

- **Browse before you download**: explore the full file and folder tree of any public GitHub repository
- **Three download modes**: grab the whole repo, select specific folders, or cherry-pick individual files
- **Keyboard-driven interface**: navigate with arrow keys, no mouse required
- **No account or token needed**: works with any public GitHub repository out of the box

---

## Installation

### macOS & Linux

```bash
curl -sL https://raw.githubusercontent.com/kevo-1/GitSubset/master/install.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/kevo-1/GitSubset/master/install.ps1 | iex
```

### Go (any platform)

```bash
go install github.com/kevo-1/GitSubset/gitsubset@latest
```

> **Requirements:** Git must be installed and available in your PATH.

---

## Usage

Run the tool:

```bash
gitsubset
```

Then follow the on-screen steps:

1. **Paste a GitHub URL**: e.g. `https://github.com/user/repo`
2. **Choose what to fetch:**
    - `Whole Repository`: download everything
    - `Select Folders`: browse and pick folders
    - `Select Files`: browse and pick individual files
3. **Confirm your selection**: files are downloaded into a local folder named after the repository

---

## Keyboard Shortcuts

| Key                    | Action                      |
| ---------------------- | --------------------------- |
| `↑` / `↓` or `k` / `j` | Navigate                    |
| `Space`                | Select / deselect item      |
| `Tab`                  | Expand or collapse a folder |
| `a`                    | Select or deselect all      |
| `Enter`                | Confirm                     |
| `Esc`                  | Go back                     |
| `r`                    | Retry after an error        |
| `q` / `Ctrl+C`         | Quit                        |

---

## How It Works

GitSubset uses Git's [sparse checkout](https://git-scm.com/docs/git-sparse-checkout) feature under the hood. When you enter a URL it performs a metadata-only clone (no file contents), fetches the file tree so you can browse it, then downloads only the paths you selected using `git sparse-checkout`.

---

## Supported Platforms

| Platform | Architecture          |
| -------- | --------------------- |
| macOS    | Intel (amd64)         |
| macOS    | Apple Silicon (arm64) |
| Linux    | amd64                 |
| Windows  | amd64                 |

---

## Contributing

Issues and pull requests are welcome. Please open an issue first if you're planning a larger change.

---

## License

MIT
