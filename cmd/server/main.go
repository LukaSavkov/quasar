package main

import (
	"fmt"
	"log"
	"net"
	"os"

	oortapi "github.com/c12s/oort/pkg/api"
	"github.com/jtomic1/config-schema-service/internal/configschema"
	"github.com/jtomic1/config-schema-service/internal/services"
	pb "github.com/jtomic1/config-schema-service/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(configschema.GetAuthInterceptor()))

	administrator, err := oortapi.NewAdministrationAsyncClient(os.Getenv("NATS_ADDRESS"))
	if err != nil {
		log.Fatalln(err)
	}
	authorizer := services.NewAuthZService(os.Getenv("SECRET_KEY"))
	configSchemaServer := configschema.NewServer(authorizer, administrator)

	pb.RegisterConfigSchemaServiceServer(grpcServer, configSchemaServer)

	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
