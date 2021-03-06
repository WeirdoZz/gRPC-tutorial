package main

import (
	"context"
	"flag"
	"fmt"
	"gRPC-tutori/pb"
	"gRPC-tutori/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

const (
	secretKey     = "weirdo"
	tokenDuration = 15 * time.Minute
)

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("--> unary interceptor:", info.FullMethod)
	return handler(ctx, req)
}

func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("-->stream interceptor:", info.FullMethod)
	return handler(srv, stream)
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "weirdo1", "123", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "weirdo2", "123", "user")
}

func main() {

	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port :%d", *port)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	userStore := service.NewInmemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users:", err)
	}
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)

	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)
	authServer := service.NewAuthServer(userStore, jwtManager)
	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	reflection.Register(grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func accessibleRoles() map[string][]string {
	const laptopServicePath = "/gRpc-tutorial.pb.LaptopService/"

	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}
