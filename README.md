# 🚀 dump_dir 📂✨

Copy a bunch of files into your clipboard to provide context for LLMs

## 🌟 Functionality

- 🔍 Search for files with specific extensions or all files
- 📋 Automatically copy file contents to clipboard
- 🚫 Skip specified directories
- 📝 Respects .gitignore rules by default


## 🚀 Usage

```
dump_dir <file_extension> <directory1> [directory2] ... [-s <skip_directory1>] [-s <skip_directory2>] ...
```

Use 'any' as file_extension to match all files.

### 📚 Examples


Get all JS files, ignoring your node_modules and dist directories:
```bash
dump_dir js ./project -s ./project/node_modules -s ./project/dist
```

Get all files in your project directory of all types
```bash
dump_dir any ./project 
```

Get all Go, JavaScript, and Python files in your project:
```bash
dump_dir go,js,py ./project
````

Get all files, including those normally ignored (e.g., files in .gitignore):
```bash
dump_dir any ./project --include-ignored
```

This repo:
```txt
dump_dir any .
Skipping ignored directory: .git
Skipping ignored directory: .idea
Skipping ignored directory: gitignore

🔍 Matching files:
  - dump_dir
  - src/colors.go
  - README.md
  - .gitignore
  - main.go
  - go.mod
  - src/ignore.go
  - src/file_processing.go
  - src/output.go
  - go.sum
  - .goreleaser.yaml
  - src/types.go
  - src/args.go
  - .github/workflows/release.yml

📚 Total files found: 14
📝 Total lines across all files: 710

✅ File contents have been copied to clipboard.
```

## 🔒 Gitignore Behavior
By default, dump_dir respects your project's .gitignore rules. This means:

Files and directories listed in your project's .gitignore will not be included in the output.
The tool also respects your global gitignore file.
Common version control directories (like .git) are automatically ignored.

To include ignored files, use the `--include-ignored` flag as shown in the examples above.

## 🛠️ Installation

See releases in github!


## 💡 Tips for LLM Development

- 📁 Use `dump_dir` to quickly gather context from multiple project files
- 🧠 Paste the copied content directly into your LLM conversation
- 🔄 Easily update context by re-running `dump_dir` with different parameters

## 🤝 Contributing

Raise an issue or make a PR, I whipped this up real quick so there's probably a lot of room for improvement.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


---------------------------------------------

Happy coding! 🎉👨‍💻👩‍💻