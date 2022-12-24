package main

import (
	"context"
	"io"
	"log"

	pb "github.com/Ishankhan21/grcp-golang/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to listen client: %v", err)
	}

	client := pb.NewTodoClient(conn)

	// Create TODO
	todoItem, err := client.CreateTodo(context.Background(), &pb.TodoItem{Id: 1, Text: "Created new Todo ITEM"})
	if err != nil {
		log.Fatalf("failed to create todo", err)
	}

	log.Printf("Created Todo: %v", todoItem)

	// Read all todos
	todoItems, err := client.ReadTodos(context.Background(), &pb.VoidNoParam{})
	log.Printf("All TODOS: ", todoItems)

	// Read todos in stream
	todoStream, err := client.ReadTodosStream(context.Background(), &pb.VoidNoParam{})
	done := make(chan bool)
	go func() {
		for {
			resp, err := todoStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}
			log.Printf("Stream response received:", resp)
		}
	}()
	<-done //we will wait until all response is received

	log.Printf("finished")

}
