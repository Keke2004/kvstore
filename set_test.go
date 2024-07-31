package main

import (
	"context"
	pb "kvstore/kv/kvs"
	"log"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func BenchmarkKVStore(b *testing.B) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		b.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKeyValueStoreClient(conn)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := c.Set(ctx, &pb.SetRequest{Key: "benchmark_key", Value: "benchmark_value"})
			if err != nil {
				log.Printf("could not set: %v", err)
			}
		}()
	}
	wg.Wait()
}
