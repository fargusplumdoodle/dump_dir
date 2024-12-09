package prompt

import "encoding/xml"

func GenerateXMLPrompt(p Prompt) XMLPrompt {
	var fileList []XMLFile
	for _, file := range p.ProcessedFiles {
		fileList = append(fileList, XMLFile{Path: file.Path})
	}

	var processed []XMLFile
	for _, file := range p.ProcessedFiles {
		processed = append(processed, XMLFile{
			Path:     file.Path,
			Contents: file.Contents,
		})
	}

	var skippedLarge []XMLFile
	for _, file := range p.SkippedLarge {
		skippedLarge = append(skippedLarge, XMLFile{Path: file.Path})
	}

	var skippedBinary []XMLFile
	for _, file := range p.SkippedBinary {
		skippedBinary = append(skippedBinary, XMLFile{Path: file.Path})
	}

	return XMLPrompt{
		FileList:      fileList,
		Processed:     processed,
		SkippedLarge:  skippedLarge,
		SkippedBinary: skippedBinary,
	}
}
func MarshalXMLPrompt(xmlPrompt XMLPrompt) (string, error) {
	output, err := xml.MarshalIndent(xmlPrompt, "", "  ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(output), nil
}
