package rate_limiter

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func DebugPrintf(config *rateLimiterConfig, format string, keyType string, key string, a ...any) (n int, err error) {
	if config.Debug {
		timeString := time.Now().UTC().Format("2000-01-02 11:01:02")
		args := []any{timeString, keyType, key}
		args = append(args, a...)
		return fmt.Printf("%s [RATE LIMITER][%s][%s] "+format+"\n", args...)
	}

	return 0, nil
}

func DebugPrintfWithoutKey(config *rateLimiterConfig, format string, a ...any) (n int, err error) {
	if config.Debug {
		timeString := time.Now().UTC().Format("2000-01-02 11:01:02")
		args := []any{timeString}
		args = append(args, a...)
		return fmt.Printf("%s [RATE LIMITER] "+format+"\n", args...)
	}

	return 0, nil
}

func GetRemainingBlockTime(block *time.Time) float64 {
	return time.Until(*block).Seconds()
}

func getStringEnv(key string) (string, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", false
	}
	if value == "" {
		return "", false
	}
	return value, true
}

func getBoolEnv(key string) (bool, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return false, false
	}
	if value == "" {
		return false, false
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, false
	}
	return parsed, true
}

func getInt64Env(key string) (int64, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, false
	}
	if value == "" {
		return 0, false
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}
