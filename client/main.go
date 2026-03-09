package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	pb "grpc-bidi-addition/grpc-bidi-addition/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewSumServiceClient(conn)

	stream, err := client.AddNumbers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// receive responses
	go func() {

		defer wg.Done()

		for {

			res, err := stream.Recv()

			if err == io.EOF {
				return
			}

			if err != nil {
				log.Fatal("Receive error:", err)
			}

			fmt.Println("Server Total:", res.Total)
		}

	}()

	numbers := []int32{10, 20, 5, 15, 56, 56, 56, 34, 2345234, 2342, 3423, 4234, 32423, 42}

	for _, num := range numbers {
		fmt.Println("Sending:", num)
		err := stream.Send(&pb.NumberRequest{
			Number: num,
		})
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)

	}

	stream.CloseSend()

	wg.Wait()
}
