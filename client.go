package main

import (
	"context"
	"log"
	"sync"
	"time"

	pb "kvstore/kv/kvs"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKeyValueStoreClient(conn)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		ra, err := c.Set(ctx, &pb.SetRequest{Key: "name", Value: "grpc"})
		if err != nil {
			log.Printf("could not set: %v", err)
		}
		log.Printf("Set result: %v", ra.Success)
	}()
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.Get(ctx, &pb.GetRequest{Key: "name"})
		if err != nil {
			log.Printf("could not get: %v", err)
		}
		log.Printf("Got: %s", r.GetValue())
	}()
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err = c.Delete(ctx, &pb.DeleteRequest{Key: "name"})
		if err != nil {
			log.Printf("could not delete: %v", err)
		}
		log.Println("Deleted")
	}()
	wg.Wait()
}
