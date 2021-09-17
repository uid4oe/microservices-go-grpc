package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/uid4oe/microservices-go-grpc/advice/advicedb"
	"github.com/uid4oe/microservices-go-grpc/advice/advicepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	timeout         = time.Second
	postgres_client *pgxpool.Pool
)

type server struct {
	advicepb.UnimplementedAdviceServiceServer
}

func (*server) CreateUpdateAdvice(ctx context.Context, req *advicepb.CreateUpdateAdviceRequest) (*advicepb.CreateUpdateAdviceResponse, error) {
	log.Println("Called CreateUpdateAdvice, Operation", req.Operation)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var err error
	if req.Operation == advicepb.Operation_CREATE {
		err = advicedb.CreateOne(postgres_client, c, &advicedb.Advice{
			UserId: req.UserId,
			Advice: req.Advice})
	} else {
		err = advicedb.UpdateOne(postgres_client, c, &advicedb.Advice{
			UserId: req.UserId,
			Advice: req.Advice})
	}
	if err != nil {
		return nil, error_response(err)
	}

	return &advicepb.CreateUpdateAdviceResponse{}, nil
}

func (*server) GetAdvice(ctx context.Context, req *advicepb.GetUserAdviceRequest) (*advicepb.GetUserAdviceResponse, error) {
	log.Println("Called GetAdvice for User Id", req.UserId)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result, err := advicedb.FindOne(postgres_client, c, req.UserId)
	if err != nil {
		return nil, error_response(err)
	}

	return &advicepb.GetUserAdviceResponse{Advice: result.Advice, CreatedAt: timestamppb.New(result.CreatedAt)}, nil
}

func error_response(err error) error {
	log.Println("ERROR:", err.Error())
	return status.Error(codes.Internal, err.Error())
}

func main() {
	log.Println("Advice Service")

	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Println("ERROR:", err.Error())
	}

	postgres_client, err = advicedb.NewClient(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer postgres_client.Close()

	s := grpc.NewServer()
	advicepb.RegisterAdviceServiceServer(s, &server{})

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}

}
