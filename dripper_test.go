package dripper

import "testing"
import "fmt"
import "time"
import "io/ioutil"
import "bytes"
import "strconv"

func TestAdd(t *testing.T) {
	d := NewDripper()
	d.AddDrop("Test", 0, 0)
	if !d.Includes("Test") {
		t.Error()
	}
}

func TestRemove(t *testing.T) {
	d := NewDripper()
	d.AddDrop("Test", 0, 0)
	d.RemoveDrop("Test")
	if d.Includes("Test") {
		t.Error()
	}
}

func TestRemoveWhileDripping(t *testing.T) {
	d := NewDripper()
	for x := 0; x < 10; x++ {
		d.AddDrop(strconv.Itoa(x), x, 5)
	}
	d.Drip()
	go printDrops(d.Faucet)
	time.Sleep(6 * time.Second)
	d.RemoveDrop("5")
	fmt.Println("5 Removed")
	time.Sleep(10 * time.Second)
	d.Stop()
}

func TestDrip(t *testing.T) {
	d := NewDripper()
	number := 0
	urlsBytes, err := ioutil.ReadFile("urls.txt")
	_ = err
	urlsBytes = urlsBytes[:len(urlsBytes)]
	urls := bytes.Split(urlsBytes, []byte("\n"))
	for _, url := range urls {
		d.AddDrop(string(url), number, 10)
	}
	d.Drip()
	go printDrops(d.Faucet)
	time.Sleep(time.Second * 10)
	d.Stop()
}

func printDrops(faucet chan *Drop) {
	for {
		drop, ok := <-faucet
		if !ok {
			return
		}
		fmt.Println(drop.Key)
	}
}
