package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
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

	linesChunkLen := 64 * 1024
	linesChunkPoolAllocated := int64(0)
	iterator := 0
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]string, 0, linesChunkLen)
		atomic.AddInt64(&linesChunkPoolAllocated, 1)
		return lines
	}}
	lines := linesPool.Get().([]string)[:0]

	br := bufio.NewReader(file)

	for {
		line, _, err  := br.ReadLine()
        if err == io.EOF {
            break
		}
		if iterator < 2000 {
			lines = append(lines, string(line))
			iterator++	
 		} else {
			wg.Add(1)
			linesToProcess := lines
			go func() {
				stdoutFile := "./huge.log"
				file, err := os.OpenFile(stdoutFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

				if err != nil {
					fmt.Printf("open file failed:", err.Error())
					return
				}

				defer file.Close()

				w := bufio.NewWriter(file) 

				for i := 0; i < len(linesToProcess); i++ {
					w.WriteString(linesToProcess[i])
					w.WriteString("\n")
				} 

				linesPool.Put(linesToProcess)
				w.Flush()
				wg.Done()
			}()
			lines = linesPool.Get().([]string)[:0]
			iterator = 1
		}
		
	}

	 if iterator > 0 {
		wg.Add(1) 
		go func() {
			stdoutFile := "./huge.log"
			file, err := os.OpenFile(stdoutFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

			if err != nil {
				fmt.Printf("open file failed:", err.Error())
				return
			}

			defer file.Close()

			w := bufio.NewWriter(file) 

			for i := 0; i < len(lines); i++ {
				w.WriteString(lines[i])
				w.WriteString("\n")
			} 
			w.Flush()
			wg.Done()
		}()
	}

	fmt.Printf("used time: %v\n", time.Since(start))
	wg.Wait()
}
