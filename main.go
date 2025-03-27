package main

import (
	"github.com/not-for-prod/implgen/internal/basic"
	"github.com/not-for-prod/implgen/internal/repo"
)

func main() {
	rootCmd := basic.NewCMD()
	rootCmd.AddCommand(repo.NewCMD())

	_ = rootCmd.Execute()
}
