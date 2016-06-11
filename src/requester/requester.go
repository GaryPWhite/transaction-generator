package requester

import (
	"fmt"
	"time"
	"sync"
	"net/http"
)

type duration struct {
	durations map[string]int
	m sync.Mutex
}

var d duration
var client = &http.Client{}

type timeoutError int

func (t timeoutError) Error() string {
	return fmt.Sprintf("Timeout Exceeded : %d", t)
}

func runRequest(method, url string) {
	// create client to run requests on
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	start := time.Now()
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	// for concurrent operations
	dur := int(time.Since(start)/time.Millisecond)
	d.m.Lock()
	d.durations[string(len(d.durations))] = dur
	d.m.Unlock()
	return // for clarity
}

func MakeRequests(method, url string, persecond, total, timeout int) (map[string]int, error) {
	d.durations = make(map[string]int)
	expectedlen := total;
	start := time.Now()
	for total > 0{ // run until total is exhausted
		for i := 0; ((i < persecond) && (i < total)); i++ {
			go runRequest(method, url); // execute requests on url with method
		}
		time.Sleep(time.Second); // sleep for a second (to ensure per/second)
		total -= persecond;
	}
	for (int(time.Since(start)) < (int(time.Millisecond) * timeout)) || (expectedlen < len(d.durations)) { // all threads successfully finished
		fmt.Printf("waiting %d", len(d.durations));
		fmt.Printf("timediff start : %d, now: %d", int(time.Since(start)/time.Millisecond), (int(time.Millisecond)*timeout))
		time.Sleep(time.Millisecond * 300) // wait a third of a second at a time
	}
	if ((int(time.Since(start)) < (timeout)) && (!(expectedlen < len(d.durations)))) {
		return nil, timeoutError(timeout);
	}
	return d.durations, nil
}
