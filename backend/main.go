package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	todov1 "backend/gen"

	todov1connect "backend/gen/todov1connect"

	connect "connectrpc.com/connect"
	"github.com/rs/cors"
)

type todoServer struct {
	mu    sync.Mutex
	tasks []todov1.Task
	id    int32
}

func (s *todoServer) AddTask(ctx context.Context, req *connect.Request[todov1.AddTaskRequest]) (*connect.Response[todov1.AddTaskResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.id++
	task := todov1.Task{Id: s.id, Text: req.Msg.Text}
	s.tasks = append(s.tasks, task)
	return connect.NewResponse(&todov1.AddTaskResponse{Task: &task}), nil
}

func (s *todoServer) GetTasks(ctx context.Context, req *connect.Request[todov1.GetTasksRequest]) (*connect.Response[todov1.GetTasksResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var tasks []*todov1.Task
	for i := range s.tasks {
		tasks = append(tasks, &s.tasks[i])
	}

	return connect.NewResponse(&todov1.GetTasksResponse{Tasks: tasks}), nil
}

func main() {
	mux := http.NewServeMux()
	server := &todoServer{}

	path, handler := todov1connect.NewTodoServiceHandler(server)
	mux.Handle(path, handler)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: newCORS().Handler(mux),
	}

	log.Println("Listening on :8000")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func newCORS() *cors.Cors {
	// To let web developers play with the demo service from browsers, we need a
	// very permissive CORS setup.
	return cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(_ /* origin */ string) bool {
			// Allow all origins, which effectively disables CORS.
			return true
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			// Content-Type is in the default safelist.
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
		},
		// Let browsers cache CORS information for longer, which reduces the number
		// of preflight requests. Any changes to ExposedHeaders won't take effect
		// until the cached data expires. FF caps this value at 24h, and modern
		// Chrome caps it at 2h.
		MaxAge: int(2 * time.Hour / time.Second),
	})
}
