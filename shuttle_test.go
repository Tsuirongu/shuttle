package shuttle

import (
	"fmt"
	"testing"
	"time"
)

func upload(list []string) {
	fmt.Println(len(list))
}

func TestShuttle(t *testing.T) {
	pool := New(
		WithFunc(func(m map[string]interface{}) error {
			var list []string
			for k := range m {
				list = append(list, k)
			}
			upload(list)
			return nil
		}),
		WithDuration(1*time.Second),
		WithMaxSize(20),
	)

	defer func() {
		if err := recover(); err != nil {
			pool.Exit()
		}
	}()
	go func() {

		// add 1
		pool.AddData("0", "value")
		time.Sleep(2 * time.Second)
		// add 2
		pool.AddData("1", "value")
		pool.AddData("2", "value")
		time.Sleep(2 * time.Second)
		pool.AddData("1", "value")
		pool.AddData("2", "value")
		pool.AddData("3", "value")
		pool.AddData("4", "value")
		time.Sleep(2 * time.Second)
		pool.AddData("1", "value")
	}()
	time.Sleep(7 * time.Second)

	// in case of panic, the retained data will be processed in time because of the Exit() above
	panic("")
}
