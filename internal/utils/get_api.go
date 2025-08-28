package util

import "os"

func GetAPIURL() string {
	if url := os.Getenv("API_URL"); url != "" {
		return url
	}

	env := os.Getenv("ENV")
	if env == "prod" {
		return "https://your-render-url.onrender.com/api"
	}
	// default dev
	return "http://127.0.0.1:8080/api/report"
}
