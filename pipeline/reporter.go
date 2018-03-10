package pipeline

import (
	"fmt"
	"time"
)

func RunReporter(pipeline *Pipeline) {
	pipeline.Run()
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for _ = range ticker.C {
			fmt.Println(pipeline.Stats())
		}
	}()

	pipeline.Wait()
	ticker.Stop()
	fmt.Println(pipeline.Stats())
}
