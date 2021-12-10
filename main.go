/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/bachelor-thesis-hown3d/chat-client/cmd"
)

func main() {

	//if clientID == "" {
	//	log.Fatal("OAUTH2_CLIENT_ID Environment Variable must be set!")
	//}
	//if clientSecret == "" {
	//	log.Fatal("OAUTH2_CLIENT_SECRET Environment Variable must be set!")
	//}

	//ctx := context.Background()
	//conn, err := grpc.DialContext(ctx, fmt.Sprintf("%v:%v", *host, *port), grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("Failed to dial %v:%v: %v", *host, *port, err)
	//}

	//rocketClient := rocketpb.NewRocketServiceClient(conn)
	//tenantClient := tenantpb.NewTenantServiceClient(conn)

	//md := metadata.Pairs("authorization", "bearer "+token)
	//ctx = metadata.NewOutgoingContext(ctx, md)

	cmd.Execute()
}
