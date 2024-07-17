# 🚀 dump_dir 📂✨

Copy a bunch of files into your clipboard to provide context for LLms

## 🌟 Functionality

- 🔍 Search for files with specific extensions or all files
- 📋 Automatically copy file contents to clipboard
- 🚫 Skip specified directories

## 🛠️ Installation

[Add installation instructions here]

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

This repo:
```txt
dump_dir go . ./README.md
🔍 Matching files:
  - README.md
  - main.go
  - src/args.go
  - src/colors.go
  - src/file_processing.go
  - src/output.go
  - src/types.go
📚 Total files found: 6
📝 Total lines across all files: 260

✅ File contents have been copied to clipboard.
```


## 💡 Tips for LLM Development

- 📁 Use `dump_dir` to quickly gather context from multiple project files
- 🧠 Paste the copied content directly into your LLM conversation
- 🔄 Easily update context by re-running `dump_dir` with different parameters

## 🤝 Contributing

Raise an issue or make a PR, I made this in 2 hours so don't judge

## 📜 License

Do whatever, I don't mind

## 🙏 Acknowledgements

- [fatih/color](https://github.com/fatih/color) for colorful console output
- [atotto/clipboard](https://github.com/atotto/clipboard) for clipboard functionality

Happy coding and LLM developing! 🎉👨‍💻👩‍💻