package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
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

	// Wrap gRPC server with gRPC-Web
	grpcWebServer := grpcweb.WrapServer(s,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true // Allow all origins for development
		}),
		grpcweb.WithWebsockets(true),
	)

	// Create HTTP server that handles both gRPC-Web and gRPC
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", *port),
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			// Handle CORS preflight
			if req.Method == http.MethodOptions {
				resp.Header().Set("Access-Control-Allow-Origin", "*")
				resp.Header().Set("Access-Control-Allow-Headers", "*")
				resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
				resp.WriteHeader(http.StatusOK)
				return
			}

			// Check if it's a gRPC-Web request
			if grpcWebServer.IsGrpcWebRequest(req) || grpcWebServer.IsAcceptableGrpcCorsRequest(req) {
				resp.Header().Set("Access-Control-Allow-Origin", "*")
				grpcWebServer.ServeHTTP(resp, req)
				return
			}

			// If not recognized
			http.NotFound(resp, req)
		}),
	}

	log.Printf("gRPC-Web server listening at %v", lis.Addr())
	if err := httpServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
