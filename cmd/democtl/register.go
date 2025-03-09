package main

import (
	"context"
	"fmt"
	"log"
	"time"

	proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/authentication"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		client := proto.NewAuthServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.Register(ctx, &proto.RegisterRequest{Email: email, Password: password})
		if err != nil {
			log.Fatalf("Registration failed: %v", err)
		}

		fmt.Printf("User registered successfully! UserID: %s, Email: %s\n", resp.UserId, resp.Email)
	},
}

func init() {
	registerCmd.Flags().String("email", "", "User Email")
	registerCmd.Flags().String("password", "", "User Password")
	registerCmd.MarkFlagRequired("email")
	registerCmd.MarkFlagRequired("password")
}
