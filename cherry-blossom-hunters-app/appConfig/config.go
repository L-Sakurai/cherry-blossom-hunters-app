package appConfig

import (
	"os"
	"strconv"
    "strings"

	"github.com/joho/godotenv"
	"cherry-blossom-hunters-app/logger"
)

func init() {
    initialize()
}

func initialize() {
    err := godotenv.Load()
    if err != nil {
        logger.Logging("[appConfig] No .env file found or failed to load:", logger.Error)
    }
}

type Config struct {
    Server struct {
        Port string
        Host string
    }
    
    Discord struct {
        Token                string
        MemberChannelID      string
        RuleChannelID        string
        MasterUserID         string
        EventChannelID       string
        GoodReactionUrlEncodeString string
    }
    
    Scraping struct {
        EventURL string
        Timeout  int // 秒
    }
    
    Storage struct {
        HashFilePath   string
        MemberDataPath string
    }
}

// 環境変数から設定を読み込み
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
            Token                string
            MemberChannelID      string
            RuleChannelID        string
            MasterUserID         string
            EventChannelID       string
            GoodReactionUrlEncodeString string
        }{
            Token:           getEnvRequired("DISCORD_BOT_CLIENT_TOKEN"),
            MemberChannelID: getEnvRequired("MEMBER_CHANNEL_ID"),
            RuleChannelID:   getEnvRequired("RULE_CHANNEL_ID"),
            MasterUserID:    getEnvRequired("MASTER_USER_ID"),
            EventChannelID:  getEnvRequired("EVENT_CHANNEL_ID"),
            GoodReactionUrlEncodeString: getEnvRequired("GOOD_REACTION_URL_ENCODE_STRING"),
        },
        
        Scraping: struct {
            EventURL string
            Timeout  int
        }{
            EventURL: getEnvWithDefault("SCRAPING_EVENT_URL", "https://www.capcom.co.jp/monsterhunter/wilds/topics/"),
            Timeout:  getEnvInt("SCRAPING_TIMEOUT", 30), // デフォルト30秒
        },
        
        Storage: struct {
            HashFilePath   string
            MemberDataPath string
        }{
            HashFilePath:   getEnvWithDefault("HASH_FILE_PATH", "sha256_output.txt"),
            MemberDataPath: getEnvWithDefault("MEMBER_DATA_PATH", "member_compliance_data.json"),
        },
    }
}

// 必須の環境変数を取得（存在しない場合はpanicする）
func getEnvRequired(key string) string {
    value := os.Getenv(key)
    if value == "" {
        panic("必須の環境変数が設定されていません: " + key)
    }
    return value
}

// デフォルト値付きで環境変数を取得
func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// 環境変数を整数として取得
func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

// 環境変数をスライスとして取得（カンマ区切り）
func getEnvSlice(key string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        // カンマ区切りで分割
        slice := strings.Split(value, ",")
        // 前後の空白を削除
        for i, v := range slice {
            slice[i] = strings.TrimSpace(v)
        }
        return slice
    }
    return defaultValue
}