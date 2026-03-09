package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "grpc-bidi-addition/grpc-bidi-addition/pkg/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSumServiceServer
}

func (s *server) AddNumbers(stream pb.SumService_AddNumbersServer) error {

	var total int32 = 0

	for {

		req, err := stream.Recv()

		if err == io.EOF {
			fmt.Println("Client finished sending")
			return nil
		}

		if err != nil {
			return err
		}

		number := req.GetNumber()
		total += number

		fmt.Println("Received:", number, "Current Total:", total)

		err = stream.Send(&pb.NumberResponse{
			Total: total,
		})

		if err != nil {
			return err
		}
	}
}

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterSumServiceServer(grpcServer, &server{})

	fmt.Println("Bidirectional gRPC server running on port 50051...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
