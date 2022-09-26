# Shuttle

> Introduction: The warehouse is a middleware summed up by individuals based on business experience. The main usage scenarios are as follows:
> 
> 1. A service with high concurrency but low database IO.
> 
> 2. A service that need to keep multiple concurrent data in sync.
> 
> 3. A service that require high-concurrency interfaces for high-cost single reporting.


Use the data pool to temporarily store high-concurrency data, and then use the shuttle function to process it in batches. After processing, the data pool is emptied and waits for the next batch of data to enter.

The project is currently go 1.18 version.

### Instructions
> 1. Initialize settings
> 
> > 1. set Shuttle function
> >
> >  WithFunc(func(m map[string]interface{}) error）
> >
> > The parameter map is the value in the entire data pool. The user can configure the handler for processing here. After processing, the data pool is cleared.
> 
> > 2. set duration
> >
> > WithDuration(time.Duration)
> >
> > Set how often to trigger the shuttle function.
> 
> > 3. Set the maximum capacity of the data pool
> >
> > WithMaxSize(int)
> >
> > After the data in the data pool reaches a certain amount, the cleaning operation will be performed immediately
> 
> 2. Just Do It
> 
> > 1. Add data to the datapool
> > 
> > AddData(key string, value interface)
> > 
> > Add the required data to the data pool, where the value of the key will be deduplicated.
> 
> > 2. Move data out of the dataPool
> >
> > DelData(key string)
> 
> > 3. Exit processing
> >
> > Exit()
> >
> > If it is sensitive to traffic (data is indispensable), if there is a panic or the service is shut down in the middle, you need to use this method to terminate the control coroutine. At the same time, this method will also trigger the configured shuttle function for the last time, so that the temporary storage But as to time the data gets processed.

### Example
```go
package main

import (
	"fmt"
	"time"

	"Tsuirongu/shuttle"
)

func main() {
	a := shuttle.New(
		shuttle.WithFunc(func(m map[string]interface{}) error {
		count := 0
		for _, _ = range m {
			count += 1
		}
		fmt.Println(count)

		return nil
	}), 
	shuttle.WithDuration(60*time.Second), 
	shuttle.WithMaxSize(100))

	a.AddData("key", "value")
	defer func() {
		if err:=recover(); err!=nil {
			a.Exit()
		}
	}()
}
```
