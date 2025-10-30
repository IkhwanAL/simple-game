package ui

import (
	"strconv"

	"github.com/ikhwanal/tinyworlds/internal/world"
)

func CeilClass(c world.CellType) string {
	switch c {
	case world.Food:
		return "w-5 h-5 border border-gray-800 bg-green-500"
	case world.AgentEn:
		return "w-5 h-5 border border-gray-800 bg-red-500"
	default:
		return "w-5 h-5 border border-gray-800 bg-gray-900"
	}
}

func GetColorCeilType(c world.CellType) string {
	switch c {
	case world.Food:
		return "bg-green-400"
	case world.AgentEn:
		return "bg-red-500"
	default:
		return "bg-white-900"
	}
}

func GetAgentViewStyple(a world.Agent) map[string]string {
	return map[string]string{
		"left": strconv.Itoa(a.X*16) + "px;",
		"top":  strconv.Itoa(a.Y*16) + "px;",
	}
}
