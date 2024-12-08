package src

func isDirectory(path string) (bool, error) {
	fileInfo, err := OsStat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}
