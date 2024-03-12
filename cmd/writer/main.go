package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/centrifugal/gocent/v3"
)

func main() {
	ctx := context.Background()
	c := gocent.New(gocent.Config{
		Addr: "http://127.0.0.1:8000/api",
		Key:  "21a951ea-5435-4a45-b690-429afee2c4af",
	})

	_, err := c.Info(ctx)
	if err != nil {
		panic(err)
	}

	ch := "public:example"

	for i := 0; i < 10000; i++ {
		result, err := c.Publish(ctx, ch, []byte(fmt.Sprintf(`{"value": "%d"}`, i)))
		if err != nil {
			panic(err)
		}
		log.Printf("Publish into channel %s successful, stream position {offset: %d, epoch: %s}", ch, result.Offset, result.Epoch)

		time.Sleep(time.Second)
	}
}
