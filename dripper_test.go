package dripper

import "testing"
import "fmt"
import "time"
import "io/ioutil"
import "bytes"

func TestCreate(t *testing.T) {
	d := NewDripper()
	number := 0
	urlsBytes, err := ioutil.ReadFile("urls.txt")
	_ = err
	urlsBytes = urlsBytes[:len(urlsBytes)-1]
	urls := bytes.Split(urlsBytes, []byte("\n"))
	for _, url := range urls {
		d.AddDrop(string(url), number)
	}
	d.Drip()
	go printDrops(d.faucet)
	time.Sleep(time.Second * 500)
}

func printDrops(faucet chan interface{}) {
	for {
		key := <-faucet
		fmt.Println(key.(string))
	}
}
