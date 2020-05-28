package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	var producer = func(outch chan<- string) {
		defer close(outch)
		stdin := bufio.NewReader(os.Stdin)
	readLoop:
		for {
			switch url, err := stdin.ReadString('\n'); err {
			case nil:
				outch <- strings.ReplaceAll(url, "\n", "")
			case (io.EOF):
				break readLoop
			default:
				panic(err)
			}
		}
	}
	var investigator = func(inch <-chan string, outch chan<- int, wg *sync.WaitGroup) {
		defer wg.Done()
		for url := range inch {
			resp, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			var c = strings.Count(string(body), "Go")
			fmt.Printf("Count for %s: %d\n", url, c)
			outch <- c
		}
	}
	var dispatcher = func(inch <-chan string, outch chan<- int) {
		defer close(outch)
		const POOL_SIZE = 5
		var (
			spawned = 0
			tasks   = make(chan string)
			wg      = sync.WaitGroup{}
		)
		for url := range inch {
			if spawned < POOL_SIZE {
				spawned += 1
				wg.Add(1)
				go investigator(tasks, outch, &wg)
			}
			tasks <- url
		}
		close(tasks)
		wg.Wait()
	}
	var consumer = func(inch <-chan int) int {
		n := 0
		for x := range inch {
			n += x
		}
		return n
	}
	// каналы
	var urls = make(chan string)
	var counts = make(chan int)
	// конвейер
	// пушит урлы в urls
	go producer(urls)
	// слушает urls, создает горутины для обработки, отправляет ответы от них дальше в counts
	go dispatcher(urls, counts)
	// редьюсит ответы
	total := consumer(counts)
	fmt.Println("Total:", total)
}
