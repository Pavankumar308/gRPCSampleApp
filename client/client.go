/*
   Created on: 18/10/21
   Author: Pavankumar Pamuru

*/

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	username        = "admin1"
	password        = "secret"
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	const greetServicePath = "/proto.Greeter/"

	return map[string]bool{
		greetServicePath + "SayHello":   true,
		greetServicePath + "GetMessage": true,
	}
}

const (
	address = "localhost:8989"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	authClient := NewAuthClient(conn, username, password)
	interceptor, err := NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}
	temp := interceptor.Unary()
	conn2, err := grpc.Dial(
		address,
		grpc.WithUnaryInterceptor(temp),
		grpc.WithStreamInterceptor(interceptor.Stream()),
		grpc.WithInsecure(), grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	greeterClient := NewGreeter(conn2)
	r := greeterClient.SayHello("Pavan")
	fmt.Println(r.GetMessage())
}
