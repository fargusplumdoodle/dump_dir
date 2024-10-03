## Plan for Implementing a Flexible Format System Supporting Multiple Output Formats

### Overview

To enhance the existing application’s output capabilities, we aim to introduce support for multiple formats—specifically XML, plain text, and Markdown. The current implementation primarily supports XML but is considered messy and not easily extensible. This plan outlines the steps to restructure the application to accommodate multiple formats in a clean, maintainable, and scalable manner. The key objectives include:

- **Modular Format Handling:** Introduce a `format` package with a well-defined interface to manage different formats.
- **Simplicity and Clarity:** Utilize straightforward constants to manage available formats without complex registration mechanisms.
- **Extensibility:** Ensure that adding new formats in the future requires minimal changes.
- **CLI Integration:** Update the command-line interface (CLI) to allow users to specify the desired format dynamically, sourcing available options from the `format` package.
- **Directory Renaming:** Rename the existing `prompts` directory to better reflect its new role within the application.

### 1. Directory Structure Changes

**Current Structure:**
```
src/
├── args.go
├── colors.go
├── file_finder.go
├── file_parser.go
├── ignore.go
├── output.go
├── prompt/
│   ├── model.go
│   └── prompt.go
├── stats.go
└── types.go
```

**Proposed Changes:**

- **Rename `prompt` Directory:**
  - **New Name:** `format`
  - **Reason:** The current name `prompt` is misleading as the directory will now encapsulate various output formats. Renaming it to `format` better represents its purpose.

**Updated Structure:**
```
src/
├── args.go
├── colors.go
├── file_finder.go
├── file_parser.go
├── ignore.go
├── output.go
├── format/
│   ├── model.go
│   ├── xml_formatter.go
│   ├── plain_formatter.go
│   ├── markdown_formatter.go
│   └── formatter.go
├── stats.go
└── types.go
```

### 2. Creating the `format` Package and Defining the Interface

**Objective:** Establish a centralized `format` package that defines a common interface for all formats, facilitating easy integration and extension.

**Steps:**

1. **Create `format` Package:**
   - Move existing `prompt/model.go` and `prompt/prompt.go` to the new `format` directory.

2. **Define `Formatter` Interface:**
   - **Location:** `src/format/formatter.go`
   - **Purpose:** Standardize the methods required for any format.
   - **Interface Definition:**
     ```go
     package format

     type Formatter interface {
         Format(stats Stats) (string, error)
         Name() string
     }
     ```

3. **Documentation:**
   - Provide clear documentation within the `format` package to guide developers on implementing new formats.

### 3. Implementing Specific Formatters

**Objective:** Create concrete implementations of the `Formatter` interface for each desired format—XML, plain text, and Markdown.

**Steps:**

1. **XML Formatter:**
   - **File:** `src/format/xml_formatter.go`
   - **Implements:** `Formatter`
   - **Responsibilities:**
     - Serialize `Stats` into XML.
     - Utilize existing logic with improvements for cleanliness.

2. **Plain Text Formatter:**
   - **File:** `src/format/plain_formatter.go`
   - **Implements:** `Formatter`
   - **Responsibilities:**
     - Convert `Stats` into human-readable plain text.
     - Ensure readability and simplicity.

3. **Markdown Formatter:**
   - **File:** `src/format/markdown_formatter.go`
   - **Implements:** `Formatter`
   - **Responsibilities:**
     - Structure `Stats` data into Markdown syntax.
     - Facilitate easy viewing in Markdown-compatible viewers.

4. **Define Available Formats Using Constants:**
   - **File:** `src/format/formatter.go`
   - **Purpose:** Define constants representing available formats.
   - **Implementation:**
     ```go
     package format

     import (
         "fmt"
         "sort"
     )

     const (
         FormatXML      = "xml"
         FormatPlain    = "plain"
         FormatMarkdown = "markdown"
     )

     // GetAvailableFormats returns a sorted list of available format names.
     func GetAvailableFormats() []string {
         formats := []string{FormatXML, FormatPlain, FormatMarkdown}
         sort.Strings(formats)
         return formats
     }

     // GetFormatter returns the Formatter implementation based on the format name.
     func GetFormatter(name string) (Formatter, error) {
         switch name {
         case FormatXML:
             return &XMLFormatter{}, nil
         case FormatPlain:
             return &PlainFormatter{}, nil
         case FormatMarkdown:
             return &MarkdownFormatter{}, nil
         default:
             return nil, fmt.Errorf("unsupported format: %s", name)
         }
     }
     ```

### 4. Refactoring Existing Code to Utilize the `Formatter` Interface

**Objective:** Modify the existing application logic to leverage the new `Formatter` interface, decoupling format generation from specific implementations.

**Steps:**

1. **Update `output.go`:**
   - Replace existing XML-specific code with calls to the `Formatter` interface.
   - **Example Modification:**
     ```go
     package src

     import (
         "fmt"
         "github.com/fargusplumdoodle/dump_dir/src/format"
     )

     func PrintDetailedOutput(stats Stats, formatter format.Formatter) {
         formattedOutput, err := formatter.Format(stats)
         if err != nil {
             fmt.Println("❌ Error formatting output:", err)
             return
         }

         fmt.Println(formattedOutput)
         if copyToClipboard(formattedOutput) {
             fmt.Println(BoldGreen("✅ Output has been copied to clipboard.\n"))
         } else {
             fmt.Println("❌ Failed to copy output to clipboard.")
         }
     }
     ```

2. **Modify `PrintDetailedOutput` Function:**
   - Accept a `format.Formatter` parameter to determine the output format dynamically.

3. **Adjust Existing Format-Specific Functions:**
   - Ensure that functions like `prompt.GenerateXMLPrompt` and related logic are encapsulated within their respective formatter implementations within the `format` package.

