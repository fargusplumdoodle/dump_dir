# 🚀 dump_dir 📂✨

Copy a bunch of files into your clipboard to provide context for LLms

## 🌟 Functionality

- 🔍 Search for files with specific extensions or all files
- 📋 Automatically copy file contents to clipboard
- 🚫 Skip specified directories

Certainly! Here's the installation instructions section in a code block, ready for you to copy as markdown:

## 🛠️ Installation

### macOS and Linux

1. Visit the [Releases](https://github.com/fargusplumdoodle/dump_dir/releases) page of the dump_dir repository.
2. Download the latest release for your operating system and architecture:
   - For macOS: `dump_dir_darwin_amd64` (Intel) or `dump_dir_darwin_arm64` (Apple Silicon)
   - For Linux: `dump_dir_linux_amd64` (64-bit) or `dump_dir_linux_arm64` (ARM)

3. Open a terminal and run the following commands, replacing `VERSION` with the version number and `ARCH` with your architecture (amd64 or arm64):

   ```bash
   # For macOS:
   curl -LO https://github.com/fargusplumdoodle/dump_dir/releases/download/vVERSION/dump_dir_darwin_ARCH
   chmod +x dump_dir_darwin_ARCH
   sudo mv dump_dir_darwin_ARCH /usr/local/bin/dump_dir

   # For Linux:
   curl -LO https://github.com/fargusplumdoodle/dump_dir/releases/download/vVERSION/dump_dir_linux_ARCH
   chmod +x dump_dir_linux_ARCH
   sudo mv dump_dir_linux_ARCH /usr/local/bin/dump_dir

### Windows

1. Visit the [Releases](https://github.com/fargusplumdoodle/dump_dir/releases) page of the dump_dir repository.
2. Download the latest release for Windows (look for `dump_dir_Windows_x86_64.zip`).
3. Extract the downloaded ZIP file.
4. Move the `dump_dir.exe` file to a directory of your choice.
5. Add the directory containing `dump_dir.exe` to your system's PATH:
   - Right-click on 'This PC' or 'My Computer' and select 'Properties'.
   - Click on 'Advanced system settings'.
   - Click on 'Environment Variables'.
   - Under 'System variables', find and select 'Path', then click 'Edit'.
   - Click 'New' and add the directory path where you placed `dump_dir.exe`.
   - Click 'OK' to close all dialogs.


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