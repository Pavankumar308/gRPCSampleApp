package main

import (
	service "gRPCSampleApp/authService"
	"gRPCSampleApp/greetService"
	"gRPCSampleApp/middleware"
	"gRPCSampleApp/proto"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"
)

func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "secret", "user")
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func accessibleRoles() map[string][]string {
	const greetServicePath = "/gRPCSampleApp.GreeterService/"

	return map[string][]string{
		greetServicePath + "SayHello":   {"admin", "user"},
		greetServicePath + "GetMessage": {"admin", "user"},
	}
}

func runGRPCServer(
	authServer proto.AuthServiceServer,
	greetServer proto.GreeterServer,
	jwtManager *service.JWTManager,
) error {
	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	}

	grpcServer := grpc.NewServer(serverOptions...)

	proto.RegisterAuthServiceServer(grpcServer, authServer)
	proto.RegisterGreeterServer(grpcServer, greetServer)
	wrappedGrpc := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(origin string) bool {
		// Allow all origins, DO NOT do this in production
		return true
	}))
	handler := middleware.NewGrpcWebMiddleware(wrappedGrpc).Handler()
	if err := http.ListenAndServe(":8989", handler); err != nil {
		panic(err)
	}
	return nil
}

func main() {
	channels := make(chan string)
	userStore := service.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users: ", err)
	}

	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)
	greetServer := &greetService.GreeterService{Channel: channels}

	err = runGRPCServer(authServer, greetServer, jwtManager)
	if err != nil {
		return
	}
}
