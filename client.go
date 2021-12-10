package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	rocketpb "github.com/bachelor-thesis-hown3d/chat-api-server/proto/rocket/v1"
	tenantpb "github.com/bachelor-thesis-hown3d/chat-api-server/proto/tenant/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	name      string = "test-rocket"
	namespace string = "testuser"
	user      string = "TestUser"
	email     string = "testUser@foo.bar"
)

var (
	port        = flag.Int("port", 10000, "Port of the api server")
	host        = flag.String("host", "", "Hostname of the api server")
	redirectURI = flag.String("redirectUrl", "http://localhost:7070", "address for the oauth server to listen on")
	issuerURI   = flag.String("issuerUrl", "https://keycloak:8443/auth/realms/kubernetes", "address for the oauth server to listen on")
)

func defaultFlow(ctx context.Context, rocketClient rocketpb.RocketServiceClient, tenantClient tenantpb.TenantServiceClient) {

	defer rocketClient.Delete(ctx, &rocketpb.DeleteRequest{Name: name, Namespace: namespace})

	versions, err := rocketClient.AvailableVersions(ctx, &rocketpb.AvailableVersionsRequest{Image: rocketpb.AvailableVersionsRequest_IMAGE_ROCKETCHAT})
	if err != nil {
		log.Fatal(fmt.Errorf("Can't get available version of rocketchat: %w", err))
	}
	for i := 0; i < 5; i++ {
		fmt.Printf("Rocket Tag: %v\n", versions.Tags[i])
	}
	versions, err = rocketClient.AvailableVersions(ctx, &rocketpb.AvailableVersionsRequest{Image: rocketpb.AvailableVersionsRequest_IMAGE_MONGODB})
	if err != nil {
		log.Fatal(fmt.Errorf("Can't get available version of mongodb: %w", err))
	}
	for i := 0; i < 5; i++ {
		fmt.Printf("Mongodb Tag: %v\n", versions.Tags[i])
	}
	//

	_, err = tenantClient.Register(ctx, &tenantpb.RegisterRequest{Size: tenantpb.RegisterRequest_SIZE_SMALL})
	if err != nil {
		log.Fatal(err)
	}

	allRockets, err := rocketClient.GetAll(ctx, &rocketpb.GetAllRequest{
		Namespace: namespace,
	})
	if err != nil {
		log.Fatalf("Can't get all rockets: %v", err)
	}
	for index, rocket := range allRockets.Rockets {
		fmt.Printf("%v: %v - %v\n", index, rocket.Name, rocket.Namespace)
	}

	_, err = rocketClient.Create(ctx, &rocketpb.CreateRequest{
		Name:         name,
		Namespace:    namespace,
		DatabaseSize: 10,
		Replicas:     1,
		Email:        email,
		Host:         "test.chat-cluster.com",
	})

	if err != nil {
		log.Fatalf("Error creating new rocket: %v", err)
	}

	newRocket, err := rocketClient.Get(ctx, &rocketpb.GetRequest{
		Name:      name,
		Namespace: namespace,
	})
	if err != nil {
		log.Fatalf("Can't get rocket: %v", err)
	}
	//
	// watch the rocket to get ready
	statusClient, err := rocketClient.Status(ctx, &rocketpb.StatusRequest{Name: newRocket.Name, Namespace: newRocket.Namespace})
	if err != nil {
		log.Fatalf("Error watching new rocket: %v", err)
	}
	var ready bool
	for ready == false {
		msg, err := statusClient.Recv()
		if status.Code(err) == codes.Canceled {
			log.Println("Context was canceled")
			break
		}
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Fatalf("Error: %v", err.Error())
		}
		fmt.Printf("StatusMessage: %v - Ready: %v\n", msg.Status, msg.Ready)
		ready = msg.Ready
	}
	//
	// logsClient, err := rocketClient.Logs(ctx, &rocketpb.LogsRequest{Name: newRocket.Name, Namespace: newRocket.Namespace})
	// for {
	// 	// blocking
	// 	msg, err := logsClient.Recv()
	// 	if status.Code(err) == codes.Canceled {
	// 		log.Println("Context was canceled")
	// 		os.Exit(0)
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("Error: %v", err.Error())
	// 	}
	// 	fmt.Printf("Pod: %v - Msg: %v\n", msg.Pod, msg.Message)
	// }
}
