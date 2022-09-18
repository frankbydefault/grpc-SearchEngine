package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	pb "grpcserver/protos"
	"log"
	"net"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

var db *sql.DB
var err error

type server struct {
	pb.UnimplementedSearchServer
}

func (s *server) GetObjects(ctx context.Context, in *pb.Message) (*pb.SearchResponse, error) {

	str := fmt.Sprintf("%%%s%%", in.Message)

	log.Println("Received: ", in.Message, "Search: ", str)

	stmtOut, err := db.Prepare("SELECT d.id, d.title, d.description, d.url, k.keyword FROM data AS d, keywords AS k WHERE d.title LIKE ? OR d.description LIKE ? OR (k.keyword LIKE ? AND k.id_data = d.id);")
	if err != nil {
		panic(err.Error())
	}

	defer stmtOut.Close()

	var res []*pb.Item
	rows, err := stmtOut.Query(str, str, str)
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	i := 0

	for rows.Next() && i < 10 {
		var id int
		var title string
		var description sql.NullString
		var url string
		var keyword string
		err = rows.Scan(&id, &title, &description, &url, &keyword)
		if err != nil {
			panic(err.Error())
		}
		res = append(res, &pb.Item{Id: int32(id), Title: title, Description: description.String, Url: url, Keywords: keyword})
	}

	return &pb.SearchResponse{Item: res}, nil
}

func main() {
	flag.Parse()

	db, err = sql.Open("mysql", "searchengine:S34rch3r_3ng1n3@tcp(localhost:3306)/search_engine")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSearchServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
