package viper

import (
	"fmt"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_Env(t *testing.T) {
	viper.AutomaticEnv()
	viper.SetConfigFile("/Users/jiang/go/src/account/config/config.yml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := godotenv.Load("/Users/jiang/go/src/account/.env")
	assert.NoError(t, err)

	err = viper.ReadInConfig()
	assert.NoError(t, err)

	fmt.Println(viper.GetString("db.max.idle.conns"))
}
