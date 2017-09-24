package utils

import (
	"github.com/google/uuid"
	"math/rand"
	"strings"

	"github.com/SestroAI/shared/config"
	"strconv"
)

func GenerateUUID() string {
	uuid := uuid.New().String()
	return uuid
}

func randomInRange(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func GenerateIDFromName(name string, addRandom bool) string {
	name = strings.Trim(strings.ToLower(name), "!.:\\/,&@")
	id := strings.Replace(name, " ", "-", -1)

	if addRandom {
		randomInt := randomInRange(1000, 9999)
		extension := strconv.FormatInt(int64(randomInt), 10)
		id += extension
	}
	return id
}

func GetServicePrefix() string {
	return "/v1/" + config.ServiceName
}
