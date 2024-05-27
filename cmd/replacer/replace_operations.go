package replacer

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb"
	"os"
	"strconv"
	"strings"
)

func createTitleToIdDict(fileName string, silent bool) (map[string]uint32, error) {
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

	pageInfo := make(map[string]uint32)
	scanner := bufio.NewScanner(file)
	bar := pb.New64(0)
	bar.ShowElapsedTime = true
	bar.ShowSpeed = true
	bar.NotPrint = silent
	bar.Start()

	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, "\t")
		pageId, err := strconv.Atoi(data[0])
		if err != nil {
			fmt.Println("Error converting line to int:", data, err)
			return nil, err
		}
		pageInfo[data[1]] = uint32(pageId) // pageInfo[title] = id
		bar.Increment()
	}

	bar.Finish()
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanner:", err)
		return nil, err
	}
	return pageInfo, nil
}

func createRedirectDict(fileName string, pageInfo map[string]uint32, silent bool) (map[uint32]uint32, error) {
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

	redirectInfo := make(map[uint32]uint32)
	scanner := bufio.NewScanner(file)
	bar := pb.New64(0)
	bar.ShowElapsedTime = true
	bar.ShowSpeed = true
	bar.NotPrint = silent
	bar.Start()

	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, "\t")
		sourcePageId, err := strconv.Atoi(data[0])
		if err != nil {
			fmt.Println("Error converting line to int:", data, err)
			return nil, err
		}
		targetPageId, ok := pageInfo[data[1]]
		if !ok {
			continue
		}

		redirectInfo[uint32(sourcePageId)] = targetPageId
		bar.Increment()
	}

	bar.Finish()
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanner:", err)
		return nil, err
	}
	return redirectInfo, nil
}

func ReplaceTitleToId(pagePath string, redirectPath string, currentPath string, silent bool) error {
	tempFilePath, err := createTempFileAndProcessData(pagePath, redirectPath, currentPath, silent)
	if err != nil {
		return err
	}

	if err := os.Rename(tempFilePath, currentPath); err != nil {
		fmt.Println("Error renaming file:", err)
		return err
	}

	return nil
}

func createTempFileAndProcessData(pagePath, redirectPath, linksPath string, silent bool) (string, error) {
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

	file, err := os.Open(linksPath)
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

	pageInfo, err := createTitleToIdDict(pagePath, silent)
	if err != nil {
		return "", err
	}

	redirectInfo, err := createRedirectDict(redirectPath, pageInfo, silent)
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

		// replace redirect target page to page which its redirect
		if redirectTo, ok := redirectInfo[pageId]; ok {
			pageId = redirectTo
		}
		data[1] = strconv.Itoa(int(pageId))

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
