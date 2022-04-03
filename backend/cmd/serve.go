/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/seth2810/go-collab/internal/config"
	"github.com/seth2810/go-collab/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Read()
		if err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}

		return server.ServeContext(cmd.Context(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
