package ui

import "github.com/ikhwanal/tinyworlds/world"

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
