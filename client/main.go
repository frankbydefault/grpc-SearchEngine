package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "grpcserver/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr    = flag.String("addr", "localhost:50051", "the address to connect to")
	message = ""
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSearchClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for {
		fmt.Print("---------------------------------------------------------------\n")
		fmt.Print("Engrese busqueda o 'exit' para salir: ")
		fmt.Scanf("%s", &message)

		if message == "exit" {
			return
		}

		objlist, err := c.GetObjects(ctx, &pb.Message{Message: *&message})
		if err != nil {
			log.Printf("No se detectaron coincidencias: %v", err)
		}

		if objlist.GetItem() == nil {
			fmt.Printf("No se encontraron coincidencias! \n")
		}

		for _, item := range objlist.GetItem() {
			fmt.Printf("{ Title: %s, Description: %s, Keywords: %s , Url: %s }\n", item.Title, item.Description, item.Keywords, item.Url)
		}

	}

}
