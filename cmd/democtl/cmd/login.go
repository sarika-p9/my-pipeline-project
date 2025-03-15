package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/authentication"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login a user",
	Run: func(cmd *cobra.Command, args []string) {
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Fatalf("Error reading email flag: %v", err)
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatalf("Error reading password flag: %v", err)
		}
		if email == "" || password == "" {
			log.Fatal("❌ Email and password are required.")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err := grpc.DialContext(ctx, "localhost:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			log.Fatalf("❌ Failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()
		client := proto.NewAuthServiceClient(conn)
		reqCtx, reqCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer reqCancel()
		fmt.Println("🔐 Attempting login...")

		resp, err := client.Login(reqCtx, &proto.LoginRequest{Email: email, Password: password})
		if err != nil {
			log.Fatalf("❌ Login failed: %v", err)
		}
		fmt.Printf("\n✅ Login successful!\nUserID: %s\nEmail: %s\nToken: %s\n", resp.UserId, resp.Email, resp.Token)
	},
}

func init() {
	loginCmd.Flags().String("email", "", "User Email")
	loginCmd.Flags().String("password", "", "User Password")
	loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagRequired("password")
}
