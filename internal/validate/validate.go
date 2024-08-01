package validate

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb"
	"os"
	"strconv"
	"strings"
)

func createRedirectDict(fileName string, silent bool) (map[uint32]uint32, error) {
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
		sourcePageId, srcErr := strconv.Atoi(data[0])
		targetPageId, tgtErr := strconv.Atoi(data[1])
		if srcErr != nil || tgtErr != nil {
			fmt.Println("Error converting line to int:", data, err)
			return nil, err
		}

		redirectInfo[uint32(sourcePageId)] = uint32(targetPageId)
		bar.Increment()
	}

	bar.Finish()
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanner:", err)
		return nil, err
	}
	return redirectInfo, nil
}

func Redirect(linksPath, redirectPath string, silent bool) error {
	tempFilePath, err := createTempFileAndProcessData(linksPath, redirectPath, silent)
	if err != nil {
		return err
	}

	if err := os.Rename(tempFilePath, linksPath); err != nil {
		fmt.Println("Error renaming file:", err)
		return err
	}

	return nil
}

func createTempFileAndProcessData(linksPath, redirectPath string, silent bool) (string, error) {
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

	redirectInfo, err := createRedirectDict(redirectPath, silent)
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

		sourcePageId, srcErr := strconv.Atoi(data[0])
		targetPageId, tgtErr := strconv.Atoi(data[1])
		if srcErr != nil || tgtErr != nil {
			fmt.Println("Error converting line to int:", data, err)
			return "", err
		}

		// replace redirect target page to page which its redirect
		if redirectTo, ok := redirectInfo[uint32(targetPageId)]; ok {
			targetPageId = int(redirectTo)
		}

		_, err = writer.WriteString(strconv.Itoa(sourcePageId) + "\t" + strconv.Itoa(targetPageId) + "\n")
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
