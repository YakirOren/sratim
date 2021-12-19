package progress

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type WriteCounter struct {
	Total uint64
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s", humanize.Bytes(wc.Total))
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// Loader prints a 30 seconds progress bar.
func Loader(seconds int) {
	// Print loading bar progress
	for i := 0; i < seconds+1; i++ {
		fmt.Printf("\r%-2s[%s]", strconv.Itoa(seconds-i), strings.Repeat("=", i)+">"+strings.Repeat(".", seconds-i))
		// Sleep for 1 second
		time.Sleep(1 * time.Second)
	}
	fmt.Println("")
}
