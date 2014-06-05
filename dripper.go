// dripper project main.go
package dripper

import (
	"crypto/md5"
	"encoding/binary"
	_ "fmt"
	"io"
	"math"
	"sync"
	"time"
)

type Dripper struct {
	Faucet   chan *Drop
	dropsMap map[interface{}][][]uint16
	drops    [][]*Drop
	cursor   int
	lock     sync.RWMutex
	running  bool
	wg       sync.WaitGroup
}

type Drop struct {
	Key  string
	Item interface{}
}

func NewDripper() *Dripper {
	return &Dripper{
		drops:    make([][]*Drop, 60, 60),
		dropsMap: make(map[interface{}][][]uint16),
		Faucet:   make(chan *Drop, 100),
	}
}

func (d *Dripper) AddDrop(key string, item interface{}, interval int) {
	if _, ok := d.dropsMap[key]; ok {
		return
	}
	if interval == 0 {
		interval = 60
	}
	h := md5.New()
	io.WriteString(h, key)
	hash := float64(binary.LittleEndian.Uint32(h.Sum(nil)[0:4]))
	modResult := math.Mod(hash, float64(interval))
	for x := 0; x < 60/interval; x++ {
		index := uint16(modResult) + (uint16(x) * uint16(interval))
		if d.drops[index] == nil {
			d.drops[index] = make([]*Drop, 0)
		}
		d.drops[index] = append(d.drops[index], &Drop{
			key,
			item,
		})
		d.dropsMap[key] = append(d.dropsMap[key], make([]uint16, 2, 2))
		d.dropsMap[key][x][0] = index
		d.dropsMap[key][x][1] = uint16(len(d.drops[index]) - 1)
	}
}

func (d *Dripper) RemoveDrop(key string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	for i := range d.dropsMap[key] {
		row := d.dropsMap[key][i][0]
		column := d.dropsMap[key][i][1]
		d.drops[row] = append(d.drops[row][:column], d.drops[row][column+1:]...)
	}
	delete(d.dropsMap, key)
}

func (d *Dripper) Includes(key string) bool {
	if _, ok := d.dropsMap[key]; ok {
		return true
	} else {
		return false
	}
}

func (d *Dripper) Drip() {
	d.cursor = int(time.Now().Second())
	d.running = true
	d.wg.Add(1)
	go func(dripper *Dripper) {
		ticks := time.Tick(time.Duration(time.Second))
		for d.running == true {
			_ = <-ticks
			now := time.Now()
			offset := int(now.Second())
			d.lock.RLock()
			for _ = dripper.cursor; dripper.cursor != offset; dripper.cursor++ {
				if d.drops[d.cursor] != nil {
					for _, drop := range dripper.drops[d.cursor] {
						if drop != nil {
							d.Faucet <- drop
						}
					}
				}
				if dripper.cursor == len(dripper.drops)-1 {
					dripper.cursor = -1
				}
			}
			d.lock.RUnlock()
		}
		d.wg.Done()
	}(d)
}

func (d *Dripper) Stop() {
	d.running = false
	d.wg.Wait()
	close(d.Faucet)
}
