package utils

func BytesToMB(bytes int) float64 {
	return float64(bytes) / (1024 * 1024)
}
