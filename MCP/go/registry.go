package main

import (
	"github.com/sonar-trading/mcp-server/config"
	"github.com/sonar-trading/mcp-server/models"
	tools_currencies "github.com/sonar-trading/mcp-server/tools/currencies"
)

func GetAll(cfg *config.APIConfig) []models.Tool {
	return []models.Tool{
		tools_currencies.CreateGet_convertTool(cfg),
		tools_currencies.CreateGet_country_currenciesTool(cfg),
		tools_currencies.CreateGet_digital_currenciesTool(cfg),
		tools_currencies.CreateGet_historyTool(cfg),
	}
}
