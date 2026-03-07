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

	"database/sql"
	"os"
	_ "github.com/lib/pq"

)

var (
    host     = os.Getenv("DB_ADDRESS")
    port     = 5432
    user     = "postgres"
    password = "password"
    dbname   = "postgres"
)




func Open() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println(psqlInfo)
    return sql.Open("postgres", psqlInfo)
}



type Message struct {
	ID        string
	Author    string
	Content   string
	Timestamp time.Time
}

// in memory server for now
type ForumServer struct {
	pb.UnimplementedForumServiceServer
	mu       sync.RWMutex
}

func (s *ForumServer) PostMessage(ctx context.Context, req *pb.PostRequest) (*pb.PostResponse, error) {
	if req.Author == "" || req.Message == "" {
		return &pb.PostResponse{Success: false}, fmt.Errorf("author and message required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	msg := &Message{
		ID:        "1",// gen random
		Author:    req.Author,
		Content:   req.Message,
		Timestamp: time.Now(),
	}

	db,err := Open()
	if err!=nil{
		return nil, err
	}
	res,err:=db.Exec("INSERT INTO forummessage (id, author, message, timestamp) VALUES ($1, $2, $3, $4)", msg.ID, msg.Author, msg.Content, msg.Timestamp)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(res)


	return &pb.PostResponse{
		Id:        msg.ID,
		Author:    msg.Author,
		Message:   msg.Content,
		Timestamp: msg.Timestamp.Format("24.12.2001 14:45"),
		Success:   true,
	}, nil
}

func (s *ForumServer) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.MessagesResponse, error) {
	fmt.Println("GET MESSAGE ENDPOINT")
	s.mu.RLock()
	defer s.mu.RUnlock()


	db,err := Open()
	if err!=nil{
		return nil, err
	}
	rows, err:=db.Query("SELECT * FROM forummessage;")
	if err!=nil{
		fmt.Println("query error",err)
		return nil, err
	}
	var (
		_id 		string
		_time 		time.Time
		author  	string
		content   	string
	)

	var rv = make([]*pb.Message,0)
	for rows.Next() {
		fmt.Println("rown")
		err := rows.Scan(&_id, &author, &content, &_time)
		msg := &pb.Message{
			Id:        _id,// gen random
			Author:    author,
			Message:   content,
			Timestamp: _time.Format("02.01.2006 15:04"),
		}
		rv = append(rv, msg)
		
		if err!=nil{
			fmt.Println("error on line ",err)
		}
	}

	return &pb.MessagesResponse{Messages: rv}, nil
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
	pb.RegisterForumServiceServer(server, &ForumServer{ })

	log.Println("gRPC server listening on :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
