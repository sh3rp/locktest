package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/etcd-io/etcd/etcdserver/api/v3lock/v3lockpb"
	"google.golang.org/grpc"
)

var lockName = "testlock"

func main() {
	etcdString := "localhost:23791"

	conn, err := grpc.Dial(etcdString, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("no connect: %v", err)
	}

	defer conn.Close()

	client := v3lockpb.NewLockClient(conn)
	ctx := context.Background()

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			fmt.Printf("[%d] Locking\n", idx)

			lockResponse, err := client.Lock(ctx, &v3lockpb.LockRequest{
				Name: []byte(lockName),
			})
			key := lockResponse.Key
			fmt.Printf("[%d] Lock: %+v, error = %v\n", idx, lockResponse, err)
			waitTime := time.Second * time.Duration(rand.Intn(5))
			fmt.Printf("[%d] Doing stuff for %v seconds\n", idx, waitTime)
			time.Sleep(waitTime)
			fmt.Printf("[%d] Unlocking\n", idx)
			unlockResponse, err := client.Unlock(ctx, &v3lockpb.UnlockRequest{
				Key: key,
			})
			fmt.Printf("[%d] Unlock: %+v, error = %v", idx, unlockResponse, err)
			wg.Done()
		}(i)
	}

	wg.Wait()

}
