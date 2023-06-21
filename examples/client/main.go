package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/nilovartem/ipc-messager/pkg/client"
)

/*
This example shows that the server is capable of reading and transmitting
information to multiple clients at the same time (in this case doubled numbers)
*/
func main() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		<-results
	}

}
func worker(id int, jobs <-chan int, results chan<- int) {
	c, err := client.Connect("/tmp/server.sock", time.Millisecond*200) //optional timeout, or constant
	if err == nil {
		for j := range jobs {
			fmt.Println("worker", id, "started  job", j)
			c.Send([]byte(strconv.Itoa(j)))
			data, ok := c.Receive()
			if ok {
				var buffer bytes.Buffer = *bytes.NewBuffer(data)
				fmt.Println(buffer.String())
				fmt.Println("worker", id, "finished job", j, "result ", buffer.String())
				result, _ := strconv.Atoi(buffer.String())
				results <- result
			}
		}
	}
}
