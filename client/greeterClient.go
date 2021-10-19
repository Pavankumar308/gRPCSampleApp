/*
   Created on: 18/10/21
   Author: Pavankumar Pamuru

*/

package main

import (
	"context"
	"gRPCSampleApp/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LaptopClient is a client to call laptop service RPCs
type GreeterClient struct {
	service proto.GreeterClient
}

// NewLaptopClient returns a new laptop client
func NewGreeter(cc *grpc.ClientConn) *GreeterClient {
	service := proto.NewGreeterClient(cc)
	return &GreeterClient{service}
}

// CreateLaptop calls create laptop RPC
func (greeterClient *GreeterClient) SayHello(name string) *proto.HelloReply {
	req := &proto.HelloRequest{
		Name: name,
	}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := greeterClient.service.SayHello(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			// not a big deal
			log.Print("laptop already exists")
		} else {
			log.Fatal("cannot create laptop: ", err)
		}
		return &proto.HelloReply{}
	}

	log.Printf("created laptop with id: %s", res.GetMessage())
	return res
}
