package sample

import (
	"context"
	"log"
	"time"
)

func DummyActivity(ctx context.Context) (string, error) {

	log.Println("sleeping for 10 seconds")

	time.Sleep(time.Second * 10)

	return "success in activity", nil
}
