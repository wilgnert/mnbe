package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wilgnert/mnbe/internal/pkg/api"
	"github.com/wilgnert/mnbe/internal/pkg/lib"
	"github.com/wilgnert/mnbe/internal/pkg/routes"
	"github.com/wilgnert/mnbe/internal/pkg/server"
)

func main() {
	db := lib.StubDatabase{}
	env_config := readConfigFile()
	cfg := api.Config{}
	cfg.SetDB(&db)
	cfg.RegisterHandlers()
	cfg.SetSecret(env_config["secret"])

	router := routes.LoadRoutes(&cfg)
	s := server.NewServer(":8080", router, nil)
	s.ListenAndServe()
}

func readConfigFile() map[string]string {
	var m = make(map[string]string)
	filePath := "./config.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	err = json .Unmarshal(data, &m)
	if err != nil {
		fmt.Println("Error unmarshalling JSON into struct:", err)
		os.Exit(1)
	}
	return m
}