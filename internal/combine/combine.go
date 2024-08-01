package combine

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"wikiparser/internal/pkg/convert"
)

type RelatedLinks struct {
	incomingLinks []uint32
	outgoingLinks []uint32
}

func addLink(pageLinks map[uint32]RelatedLinks, sourcePageId, targetPageId uint32) {
	pl, ok := pageLinks[sourcePageId]
	if !ok {
		pl = RelatedLinks{}
	}
	pl.outgoingLinks = append(pl.outgoingLinks, targetPageId)
	pageLinks[sourcePageId] = pl

	pl, ok = pageLinks[targetPageId]
	if !ok {
		pl = RelatedLinks{}
	}
	pl.incomingLinks = append(pl.incomingLinks, sourcePageId)
	pageLinks[targetPageId] = pl
}

func Links(pageLinksPath string, silent bool) error {
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

	pageLinks := make(map[uint32]RelatedLinks)

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

		addLink(pageLinks, uint32(sourcePageId), uint32(targetPageId))
		bar.Increment()
	}

	if err := scannerLinks.Err(); err != nil {
		fmt.Println("Error scanner:", err)
		return err
	}

	for pageId, links := range pageLinks {
		_, err = writer.WriteString(strconv.Itoa(int(pageId)) + "\t" + convert.Uint32SliceToString(links.incomingLinks, ",") + "\t" + convert.Uint32SliceToString(links.outgoingLinks, ",") + "\n")
		if err != nil {
			fmt.Println("Error writing file:", err)
			return err
		}
		bar.Increment()
	}

	bar.Finish()
	return nil
}
