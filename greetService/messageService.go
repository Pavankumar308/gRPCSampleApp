/*
   Created on: 18/10/21
   Author: Pavankumar Pamuru

*/

package greetService

import (
	"fmt"
	helloworld "gRPCSampleApp/proto"
)

func (s *GreeterService) GetMessage(req *helloworld.MessageRequest, stream helloworld.Greeter_GetMessageServer) error {
	channel := s.Channel
	defer close(channel)
	for message := range channel {
		fmt.Println("Message: ", message)
		err := stream.Send(&helloworld.MessageReply{
			ReplyMsg: message,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
