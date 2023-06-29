package collector

import (
	"log"
	"time"
)

type Collector interface {
	Name() string
	Collect(config CollectConfig) error
}

type CollectConfig struct {
	CollectDest    bool
	CollectTracing bool
	ClashHost      string
	ClashToken     string
}

var collectors []Collector

func Register(c Collector) {
	collectors = append(collectors, c)
}

var maxRetry = 3

func Start(config CollectConfig) {
	for _, c := range collectors {
		go func(c Collector) {
			retryCount := 0
			var err error
			for retryCount < maxRetry {
				if err = c.Collect(config); err != nil {
					retryCount++
					sleepDuration := time.Duration(retryCount) * 10 * time.Second
					log.Println("collector:", c.Name(), "failed error: ", err, "retry after", sleepDuration)
					time.Sleep(sleepDuration)
					continue
				}
				break
			}
			if err != nil {
				log.Fatalf("error start collector %s, max retry exceeded will exit", c.Name())
			}
		}(c)
	}
}
