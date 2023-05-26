package main

import (
	pb "github.com/jel221/FanTheoryDB/src/pb"
	"github.com/jel221/FanTheoryDB/src/driver.go"
	"context"
	"log"
	"net"
	"time"
    "os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type server struct {
	pb.UnimplementedTheoryDBServer
	coll *mongo.Collection
}

func (s *server) GetTheory(ctx context.Context, in *pb.GetTheoryRequest) (*pb.GetTheoryReply, error) {
	log.Printf("Received request: %v", in.ProtoReflect().Descriptor().FullName())
	return &pb.GetTheoryReply{
	}, nil
}

func (s *server) PutTheory(ctx context.Context, in *pb.PutTheoryRequest) (*pb.PutTheoryReply, error) {
	log.Printf("Received request: %v", in.ProtoReflect().Descriptor().FullName())
	
	var err error
	if err == nil {
		return &pb.PutTheoryReply{
			Success: true,
			Error: "",
		}, nil
	} else {
		return &pb.PutTheoryReply{
			Success: false,
			Error: err.Error(),
		}, nil
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	/* Establish connection to DB server */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cred, err := os.ReadFile("/tmp/dat")
	if err != nil {
		panic(err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(string(cred)))
	if err != nil {
		panic(err)
	}

	collection := client.Database("TestDB").Collection("Theory")

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterTheoryDBServer(s, &server{coll: collection})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}