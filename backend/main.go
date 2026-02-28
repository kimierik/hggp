package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	pb "github.com/kimierik/hggp/backend/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Message struct {
	ID        string
	Author    string
	Content   string
	Timestamp string
}

// in memory server for now
type ForumServer struct {
	pb.UnimplementedForumServiceServer
	messages []*Message
	mu       sync.RWMutex
}

func (s *ForumServer) PostMessage(ctx context.Context, req *pb.PostRequest) (*pb.PostResponse, error) {
	if req.Author == "" || req.Message == "" {
		return &pb.PostResponse{Success: false}, fmt.Errorf("author and message required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	msg := &Message{
		ID:        "1",
		Author:    req.Author,
		Content:   req.Message,
		Timestamp: time.Now().Format("03:04 PM"),
	}

	s.messages = append([]*Message{msg}, s.messages...)

	return &pb.PostResponse{
		Id:        msg.ID,
		Author:    msg.Author,
		Message:   msg.Content,
		Timestamp: msg.Timestamp,
		Success:   true,
	}, nil
}

func (s *ForumServer) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.MessagesResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	limit := int(req.Limit)
	if limit <= 0 || limit > len(s.messages) {
		limit = len(s.messages)
	}

	pbMessages := make([]*pb.Message, limit)
	for i := 0; i < limit; i++ {
		msg := s.messages[i]
		pbMessages[i] = &pb.Message{
			Id:        msg.ID,
			Author:    msg.Author,
			Message:   msg.Content,
			Timestamp: msg.Timestamp,
		}
	}

	return &pb.MessagesResponse{Messages: pbMessages}, nil
}

func startGateway(){
	ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := runtime.NewServeMux()
	insecure.NewCredentials()
    opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	err:= pb.RegisterForumServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
    if err != nil {
        log.Fatalf("failed to start HTTP gateway: %v", err)
    }
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // i love cors
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "86400")
    	w.Header().Set("Vary", "Origin")  // ← Add this line

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        mux.ServeHTTP(w, r)
    })


    log.Println("HTTP server running on :8000")
    http.ListenAndServe(":8000", handler)
}

func main() {
	go startGateway()


	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterForumServiceServer(server, &ForumServer{
		messages: make([]*Message, 0),
	})

	log.Println("gRPC server listening on :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
