package main

import (
	"context"
	pb "kvstore/kv/kvs"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"sync"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedKeyValueStoreServer
	store map[string]string
	mu    sync.Mutex
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetResponse, error) {
	s.mu.Lock()
	s.store[in.Key] = in.Value
	s.mu.Unlock()
	return &pb.SetResponse{Success: true}, nil
}

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	s.mu.Lock()
	value, err := s.store[in.Key]
	s.mu.Unlock()
	if !err {
		return &pb.GetResponse{Value: ""}, nil
	}
	return &pb.GetResponse{Value: value}, nil
}

func (s *server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	delete(s.store, in.Key)
	s.mu.Unlock()
	return &pb.DeleteResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	kvServer := &server{store: make(map[string]string)}
	pb.RegisterKeyValueStoreServer(s, kvServer)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(os.Stdout)
	runtime.GOMAXPROCS(1)
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
