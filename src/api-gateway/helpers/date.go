package helpers

import "time"

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02") // Example: "2023-01-12"
}
