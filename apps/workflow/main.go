package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/luis12loureiro/neurun/apps/workflow/gen"
	ws "github.com/luis12loureiro/neurun/apps/workflow/internal/workflow"
	wh "github.com/luis12loureiro/neurun/apps/workflow/internal/workflow/handler"
	wr "github.com/luis12loureiro/neurun/apps/workflow/internal/workflow/repository"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	// repo, err := wr.NewSQLiteRepository("./internal/workflow/repository/storage", "data.sqlite")
	// if err != nil {
	// 	log.Fatalf("failed to create repository: %v", err)
	// }
	// defer func() {
	// 	if sqliteRepo, ok := repo.(*wr.SQLiteRepo); ok {
	// 		sqliteRepo.Close()
	// 	}
	// }()
	repo := wr.NewMemoryRepository()
	te := ws.NewTaskExecutor()
	we := ws.NewWorkflowExecutor(repo, te)
	svc := ws.NewService(repo, we)
	handler := wh.NewServer(svc)
	pb.RegisterWorkflowServiceServer(s, handler)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
