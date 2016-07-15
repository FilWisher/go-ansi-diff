package differ

import (
	"time"
	"fmt"
)

func ExampleDiffer() {

	d := Differ{}
	
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
				case <- ticker.C:
					fmt.Printf("%s", d.Diff([]byte(fmt.Sprintf("%s", time.Now().String()))))
				case <- quit:
					ticker.Stop()
					return
			}	
		}	
	}()
	
	<-quit
}
