package main

import "auth-service/internal/api"

const configDIR = "api/configs"
const envDIR = "api/.env"

func main() {
	api.Run(configDIR, envDIR)
}
