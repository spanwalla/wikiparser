package combiner

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Links struct {
	incomingLinks []int
	outgoingLinks []int
}

func convertIntSliceToString(slice []int, sep string) string {
	var sb strings.Builder
	for i, num := range slice {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(strconv.Itoa(num))
	}
	if len(slice) == 0 {
		sb.WriteString("-")
	}
	return sb.String()
}

func addLink(pageLinks map[int]Links, sourcePageId, targetPageId int) {
	pl, ok := pageLinks[sourcePageId]
	if !ok {
		pl = Links{}
	}
	pl.outgoingLinks = append(pl.outgoingLinks, targetPageId)
	pageLinks[sourcePageId] = pl

	pl, ok = pageLinks[targetPageId]
	if !ok {
		pl = Links{}
	}
	pl.incomingLinks = append(pl.incomingLinks, sourcePageId)
	pageLinks[targetPageId] = pl
}

func CombineLinks(pageLinksPath string, silent bool) error {
	linksFile, err := os.Open(pageLinksPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(linksFile) // close before return

	// open output file
	linksFileName := filepath.Base(pageLinksPath)
	output, err := os.Create("output/combined-" + linksFileName)
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

	scannerLinks := bufio.NewScanner(linksFile)
	writer := bufio.NewWriter(output)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer:", err)
		}
	}(writer) // flush buffer before return

	pageLinks := make(map[int]Links)

	bar := pb.New64(0)
	bar.ShowElapsedTime = true
	bar.ShowSpeed = true
	bar.NotPrint = silent
	bar.Start()

	for scannerLinks.Scan() {
		line := scannerLinks.Text()
		data := strings.Split(line, "\t")
		sourcePageId, srcErr := strconv.Atoi(data[0])
		targetPageId, tgtErr := strconv.Atoi(data[1])

		if srcErr != nil || tgtErr != nil {
			fmt.Println("Error converting line to int:", data, err)
			return err
		}

		addLink(pageLinks, sourcePageId, targetPageId)
		bar.Increment()
	}

	if err := scannerLinks.Err(); err != nil {
		fmt.Println("Error scanner:", err)
		return err
	}

	for pageId, links := range pageLinks {
		_, err = writer.WriteString(strconv.Itoa(pageId) + "\t" + convertIntSliceToString(links.incomingLinks, ",") + "\t" + convertIntSliceToString(links.outgoingLinks, ",") + "\n")
		if err != nil {
			fmt.Println("Error writing file:", err)
			return err
		}
		bar.Increment()
	}

	bar.Finish()
	return nil
}