### 5. Updating the Command-Line Interface (CLI)

**Objective:** Introduce a new CLI option to allow users to specify the desired format, sourcing available options from the `format` package.

**Steps:**

1. **Define New CLI Flag:**
   - **Flag Name:** `--format` or `-f`
   - **Usage:** Specify the desired output format (e.g., `--format=xml`, `--format=plain`, `--format=markdown`).

2. **Integrate with `format` Package:**
   - Retrieve available formats from the `format` package’s constants.
   - **Example Implementation:**
     ```go
     // In args.go or relevant CLI parsing file
     package src

     import (
         "flag"
         "fmt"
         "github.com/fargusplumdoodle/dump_dir/src/format"
     )

     type Config struct {
         // Existing fields...
         OutputFormat string
         // Other fields...
     }

     func ParseArgs(args []string) (Config, error) {
         config := Config{
             // Initialize default values
             OutputFormat: format.FormatXML, // Default format
             // Other initializations...
         }

         // Define flags
         formatFlag := flag.String("format", format.FormatXML, "Specify output format. Available formats: "+format.JoinFormats())
         flag.StringVar(formatFlag, "f", format.FormatXML, "Specify output format. Available formats: "+format.JoinFormats())

         // Parse flags
         flag.CommandLine.Parse(args)

         // Validate output format
         if _, err := format.GetFormatter(*formatFlag); err != nil {
             return config, err
         }
         config.OutputFormat = *formatFlag

         // Process other arguments...

         return config, nil
     }

     // In format/formatter.go, add a helper function to join formats
     func JoinFormats() string {
         available := GetAvailableFormats()
         return fmt.Sprintf("%s", joinWithComma(available))
     }

     func joinWithComma(items []string) string {
         result := ""
         for i, item := range items {
             if i > 0 {
                 result += ", "
             }
             result += item
         }
         return result
     }
     ```

3. **Provide Help and Usage Information:**
   - Dynamically list available formats in the help section by querying the `format` package’s constants.
   - **Example Help Output:**
     ```
     -f, --format string    Specify output format. Available formats: xml, markdown, plain
     ```

4. **Handle Default Output Format:**
   - Define a default output format (e.g., XML) if the user does not specify one.

### 6. Ensuring Extensibility for Future Formats

**Objective:** Design the system to allow seamless addition of new formats without modifying existing core logic.

**Steps:**

1. **Interface Compliance:**
   - Ensure that any new formatter adheres to the `Formatter` interface defined in the `format` package.

2. **Adding New Formatter:**
   - Create a new formatter file within the `format` package (e.g., `new_formatter.go`).
   - Implement the `Formatter` interface.
   - Update the `GetFormatter` function in `formatter.go` to include the new format.

3. **Minimal Codebase Changes:**
   - Adding a new format should only require creating a new formatter file and updating the `GetFormatter` function with a new case for the format.

4. **Documentation and Guidelines:**
   - Maintain comprehensive documentation within the `format` package to guide developers on implementing and adding new formats.

### 7. Testing and Validation

**Objective:** Ensure that the new format system functions correctly across all supported formats and maintains the application's integrity.

**Steps:**

1. **Unit Tests:**
   - Write unit tests for each `Formatter` implementation to verify correct formatting.

2. **Integration Tests:**
   - Test the end-to-end flow of generating output in different formats through the CLI.

3. **CLI Tests:**
   - Validate that the CLI correctly handles the `--format` flag, including default behavior and error handling for unsupported formats.

4. **Performance Testing:**
   - Ensure that the introduction of multiple formats does not adversely affect application performance.

### 8. Documentation and Developer Guidelines

**Objective:** Provide clear guidance for both users and developers on using and extending the format system.

**Steps:**

1. **User Documentation:**
   - Update the application's README and help sections to reflect the new format options.
   - Provide examples demonstrating how to use different formats.
   - **Example README Section:**
     ```markdown
     ## Output Formats

     The application supports multiple output formats. You can specify the desired format using the `--format` or `-f` flag.

     **Available Formats:**
     - `xml`: Outputs data in XML format.
     - `plain`: Outputs data in plain text.
     - `markdown`: Outputs data in Markdown format.

     **Examples:**
     ```bash
     dump_dir --format=xml ./project
     dump_dir -f markdown ./project
     dump_dir --format=plain ./README.md ./main.go
     ```
     ```

2. **Developer Documentation:**
   - Document the `format` package’s structure, interface, and procedures for adding new formats.
   - Include code comments and usage examples within the `format` package files.

### 9. Migration Strategy

**Objective:** Transition the existing application to the new format system without disrupting current functionality.

**Steps:**

1. **Parallel Implementation:**
   - Implement the new `format` package alongside the existing codebase.
   - Ensure that existing XML output functionality remains operational during the transition.

2. **Gradual Refactoring:**
   - Incrementally replace XML-specific output code with the new `Formatter` interface.
   - Verify each step through testing before proceeding.

3. **Deprecation and Cleanup:**
   - Once the new system is fully integrated and tested, deprecate and remove the old XML-specific code.
   - Clean up any redundant files or logic to maintain codebase cleanliness.

### 10. Conclusion

By restructuring the application to include a dedicated `format` package with a standardized `Formatter` interface, we achieve a modular and extensible system capable of supporting multiple output formats with ease. This design not only addresses the current need to support XML, plain text, and Markdown but also lays the groundwork for incorporating additional formats in the future with minimal effort. Integrating the format selection into the CLI ensures user flexibility, while the use of constants simplifies format management, promoting maintainability and scalability.

Implementing this plan will result in a cleaner codebase, improved maintainability, and enhanced user experience through versatile output options.