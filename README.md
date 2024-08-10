# 🚀 dump_dir 📂✨

Copy a bunch of files into your clipboard to provide context for LLMs

## 🌟 Functionality

- 🔍 Search for files with specific extensions or all files
- 📋 Automatically copy file contents to clipboard
- 🚫 Skip specified directories
- 📝 Respects .gitignore rules by default
- 🚀 Fast! Can copy 1.6 million tokens in ~0.8 seconds

_Example: dump_dir-ing this repo_

![image](https://github.com/user-attachments/assets/ae8bc680-8da0-4f50-9092-6b6f89a2a9ad)

## 🚀 Usage


```bash
dump_dir [options] <file_extension1> [,<file_extension2>,...] <directory1> [directory2] ... [-s <skip_directory1>] [-s <skip_directory2>] ... [--include-ignored]
```

Use 'any' as file_extension to match all files.


#### 📚 Options

- `-h`, `--help`: Display help information
- `-v`, `--version`: Display the version of `dump_dir`
- `-s <directory>`: Skip specified directory
- `--include-ignored`: Include files that would normally be ignored (e.g., those in `.gitignore`)

#### 📑 Examples

Get all JS files, ignoring your node_modules and dist directories:
```bash
dump_dir js ./project -s ./project/node_modules -s ./project/dist
```
Get all files in your project directory of all types:
```bash
dump_dir any ./project
````
Get all Go, JavaScript, and Python files in your project:
```bash
dump_dir go,js,py ./project
```
Get all files, including those normally ignored (e.g., files in .gitignore):
```bash
dump_dir any ./project --include-ignored
```
Get specific files regardless of their extension:
```bash
dump_dir any ./README.md ./main.go
```

## 🔒 Gitignore Behavior
By default, dump_dir respects your project's .gitignore rules. This means:

Files and directories listed in your project's .gitignore will not be included in the output.
The tool also respects your global gitignore file.
Common version control directories (like .git) are automatically ignored.

To include ignored files, use the `--include-ignored` flag as shown in the examples above.

## 👉 Special Files Behavior

By default, files are too large if they are >500KB. 

| File type       | Output                          |
|-----------------|---------------------------------|
| Binary files    | `<BINARY SKIPPED>`              |
| File too large  | `<FILE TOO LARGE: %d bytes>`    |
| Empty files     | `<EMPTY FILE>`                  |


## 🛠️ Installation

**Mac and Linux**

You can easily install the latest version of dump_dir using curl. Run the following command in your terminal:
```bash
curl -sfL https://raw.githubusercontent.com/fargusplumdoodle/dump_dir/main/install.sh | bash
```

This script will:
- Detect your operating system and architecture
- Download the latest release of dump_dir
- Install it to /usr/local/bin (you may need to use sudo for this)

Alternatively, you can manually download the latest release from the GitHub Releases page and place it in a directory in your PATH.

**Windows**

For Windows users, please download the latest release from the GitHub Releases page and add it to your PATH manually.

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
