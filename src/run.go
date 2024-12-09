package src

import "fmt"

func Run(args []string, config RunConfig) error {
	if !ValidateArgs(args) {
		PrintUsage()
		return nil
	}

	cliConfig, err := ParseArgs(args)
	if err != nil {
		PrintUsage()
		return fmt.Errorf("error parsing arguments: %v", err)
	}

	switch cliConfig.Action {
	case "help":
		PrintUsage()
		return nil
	case "version":
		PrintVersion(config)
		return nil
	case "dump_dir":
		return performDumpDir(cliConfig, config)
	default:
		return fmt.Errorf("unknown action: %s", cliConfig.Action)
	}
}

func performDumpDir(cliArgumentsConfig Config, runConfig RunConfig) error {
	configLoader := NewConfigLoader(runConfig.Fs)
	config, err := configLoader.LoadAndMergeConfig(cliArgumentsConfig)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	fileFinder := NewFileFinder(config, runConfig.Fs)
	fileProcessor := NewFileProcessor(runConfig.Fs, config)

	filePaths := fileFinder.DiscoverFiles()
	processedFiles := fileProcessor.ProcessFiles(filePaths)
	stats := CalculateStats(processedFiles)
	PrintDetailedOutput(stats, runConfig)
	return nil
}

func PrintVersion(cfg RunConfig) {
	println("dump_dir version:", cfg.Version)
	println("commit:", cfg.Commit)
	println("built at:", cfg.Date)
}
