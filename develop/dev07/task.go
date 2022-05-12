package main

import (
	"fmt"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

// channelProcessing - функция обработки канал, если канал закрыт то пишем в single nil
func channelProcessing(channel <-chan interface{}, single chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case _, open := <-channel:
			if !open {
				single <- nil
				break
			}
		}
	}
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	single := make(chan interface{})
	wg := &sync.WaitGroup{} // для ожидания обработки всех каналов и закрытия single канала
	wg.Add(len(channels))

	for _, done := range channels {
		go channelProcessing(done, single, wg)
	}

	go func() {
		wg.Wait()
		close(single)
	}()

	return single
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-or(
		sig(2*time.Hour),
		sig(3*time.Second),
		//sig(1*time.Second),
		sig(2*time.Second),
		sig(2*time.Minute),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("fone after %v ", time.Since(start))
}
