package appConfig

import (
	"os"
	"strconv"
	"strings"
	"time"
	"fmt"
	"github.com/joho/godotenv"
	"cherry-blossom-hunters-app/logger"
)

func init() {
	initialize()
}

func initialize() {
	fmt.Println(os.Environ())
    if os.Getenv("APP_ENV") != "production" {
        err := godotenv.Load()
        if err != nil {
            logger.Logging("[appConfig] No .env file found or failed to load:", logger.Error)
        }
    }
}

type Config struct {
	Server struct {
		Port string
		Host string
	}

	Discord struct {
		Token                             string
		MemberChannelID                   string
		RuleChannelID                     string
		MasterUserID                      string
		EventChannelID                    string
		GoodReactionUrlEncodeString       string
		SubjectToUrlExclusionFromJedgment string
	}

	Event struct {
		ScriptPath     string
		PythonCommand  string
		OutputFilePath string
		DefaultTimeout time.Duration
	}

	Scraping struct {
		EventURL string
		Timeout  int
	}

	Storage struct {
		HashFilePath   string
		MemberDataPath string
	}

	Notify struct {
		Webhooks map[string]string
	}
}

// GetConfig returns the full app configuration
func GetConfig() *Config {
	return &Config{
		Server: struct {
			Port string
			Host string
		}{
			Port: getEnvWithDefault("SERVER_PORT", "8080"),
			Host: getEnvWithDefault("SERVER_HOST", "localhost"),
		},

		Discord: struct {
			Token                             string
			MemberChannelID                   string
			RuleChannelID                     string
			MasterUserID                      string
			EventChannelID                    string
			GoodReactionUrlEncodeString       string
			SubjectToUrlExclusionFromJedgment string
		}{
			Token:                             getEnvRequired("DISCORD_BOT_CLIENT_TOKEN"),
			MemberChannelID:                   getEnvRequired("MEMBER_CHANNEL_ID"),
			RuleChannelID:                     getEnvRequired("RULE_CHANNEL_ID"),
			MasterUserID:                      getEnvRequired("MASTER_USER_ID"),
			EventChannelID:                    getEnvRequired("EVENT_CHANNEL_ID"),
			GoodReactionUrlEncodeString:       getEnvRequired("GOOD_REACTION_URL_ENCODE_STRING"),
			SubjectToUrlExclusionFromJedgment: getEnvRequired("SUBJECT_TO_EXCLUSION_FROM_JUDGMENT"),
		},

		Event: struct {
			ScriptPath     string
			PythonCommand  string
			OutputFilePath string
			DefaultTimeout time.Duration
		}{
			ScriptPath:     getEnvWithDefault("EVENT_SCRIPT_PATH", "./script/event-scraper.py"),
			PythonCommand:  getEnvWithDefault("EVENT_PYTHON_COMMAND", "python3"),
			OutputFilePath: getEnvWithDefault("EVENT_OUTPUT_FILE_PATH", "./sha256_output.txt"),
			DefaultTimeout: getEnvDuration("EVENT_DEFAULT_TIMEOUT", 10*time.Second),
		},

		Scraping: struct {
			EventURL string
			Timeout  int
		}{
			EventURL: getEnvWithDefault("SCRAPING_EVENT_URL", "https://www.capcom.co.jp/monsterhunter/wilds/topics/"),
			Timeout:  getEnvInt("SCRAPING_TIMEOUT", 30),
		},

		Storage: struct {
			HashFilePath   string
			MemberDataPath string
		}{
			HashFilePath:   getEnvWithDefault("HASH_FILE_PATH", "sha256_output.txt"),
			MemberDataPath: getEnvWithDefault("MEMBER_DATA_PATH", "member_compliance_data.json"),
		},

		Notify: struct {
			Webhooks map[string]string
		}{
			Webhooks: getWebhookMap([]string{
				"DISCORD_WEBHOOK_MEMBER",
				"DISCORD_WEBHOOK_EVENT",
				"DISCORD_WEBHOOK_LOG",
			}),
		},
	}
}

func getEnvRequired(key string) string {
	val := os.Getenv(key)
	if val == "" {
		fmt.Println("Available environment variables:")
		for _, env := range os.Environ() {
			if strings.Contains(env, "DISCORD") || strings.Contains(env, "TOKEN") {
				fmt.Println(env)
			}
		}
		panic("Required environment variable is not set: " + key)
	}
	return val
}
func getEnvWithDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if sec, err := strconv.Atoi(val); err == nil {
			return time.Duration(sec) * time.Second
		}
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if val := os.Getenv(key); val != "" {
		slice := strings.Split(val, ",")
		for i := range slice {
			slice[i] = strings.TrimSpace(slice[i])
		}
		return slice
	}
	return defaultValue
}

// getWebhookMap loads multiple DISCORD_WEBHOOK_* env vars into a map
func getWebhookMap(keys []string) map[string]string {
	webhooks := make(map[string]string)
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			name := strings.ToLower(strings.TrimPrefix(key, "DISCORD_WEBHOOK_")) // e.g. MEMBER → member
			webhooks[name] = val
		}
	}
	return webhooks
}
