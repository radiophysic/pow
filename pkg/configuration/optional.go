package configuration

import (
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Optional struct {
	viper *viper.Viper
}

func (o *Optional) Bool(key string, value bool) bool {
	if !o.viper.IsSet(key) {
		return value
	}
	return o.viper.GetBool(key)
}

func (o *Optional) Int(key string, value int) int {
	if !o.viper.IsSet(key) {
		return value
	}
	return o.viper.GetInt(key)
}

func (o *Optional) String(key string, value string) string {
	if !o.viper.IsSet(key) {
		return value
	}
	return o.viper.GetString(key)
}

func (o *Optional) Duration(key string, value time.Duration) time.Duration {
	if !o.viper.IsSet(key) {
		return value
	}
	return o.viper.GetDuration(key)
}

// DurationSlice reads duration slice from env (list should be space separated, for example "1h 2s 3m").
func (o *Optional) DurationSlice(key string, value []time.Duration) []time.Duration {
	if !o.viper.IsSet(key) {
		return value
	}
	return cast.ToDurationSlice(o.viper.GetStringSlice(key))
}
