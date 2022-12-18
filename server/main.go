package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Ishankhan21/grcp-golang/server/pb"
	"google.golang.org/grpc"
)

type Myserver struct {
	pb.UnimplementedTodoServer
	Todos []*pb.TodoItem
}

func (s *Myserver) CreateTodo(ctx context.Context, in *pb.TodoItem) (*pb.TodoItem, error) {
	log.Printf("Received: %v", in.Text)
	item := &pb.TodoItem{Id: 1, Text: in.GetText()}
	t := append(s.Todos, item)
	s.Todos = t
	log.Println(s.Todos)

	return item, nil
}

func (s *Myserver) ReadTodos(ctx context.Context, noParams *pb.VoidNoParam) (*pb.TodoItems, error) {
	data := &pb.TodoItems{Items: s.Todos}
	return data, nil
}

func (s *Myserver) ReadTodosStream(noParams *pb.VoidNoParam, stream pb.Todo_ReadTodosStreamServer) error {
	for _, todo := range s.Todos {
		if err := stream.Send(todo); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	log.Println("First line")
	t := []*pb.TodoItem{
		{Id: 1, Text: "First"},
	}
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTodoServer(grpcServer, &Myserver{Todos: t})
	log.Println("2")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	log.Println("3")

}
