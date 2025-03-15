package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the main API server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting main_api_server...")
		err := godotenv.Load("../main_server/.env")
		if err := godotenv.Load(".env"); err != nil {
			log.Println("Warning: Could not load .env file. Using system environment variables.")
		}
		mainServerPath, err := filepath.Abs("../main_server/main_api_server")
		if err != nil {
			fmt.Println("Error resolving main_api_server path:", err)
			os.Exit(1)
		}
		if _, err := os.Stat(mainServerPath); os.IsNotExist(err) {
			fmt.Println("Error: main_api_server does not exist at", mainServerPath)
			os.Exit(1)
		}
		postgresDSN := os.Getenv("POSTGRES_DSN")
		if postgresDSN == "" {
			fmt.Println("Error: POSTGRES_DSN environment variable is not set")
			os.Exit(1)
		}
		serverCmd := exec.Command(mainServerPath)
		serverCmd.Env = append(os.Environ(), "POSTGRES_DSN="+postgresDSN)
		serverCmd.Stdout = os.Stdout
		serverCmd.Stderr = os.Stderr

		err = serverCmd.Start()
		if err != nil {
			fmt.Println("Error starting main_api_server:", err)
			os.Exit(1)
		}

		fmt.Println("main_api_server started successfully!")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
