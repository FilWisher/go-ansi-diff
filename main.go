package main

import (
	"fmt"
	"bytes"
	//"time"
//	"strings"
)

const (
	ESC_CURSOR_SET string = "\x1b[%d;%dH"
	ESC_CURSOR_UP string = "\x1b[%dA"
	ESC_CURSOR_DOWN string = "\x1b[%dB"
	ESC_CURSOR_RIGHT string = "\x1b[%dC"
	ESC_CURSOR_LEFT string = "\x1b[%dD"
	ESC_CLEAR_RIGHT string = "\x1b[K"
	ESC_CLEAR_LINE string = "\x1b[2K"
	ESC_NEWLINE string = "\n"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Differ struct {
	old [][]byte
}

func diffLine(old []byte, new []byte, changes [][]byte) [][]byte {

	min := min(len(old), len(new))
	
	i := 0
	for ; i < min; i++ {
		if old[i] != new[i] {
			break
		}
	}
	if (i == 0) {
		changes = append(changes, []byte(new))
	} else {
		changes = append(changes, []byte(fmt.Sprintf(ESC_CURSOR_RIGHT, i)))
		changes = append(changes, []byte(fmt.Sprintf(ESC_CLEAR_RIGHT)))
		changes = append(changes, []byte(fmt.Sprintf("%s", new[i:])))
		changes = append(changes, []byte(fmt.Sprintf(ESC_NEWLINE)))
	}
	
	return changes
}

func (d *Differ) Diff(new []byte) []byte {

	changes := [][]byte{}
	
	lines := bytes.Split(new, []byte("\n"))
	//diff := make([]int, len(lines))
	
	for i := len(lines); i < len(d.old); i++ {
		changes = append(changes, []byte(fmt.Sprintf(ESC_CURSOR_UP, 1)))
		changes = append(changes, []byte(fmt.Sprintf(ESC_CLEAR_LINE)))
	}

	changes = append(changes, []byte(fmt.Sprintf(ESC_CURSOR_UP, len(d.old))))

	for i := 0; i < len(lines); i++ {
		if i < len(d.old) {
			changes = diffLine(d.old[i], lines[i], changes)
		} else {
			changes = append(changes, lines[i])
		  changes = append(changes, []byte(fmt.Sprintf(ESC_NEWLINE)))
		}
		//changes = append(changes, []byte(fmt.Sprintf(ESC_CURSOR_DOWN, 1)))
	}

	d.old = lines
	return bytes.Join(changes, []byte(""))
}

func main() {

	d := Differ{}
	
	changes := d.Diff([]byte("hello\ncool!\naoeu"))
	fmt.Printf("%s", changes)
	
	changes = d.Diff([]byte("huiio\ncool!"))
	fmt.Printf("%s", changes)
	//
	//changes = d.Diff([]byte("heuuiio"))
	//fmt.Printf("%s", changes)
	
	//ticker := time.NewTicker(1 * time.Second)
	//quit := make(chan struct{})
	//go func() {
	//	for {
	//		select {
	//			case <- ticker.C:
	//				changes := d.Diff([]byte(fmt.Sprintf("%s", time.Now().String())))
	//				fmt.Printf("%s", changes)
	//			case <- quit:
	//				ticker.Stop()
	//				return
	//		}	
	//	}	
	//}()
	//
	//<-quit
}
