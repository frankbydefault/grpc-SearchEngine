/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
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

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedSearchServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

func (s *server) GetObjects(ctx context.Context, in *pb.Message) (*pb.SearchResponse, error) {

	str := fmt.Sprintf("%%%s%%", in.Message)

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

	//defer db.Close()
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
