package main

import (
	"bitbucket.org/minamartinteam/myevents/src/bookingservice/rest"
	"bitbucket.org/minamartinteam/myevents/src/bookingservice/listener"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go listener.ProcessEvents(&wg)
	go rest.ServeAPI(&wg)

	wg.Wait()
}
