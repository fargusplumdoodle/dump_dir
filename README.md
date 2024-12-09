# 🚀 dump_dir 📂✨

Copy a bunch of files into your clipboard to provide context for LLMs

## 🌟 Functionality

- 🔍 Search for files with specific extensions or all files
- 📋 Automatically copy file contents to clipboard
- 🚫 Skip specified directories
- 📝 Respects .gitignore rules by default
- 🚀 Fast! Can copy 1.6 million tokens in ~0.8 seconds

_Example: dump_dir-ing this repo_

![image](https://github.com/user-attachments/assets/17e52273-871b-44a0-8f0b-10f91eb4ad25)



## 🚀 Usage


```bash
dump_dir <directory1> [directory2] ...
```

#### 📚 Options

- `-h`, `--help`: Display help information
- `-v`, `--version`: Display the version of `dump_dir`
- `-s <directory>, --skip <directory>`: Skip specified directory
- `-e <extension[s]>, --extension <extension[s]>`: Filter by specific file extensions
- `--include-ignored`: Include files that would normally be ignored (e.g., those in `.gitignore`)
- `-m <size>`, `--max-filesize <size>`: Specify the maximum file size to process. You can use units like B, KB, or MB (e.g., 500KB, 2MB). If no unit is specified, it defaults to bytes.

#### 📑 Examples

Get all files in your project directory of all types:
```bash
dump_dir ./project
````
Get all JS files, ignoring your dist directory:
```bash
dump_dir ./project -e js  --skip ./project/dist
```
Get all Go, JavaScript, and Python files in your project:
```bash
dump_dir ./project --extension go,js,py 
```
Get all files, including those normally ignored (e.g., files in .gitignore):
```bash
dump_dir ./project --include-ignored
```
Get specific files regardless of their extension:
```bash
dump_dir ./README.md ./main.go --extension py
```
Set a maximum file size of 1MB:
```bash
dump_dir ./project --max-filesize 1MB
```

## 🔒 Gitignore Behavior
By default, dump_dir respects your project's .gitignore rules. This means:

Files and directories listed in your project's .gitignore will not be included in the output.
The tool also respects your global gitignore file.
Common version control directories (like .git) are automatically ignored.

To include ignored files, use the `--include-ignored` flag as shown in the examples above.

## 👉 Special Files Behavior

By default, files are too large if they are >500KB. You can adjust this limit using the `-m` or `--max-filesize` option.
If any files are skipped for any reason, `dump_dir` will inform you.

| File type       | Output                          |
|-----------------|---------------------------------|
| Binary files    | `<BINARY SKIPPED>`              |
| File too large  | `<FILE TOO LARGE: %d bytes>`    |
| Empty files     | `<EMPTY FILE>`                  |

## 📝 Configuration File

`dump_dir` will check for a file in your current directory
called `.dump_dir.yml` which can be used to specify default options. 

Example `.dump_dir.yml` file:
```yaml
---
# Always include these paths
include:
  - ./README.md
  - ./prompts

# Always skip these directories
# (even if they aren't in your gitignore)
ignore:
  - ./dist 
  - ./vendor
```

You can check the config file of this repo as another example.

**Purpose**

Including things like coding standards, general architecture patterns,
and high level explanations of your work can produce dramatically
better results with LLMs.


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
- 📚 Include documentation in your prompts
- 🧠 Paste the copied content directly into your LLM conversation
- 🔄 Easily update context by re-running `dump_dir` with different parameters

## 🤝 Contributing

Raise an issue or make a PR. Run tests with `go test ./tests/...`

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


---------------------------------------------

Happy coding! 🎉👨‍💻👩‍💻

##################
##################
##################


Gitignore Integration Tests:

Test respecting basic .gitignore patterns
Test with --include-ignored flag
Test with nested .gitignore files
Test with global gitignore rules
Test with complex gitignore patterns


Special File Handling Tests:

Test handling of binary files
Test handling of symbolic links
Test handling of files with special characters in names
Test handling of hidden files
Test handling of files with no extension


Configuration File Tests:

Test with valid .dump_dir.yml configuration
Test with invalid YAML syntax in config
Test config file directory inclusion rules
Test config file directory exclusion rules
Test precedence of CLI arguments over config file


Output Format Tests:

Test file content markers (START FILE/END FILE)
Test correct line counting
Test token estimation accuracy
Test handling of different line endings (CRLF/LF)
Test output formatting of skipped files


CLIPBOARD OUTPUT TESTS:
- Binary files
- Too large files
- Empty files
- Normal files

Path Handling Tests:

Test with absolute paths
Test with relative paths
Test with multiple input directories
Test with duplicate paths
Test with nonexistent paths


Error Handling Tests:

Test with insufficient permissions
Test with inaccessible directories
Test with corrupted files
Test with invalid command flags
Test with conflicting arguments


Performance Tests:

Test with large number of small files
Test with few large files
Test with deep directory structures
Test with mixed file sizes and types
Test token counting performance

