package util

import (
	"sync"
)

var count = 5

func DoMap(data []interface{}, mapper func(interface{}) interface{}) []interface{} {
	tokens := make(chan bool, count)
	result := make([]interface{}, len(data))
	wg := sync.WaitGroup{}
	wg.Add(len(data))
	for i := 0; i < count; i++ {
		tokens <- true
	}
	for i := 0; i < len(data); i++ {
		go func(index int) {
			t := <-tokens
			defer wg.Done()
			result[index] = mapper(data[index])
			tokens <- t
		}(i)
	}
	wg.Wait()
	return result
}
