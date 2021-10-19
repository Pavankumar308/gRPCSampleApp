/*
   Created on: 18/10/21
   Author: Pavankumar Pamuru

*/

package greetService

import (
	"context"
	"encoding/json"
	"fmt"
	helloworld "gRPCSampleApp/proto"
	"io/ioutil"
)

// server is used to implement helloworld.GreeterServer.

// SayHello implements helloworld.GreeterServer

func (s *GreeterService) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	fmt.Println("\n\n\n############\n\nEntered To SayHello: ", in.GetName())
	file, _ := ioutil.ReadFile("messages.json")

	data := Message{}

	_ = json.Unmarshal([]byte(file), &data)
	data.Messages = append(data.Messages, in.GetName())
	file1, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile("messages.json", file1, 0644)
	fmt.Println("#####\n\nBefore Pushed to Channel")
	//s.Channel <- in.GetName()
	fmt.Println("#####\n\nPushed to Channel")
	return &helloworld.HelloReply{Message: in.GetName()}, nil
}
