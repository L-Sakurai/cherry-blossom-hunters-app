package jsonUtil

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cherry-blossom-hunters-app/logger"
)

func PersistJSONToHash256(jsonBytes []byte, outpath string) error {
	var jsonMap interface{}
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		logger.Logging("Failed to parse json: " + err.Error(), logger.Error)
		return err
	}

	normalizedJSON, err := json.Marshal(jsonMap)
	if err != nil {
		logger.Logging("Failed to re-marshal json.", logger.Error)
		return err
	}

	hash := sha256.Sum256(normalizedJSON)
	hashString := fmt.Sprintf("%x", hash)

	isMatch, err := compareJSON_HashCode(hashString)
	if !isMatch {
		logger.Logging("[difference hashcode]", logger.Debug)
		if err := os.WriteFile(outpath, []byte(hashString), 0644); err != nil {
			logger.Logging("Failed to write file.", logger.Error)
			return err
		}
	}
	logger.Logging("HashCode: " + hashString)
	return nil
}

func compareJSON_HashCode(hashString string) (bool, error) {
	hashCode, err := os.ReadFile("./sha256_output.txt")
	if err != nil {
		logger.Logging("Failed to read hash file.", logger.Error)
		return false, err
	}

	storedHashCode := strings.TrimSpace(string(hashCode))
	logger.Logging("[hashString]" + hashString, logger.Debug)
	logger.Logging("[storedHashCode]"+ storedHashCode, logger.Debug)
	isMatch := storedHashCode == hashString
	return isMatch, nil
}
