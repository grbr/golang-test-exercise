package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	//
	var unexpected = func(err error) {
		if err != nil {
			panic(err)
		}
	}
	//
	var investigateURL = func(url string, out chan<- int) {
		resp, e := http.Get(url)
		unexpected(e)
		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)
		unexpected(e)
		var c = strings.Count(string(body), "Go")
		fmt.Printf("Count for %s: %d\n", url, c)
		out <- c
	}

	const POOL_SIZE = 5
	var (
		// результат
		total = 0
		// живых горутин
		thingsToDo = 0
		// случился ли конец URLам
		isEOF = false
		//
		ch = make(chan int)
	)
	stdin := bufio.NewReader(os.Stdin)

mainCicle:
	for {
		if isEOF == false {
			switch url, err := stdin.ReadString('\n'); err {
			case nil:
				thingsToDo += 1
				go investigateURL(strings.ReplaceAll(url, "\n", ""), ch)
			case (io.EOF):
				isEOF = true
			default:
				panic(err)
			}
		}
		if thingsToDo == POOL_SIZE || (isEOF && thingsToDo > 0) {
			v := <-ch
			total += v
			thingsToDo -= 1
		}
		if isEOF && thingsToDo == 0 {
			break mainCicle
		}
	}

	close(ch)
	fmt.Println("Total:", total)
}
