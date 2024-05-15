package replacer

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb"
	"os"
	"strings"
)

func fillTitleToIdDict(fileName string, silent bool) (map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	pageInfo := make(map[string]string)
	scanner := bufio.NewScanner(file)
	bar := pb.New64(0)
	bar.ShowElapsedTime = true
	bar.ShowSpeed = true
	bar.NotPrint = silent
	bar.Start()

	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, "\t")
		pageInfo[data[1]] = data[0] // pageInfo[title] = id
		bar.Increment()
	}

	bar.Finish()
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanner:", err)
		return nil, err
	}
	return pageInfo, nil
}

func ReplaceTitleToId(pagePath string, currentPath string, silent bool) error {
	tempFilePath, err := createTempFileAndProcessData(pagePath, currentPath, silent)
	if err != nil {
		return err
	}

	if err := os.Rename(tempFilePath, currentPath); err != nil {
		fmt.Println("Error renaming file:", err)
		return err
	}

	return nil
}

func createTempFileAndProcessData(pagePath string, currentPath string, silent bool) (string, error) {
	tempFile, err := os.CreateTemp("", "tempfile.*.csv")
	if err != nil {
		fmt.Println("Error creating tempfile:", err)
		return "", err
	}
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			fmt.Println("Error closing tempfile:", err)
		}
	}(tempFile)

	file, err := os.Open(currentPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file) // close before return

	pageInfo, err := fillTitleToIdDict(pagePath, silent)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(tempFile)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer:", err)
		}
	}(writer) // flush buffer before return

	bar := pb.New64(0)
	bar.ShowElapsedTime = true
	bar.ShowSpeed = true
	bar.NotPrint = silent
	bar.Start()

	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, "\t")
		bar.Increment()
		pageId, ok := pageInfo[data[1]]
		if !ok {
			continue
		}
		data[1] = pageId

		_, err = writer.WriteString(strings.Join(data, "\t") + "\n")
		if err != nil {
			fmt.Println("Error writing file:", err)
			return "", err
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanner:", err)
		return "", err
	}

	bar.Finish()
	return tempFile.Name(), nil
}
