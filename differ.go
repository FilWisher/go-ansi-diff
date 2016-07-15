package differ

import (
	"fmt"
	"bytes"
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

	bufs := [][]byte{}
	
	lines := bytes.Split(new, []byte("\n"))
	diff := make([]int, len(lines))
	pos := -1
	
	for i := 0; i < len(lines); i++ {
		if i < len(d.old) {
			old := d.old[i]
			if string(old) != string(lines[i]) {
				if pos == -1 {
					pos = i	
				}
				diff[i] = 1
			} else {
				diff[i] = 0		
			}
		}	else {
			diff[i] = 2	
		}
	}
	
	movedOnce := false
	
	for i := len(lines); i < len(d.old); i++ {
		if !movedOnce {
			bufs = append(bufs, []byte(fmt.Sprintf(ESC_CURSOR_LEFT, 1)))
			movedOnce = true
		}
		bufs = append(bufs, []byte(fmt.Sprintf(ESC_CURSOR_UP, 1)))
		bufs = append(bufs, []byte(fmt.Sprintf(ESC_CLEAR_LINE)))
	}

	if pos > -1 {
		missing := min(len(d.old), len(lines)) - pos
		bufs = append(bufs, []byte(fmt.Sprintf(ESC_CURSOR_UP, missing)))
	} else {
		pos = len(d.old)	
	}

	for ; pos < len(lines); pos++ {
		if !movedOnce {
			bufs = append(bufs, []byte(fmt.Sprintf(ESC_CURSOR_LEFT, 1)))
			movedOnce = true
		}	
		if diff[pos] > 0 {
			if diff[pos] == -1 {
				bufs = diffLine(d.old[pos], lines[pos], bufs)
			}	else {
				bufs = append(bufs, []byte(fmt.Sprintf("%s\n", lines[pos])))
			}
		} else {
			bufs = append(bufs, []byte(fmt.Sprintf(ESC_CURSOR_DOWN, 1)))
		}
	}

	d.old = lines
	return bytes.Join(bufs, []byte(""))
}
