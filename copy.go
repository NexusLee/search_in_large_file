package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var(
	wg = sync.WaitGroup{}
)
func main() {
	start := time.Now()
	stdinFile := "./stdout.log"
	
 	file, err := os.Open(stdinFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	br := bufio.NewReader(file)

	stdoutFile := "./huge.log"
	outfile, err1 := os.OpenFile(stdoutFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	w := bufio.NewWriter(outfile) 

	if err1 != nil {
		fmt.Printf("open file failed:", err.Error())
		return
	}

	defer outfile.Close()

	for {
		line, _, err  := br.ReadLine()
        if err == io.EOF {
            break
		}

		w.WriteString(string(line))
		w.WriteString("\n")

		w.Flush()
		
	}

	fmt.Printf("used time: %v\n", time.Since(start))
}
