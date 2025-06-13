package hyperliquid

import (
	"github.com/joho/godotenv"
	"os"
)

func init() {
	if _, err := os.Stat(".test.env"); err == nil {
		_ = godotenv.Load(".test.env")
	} else if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load()
	}
}
