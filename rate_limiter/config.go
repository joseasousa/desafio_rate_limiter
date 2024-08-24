package rate_limiter

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/joseasousa/rate_limiter/rate_limiter/adapter"
	"github.com/joseasousa/rate_limiter/rate_limiter/response_writer"
)

const envKeyIPMaxRequestsPerSecond = "RATE_LIMITER_IP_MAX_REQUESTS"
const envKeyIPBlockTimeMilliseconds = "RATE_LIMITER_IP_BLOCK_TIME"
const envKeyTokenMaxRequestsPerSecond = "RATE_LIMITER_TOKEN_MAX_REQUESTS"
const envKeyTokenBlockTimeMilliseconds = "RATE_LIMITER_TOKEN_BLOCK_TIME"
const envKeyDebug = "RATE_LIMITER_DEBUG"
const envUseRedis = "RATE_LIMITER_USE_REDIS"
const envRedisAddress = "RATE_LIMITER_REDIS_ADDRESS"
const envRedisPassword = "RATE_LIMITER_REDIS_PASSWORD"
const envRedisDB = "RATE_LIMITER_REDIS_DB"

type rateLimiterRateConfig struct {
	MaxRequestsPerSecond  int64 `json:"maxRequestsPerSecond"`
	BlockTimeMilliseconds int64 `json:"blockTimeMilliseconds"`
}

type rateLimiterConfig struct {
	IP             *rateLimiterRateConfig                            `json:"ip"`
	Token          *rateLimiterRateConfig                            `json:"token"`
	CustomTokens   *map[string]*rateLimiterRateConfig                `json:"tokens"`
	StorageAdapter adapter.RateLimitStorageAdapter                   `json:"-"`
	ResponseWriter *response_writer.RateLimiterDefaultResponseWriter `json:"-"`
	Debug          bool                                              `json:"debug"`
	DisableEnvs    bool                                              `json:"disableEnvs"`
}

func (c *rateLimiterConfig) GetrateLimiterRateConfigForToken(token string) (*rateLimiterRateConfig, bool) {
	customTokenConfig, ok := (*c.CustomTokens)[token]
	if ok {
		return customTokenConfig, true
	} else {
		return c.Token, false
	}
}

func getDefaultConfiguration() *rateLimiterConfig {
	return &rateLimiterConfig{
		IP: &rateLimiterRateConfig{
			MaxRequestsPerSecond:  100,
			BlockTimeMilliseconds: 1000,
		},
		Token: &rateLimiterRateConfig{
			MaxRequestsPerSecond:  200,
			BlockTimeMilliseconds: 500,
		},
		CustomTokens:   &map[string]*rateLimiterRateConfig{},
		StorageAdapter: adapter.NewRateLimitMemoryStorageAdapter(),
		ResponseWriter: response_writer.NewRateLimiterDefaultResponseWriter(),
		Debug:          false,
	}
}

func setConfiguration(config *rateLimiterConfig) *rateLimiterConfig {
	defaultConfiguration := getDefaultConfiguration()

	if config == nil {
		config = defaultConfiguration
	}

	if !config.DisableEnvs {
		debug, ok := getBoolEnv(envKeyDebug)
		if ok {
			config.Debug = debug
			DebugPrintfWithoutKey(config, "using env %s", envKeyDebug)
		}
	}

	configureIP(config, defaultConfiguration)
	configureToken(config, defaultConfiguration)
	configureCustomTokens(config, defaultConfiguration)
	configureStorageAdapter(config, defaultConfiguration)
	configureResponseWriter(config, defaultConfiguration)

	if config.Debug {
		jsonConfiguration, err := json.Marshal(config)
		if err == nil {
			DebugPrintfWithoutKey(config, "using configuration: %s", jsonConfiguration)
		}
	}

	return config
}

func configureIP(config *rateLimiterConfig, defaultConfiguration *rateLimiterConfig) {
	if config.IP == nil {
		config.IP = defaultConfiguration.IP
	}

	if !config.DisableEnvs {
		mrps, ok := getInt64Env(envKeyIPMaxRequestsPerSecond)
		if ok {
			config.IP.MaxRequestsPerSecond = mrps
			DebugPrintfWithoutKey(config, "using env %s", envKeyIPMaxRequestsPerSecond)
		}

		bt, ok := getInt64Env(envKeyIPBlockTimeMilliseconds)
		if ok {
			config.IP.BlockTimeMilliseconds = bt
			DebugPrintfWithoutKey(config, "using env %s", envKeyIPBlockTimeMilliseconds)
		}
	}
}

