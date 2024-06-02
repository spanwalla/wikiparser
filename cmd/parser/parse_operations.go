package parser

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/cheggaaa/pb"
	"os"
	"path/filepath"
	"strings"
)

func processValue(data string, table string) ([]string, error) {
	info := tableInfo[table]
	matches := info.re.FindStringSubmatch(data)
	if matches == nil {
		return nil, nil
	}

	for i, column := range info.columns {
		matches[i] = matches[column+1]
	}
	return matches[:len(info.columns)], nil
}

func processString(data string) ([][]string, error) {
	tableName, values, found := strings.Cut(data, "` VALUES (")
	if !found {
		return nil, errors.New("incorrect string (doesn't contain '` VALUES (')")
	}

	// format "INSERT INTO `table` VALUES (...),(...),(...);"
	tableName = strings.TrimPrefix(tableName, "INSERT INTO `")
	values = strings.TrimSuffix(values, ");")
	valuesList := strings.Split(values, "),(")

	// call process value and return results
	var result [][]string
	for _, value := range valuesList {
		extractedValues, err := processValue(value, tableName)
		if err != nil {
			return nil, err
		}
		if extractedValues != nil {
			result = append(result, extractedValues)
		}
	}

	return result, nil
}

func ProcessFile(input string, silent bool) error {
	// open .gz file
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	// unpack .gz file
	gz, err := gzip.NewReader(file)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func(gz *gzip.Reader) {
		err := gz.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(gz)

	scanner := bufio.NewScanner(gz)
	const bufferSize = 1024 * 1024
	scanner.Buffer(make([]byte, bufferSize), bufferSize)

	// open output file
	fileNameWithExtension := filepath.Base(input)
	fileNameWithoutExtension := strings.TrimSuffix(fileNameWithExtension, filepath.Ext(fileNameWithExtension))

	output, err := os.Create("output/" + fileNameWithoutExtension + ".csv")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(output)

	writer := bufio.NewWriter(output)

	bar := pb.New64(0)
	bar.ShowElapsedTime = true
	bar.ShowSpeed = true
	bar.NotPrint = silent
	bar.Start()

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "INSERT INTO") {
			// call process line function
			data, err := processString(line)
			if err != nil {
				fmt.Println("Error processing:", err)
				return err
			}

			// use results of process line and write it
			for _, row := range data {
				_, err := writer.WriteString(strings.Join(row, "\t") + "\n")
				if err != nil {
					fmt.Println("Error writing file:", err)
					break
				}
			}
			bar.Increment()
		}
	}

	// make sure that all data written before closing
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing buffer:", err)
		return err
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanner:", err)
		return err
	}

	bar.Finish()
	return nil
}
