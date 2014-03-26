// dripper project main.go
package dripper

import (
	"crypto/md5"
	"encoding/binary"
	"io"
	"math"
	"time"
)

type Dripper struct {
	Faucet   chan interface{}
	dropsMap map[string]uint16
	drops    [][]*Drop
	cursor   int
}

type Drop struct {
	Key  string
	Item interface{}
}

func NewDripper() *Dripper {
	return &Dripper{
		drops:    make([][]*Drop, 60, 60),
		dropsMap: make(map[string]uint16),
		Faucet:   make(chan interface{}, 100),
	}
}

func (d *Dripper) AddDrop(key string, item interface{}) {
	h := md5.New()
	io.WriteString(h, key)
	hash := float64(binary.LittleEndian.Uint32(h.Sum(nil)[0:4]))
	modResult := math.Mod(hash, 60)
	index := uint16(modResult)
	if d.drops[index] == nil {
		d.drops[index] = make([]*Drop, 0)
	}
	d.drops[index] = append(d.drops[index], &Drop{
		key,
		item,
	})
	d.dropsMap[key] = index
}

func (d *Dripper) Drip() {
	go func(dripper *Dripper) {
		ticks := time.Tick(time.Duration(time.Second))
		for {
			_ = <-ticks
			now := time.Now()
			offset := int(now.Second())
			for _ = dripper.cursor; dripper.cursor != offset; dripper.cursor++ {
				if d.drops[d.cursor] != nil {
					for _, drop := range dripper.drops[d.cursor] {
						d.Faucet <- drop.Item
					}
				}
				if dripper.cursor == len(dripper.drops)-1 {
					dripper.cursor = -1
				}
			}
		}
	}(d)
}
