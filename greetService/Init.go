/*
   Created on: 18/10/21
   Author: Pavankumar Pamuru

*/

package greetService

import "gRPCSampleApp/proto"

type Message struct {
	Messages []string `json:"messages"`
}

type Connection struct {
	Id      string
	Active  bool
	Channel chan string
}

type GreeterService struct {
	Connection []*Connection
	proto.UnimplementedGreeterServer
	Channel chan string
}
