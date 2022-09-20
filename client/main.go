package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	pb "grpcserver/protos"

	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr     = flag.String("addr", "localhost:50051", "the address to connect to")
	message  = ""
	aciertos = 0
	total    = 0
)

func main() {
	flag.Parse()

	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          []string{"localhost:7000", "localhost:7001", "localhost:7002", "localhost:7003", "localhost:7004", "localhost:7005"},
		RouteByLatency: true,
	})

	defer redisClient.Close()

	_, err := redisClient.Ping().Result()

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("---------------------------------------------------------------\n")
		fmt.Print("Ingrese búsqueda o 'exit' para salir: ")
		scanner.Scan()
		message = scanner.Text()

		if message == "exit" {
			return
		}

		total++

		nameInRedis, err := redisClient.Get(message).Result()

		if err != nil || err == redis.Nil {
			start := time.Now()

			log.Println("No existe en redis")

			res, err := HandleGRPC()

			if err != nil {
				log.Println(err)
			}

			redisClient.Append(message, res)
			elapsed := time.Since(start)
			fmt.Printf("\nTiempo de búsqueda en base de datos: %s \n", elapsed)

		} else {
			aciertos++
			start := time.Now()

			log.Println("Encontrado en redis")

			fmt.Print(nameInRedis)
			elapsed := time.Since(start)
			fmt.Printf("\nTiempo de búsqueda en Redis: %s \n", elapsed)
		}

		fmt.Println("\nPorcentaje de aciertos: ", (aciertos*100)/total, "%")

	}

}

func HandleGRPC() (result string, err error) {

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSearchClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	objlist, err := c.GetObjects(ctx, &pb.Message{Message: message})

	if err != nil {
		log.Printf("No se detectaron coincidencias: %v", err)
		return
	}

	if objlist.GetItem() == nil {
		fmt.Printf("No se encontraron coincidencias! \n")
		result = "No se encontraron coincidencias! \n"
		return
	}

	for _, item := range objlist.GetItem() {
		fmt.Printf("{ Title: %s, Description: %s, Keywords: %s , Url: %s }\n", item.Title, item.Description, item.Keywords, item.Url)
		result += fmt.Sprintf("{ Title: %s, Description: %s, Keywords: %s , Url: %s }\n", item.Title, item.Description, item.Keywords, item.Url)
	}

	return

}