func configureToken(config *rateLimiterConfig, defaultConfiguration *rateLimiterConfig) {
	if config.Token == nil {
		config.Token = defaultConfiguration.Token
	}

	if !config.DisableEnvs {
		mrps, ok := getInt64Env(envKeyTokenMaxRequestsPerSecond)
		if ok {
			config.Token.MaxRequestsPerSecond = mrps
			DebugPrintfWithoutKey(config, "using env %s", envKeyTokenMaxRequestsPerSecond)
		}

		bt, ok := getInt64Env(envKeyTokenBlockTimeMilliseconds)
		if ok {
			config.Token.BlockTimeMilliseconds = bt
			DebugPrintfWithoutKey(config, "using env %s", envKeyTokenBlockTimeMilliseconds)
		}
	}
}

func configureCustomTokens(config *rateLimiterConfig, defaultConfiguration *rateLimiterConfig) {
	if config.CustomTokens == nil {
		config.CustomTokens = defaultConfiguration.CustomTokens
	}

	for key := range *config.CustomTokens {
		value, ok := (*config.CustomTokens)[key]
		if !ok || value == nil {
			(*config.CustomTokens)[key] = config.Token
		}
	}

	customTokens := getCustomTokenList()
	for _, customToken := range *customTokens {
		configureCustomToken(config, customToken)
	}
}

func getCustomTokenList() *[]string {
	envKeyRegex := regexp.MustCompile("^RATE_LIMITER_TOKEN_(.*)_(MAX_REQUESTS|BLOCK_TIME)$")

	foundTokens := map[string]bool{}

	envs := os.Environ()
	for _, env := range envs {
		envPair := strings.SplitN(env, "=", 2)
		envKey := envPair[0]
		if envKeyRegex.Match([]byte(envKey)) {
			foundTokens[envKeyRegex.FindStringSubmatch(envKey)[1]] = true
		}
	}

	tokens := []string{}
	for k := range foundTokens {
		tokens = append(tokens, k)
	}

	return &tokens
}

func configureCustomToken(config *rateLimiterConfig, customToken string) {

	DebugPrintfWithoutKey(config, "configuring custom token \"%s\"", customToken)

	maxRequestsPerSecondEnvKey := fmt.Sprintf("RATE_LIMITER_TOKEN_%s_MAX_REQUESTS", customToken)
	maxRequestsPerSecond, ok := getInt64Env(maxRequestsPerSecondEnvKey)
	if !ok {
		defaultValue := config.Token.MaxRequestsPerSecond
		DebugPrintfWithoutKey(config, "env \"%s\" not found: using default value %d", maxRequestsPerSecondEnvKey, defaultValue)
		maxRequestsPerSecond = defaultValue
	}

	blockTimeMillisecondEnvKey := fmt.Sprintf("RATE_LIMITER_TOKEN_%s_BLOCK_TIME", customToken)
	blockTimeMilliseconds, ok := getInt64Env(blockTimeMillisecondEnvKey)
	if !ok {
		defaultValue := config.Token.BlockTimeMilliseconds
		DebugPrintfWithoutKey(config, "env \"%s\" not found: using default value %d", blockTimeMillisecondEnvKey, defaultValue)
		blockTimeMilliseconds = defaultValue
	}

	(*config.CustomTokens)[customToken] = &rateLimiterRateConfig{
		MaxRequestsPerSecond:  maxRequestsPerSecond,
		BlockTimeMilliseconds: blockTimeMilliseconds,
	}
}

func configureStorageAdapter(config *rateLimiterConfig, defaultConfiguration *rateLimiterConfig) {
	if config.StorageAdapter == nil {
		config.StorageAdapter = defaultConfiguration.StorageAdapter
	}

	useRedis, ok := getBoolEnv(envUseRedis)
	if ok && useRedis {
		configureRedisStorageAdapter(config)
	} else if config.StorageAdapter != defaultConfiguration.StorageAdapter {
		DebugPrintfWithoutKey(config, "using StorageAdapter Custom")
	} else {
		DebugPrintfWithoutKey(config, "using StorageAdapter Default")
	}
}

func configureRedisStorageAdapter(config *rateLimiterConfig) {
	DebugPrintfWithoutKey(config, "using StorageAdapter Redis")

	redisAddress, ok := getStringEnv(envRedisAddress)
	if !ok {
		panic(fmt.Sprintf("%s env is required when using redis adapter with env configuration", envRedisAddress))
	}

	redisPassword, ok := getStringEnv(envRedisPassword)
	if !ok {
		redisPassword = ""
	}

	redisDB, ok := getInt64Env(envRedisDB)
	if !ok {
		redisDB = 0
	}

	config.StorageAdapter = adapter.NewRateLimitRedisStorageAdapter(redisAddress, redisPassword, redisDB)
}

func configureResponseWriter(config *rateLimiterConfig, defaultConfiguration *rateLimiterConfig) {
	if config.ResponseWriter == nil {
		config.ResponseWriter = defaultConfiguration.ResponseWriter
	}

	if config.ResponseWriter != defaultConfiguration.ResponseWriter {
		DebugPrintfWithoutKey(config, "using ResponseWriter Custom")
	} else {
		DebugPrintfWithoutKey(config, "using ResponseWriter Default")
	}
}
