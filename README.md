//串行
package main

import (
 "context"
 "log"
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
 ctx, cancel := context.WithTimeout(context.Background(), time.Second)
 defer cancel()
 ra, err := c.Set(ctx, &pb.SetRequest{Key: "name", Value: "grpc"})
 if err != nil {
  log.Fatalf("could not set: %v", err)
 }
 log.Printf("Set result: %v", ra.Success)
 r, err := c.Get(ctx, &pb.GetRequest{Key: "name"})
 if err != nil {
  log.Fatalf("could not get: %v", err)
 }
 log.Printf("Got: %s", r.GetValue())
 _, err = c.Delete(ctx, &pb.DeleteRequest{Key: "name"})
 if err != nil {
  log.Fatalf("could not delete: %v", err)
 }
 log.Println("Deleted")
}
//并行1
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
//并行2
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"

 pb "kvstore/kv/kvs"

 "google.golang.org/grpc"
)

type OperationResult struct {
 Op    string
 Error error
 Data  interface{}
}

func main() {
 conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
 if err != nil {
  log.Fatalf("did not connect: %v", err)
 }
 defer conn.Close()
 c := pb.NewKeyValueStoreClient(conn)
 results := make(chan OperationResult, 3)
 var wg sync.WaitGroup
 wg.Add(1)
 go func() {
  defer wg.Done()
  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()
  ra, err := c.Set(ctx, &pb.SetRequest{Key: "name", Value: "grpc"})
  results <- OperationResult{Op: "Set", Error: err, Data: ra}
 }()
 wg.Add(1)
 go func() {
  defer wg.Done()
  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()
  r, err := c.Get(ctx, &pb.GetRequest{Key: "name"})
  results <- OperationResult{Op: "Get", Error: err, Data: r}
 }()
 wg.Add(1)
 go func() {
  defer wg.Done()
  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()
  _, err := c.Delete(ctx, &pb.DeleteRequest{Key: "name"})
  results <- OperationResult{Op: "Delete", Error: err}
 }()
 wg.Wait()
 close(results)
 for result := range results {
  if result.Error != nil {
   log.Printf("Failed %s: %v", result.Op, result.Error)
  } else {
   switch result.Op {
   case "Set":
    ra := result.Data.(*pb.SetResponse)
    fmt.Printf("Set result: %v\n", ra.Success)
   case "Get":
    r := result.Data.(*pb.GetResponse)
    fmt.Printf("Got: %s\n", r.GetValue())
   case "Delete":
    fmt.Println("Deleted")
   }
  }
 }
}
