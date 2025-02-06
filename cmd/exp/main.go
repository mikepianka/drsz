package main

import (
	"fmt"
	"sync"
	"time"

	"math/rand"

	"slices"
)

func flush(rec int, buf map[int]string, output *[]int) {
	fmt.Printf("Flushing record %d\n", rec)
	*output = append(*output, rec)
	delete(buf, rec)
}

// flushBuf flushes keys from the buffer under the following conditions:
// if the buffer has length 1, flush the record
// if the buffer has length > 1, beginning at the first record check for consecutive keys
// flush if consecutive keys were found starting at the first record
func flushBuf(buf map[int]string, output *[]int) {
	if len(buf) == 0 {
		return
	}

	if len(buf) == 1 {
		// last record remaining in buf, flush it
		for k := range buf {
			flush(k, buf, output)
		}
		return
	}

	// pull keys from map and sort
	var keys []int
	for k := range buf {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	// we want to flush the buffer in ascending order, stopping once the first empty record is found
	for _, k := range keys {
		currentData := buf[k]
		if currentData == "" {
			// record is empty, stop flushing
			return
		}
		flush(k, buf, output)
	}
}

// writeToBuf writes a record to the buffer and flushes records if possible.
func writeToBuf(buf map[int]string, key int, val string, mu *sync.Mutex, output *[]int) {
	mu.Lock()
	// add record to buffer
	buf[key] = val
	// flush records from buffer if possible
	flushBuf(buf, output)
	mu.Unlock()
}

func processing(jobs int, buf map[int]string, output *[]int) {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(jobs)

	for i := 0; i < jobs; i++ {
		go func(id, amt int) {
			// simulate processing time
			time.Sleep(time.Duration(amt) * time.Second)
			fmt.Printf("Processed job %d\n", id)
			// store result in buffer
			writeToBuf(buf, id, fmt.Sprintf("%ds", amt), &mu, output)
			wg.Done()
		}(i, rand.Intn(10))
	}

	wg.Wait()
	fmt.Println("Processing complete")
}

func main() {

	jobs := 20
	buf := make(map[int]string)

	// fill slots in buffer to start
	for i := 0; i < jobs; i++ {
		buf[i] = ""
	}

	var output []int

	processing(jobs, buf, &output)
	fmt.Println(output)

	if len(output) != jobs {
		fmt.Println("Error: not all jobs were processed")
	}
}
