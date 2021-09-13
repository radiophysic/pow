package configuration

import (
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Mandatory struct {
	viper       *viper.Viper
	missingKeys []string
}

func (m *Mandatory) String(key string) string {
	if !m.viper.IsSet(key) {
		m.missingKeys = append(m.missingKeys, key)
		return ""
	}
	return m.viper.GetString(key)
}

func (m *Mandatory) Bool(key string) bool {
	if !m.viper.IsSet(key) {
		m.missingKeys = append(m.missingKeys, key)
		return false
	}
	return m.viper.GetBool(key)
}

func (m *Mandatory) Int(key string) int {
	if !m.viper.IsSet(key) {
		m.missingKeys = append(m.missingKeys, key)
		return 0
	}
	return m.viper.GetInt(key)
}

func (m *Mandatory) Duration(key string) time.Duration {
	if !m.viper.IsSet(key) {
		m.missingKeys = append(m.missingKeys, key)
		return 0
	}
	return m.viper.GetDuration(key)
}

// DurationSlice reads duration slice from env (list should be space separated, for example "1h 2s 3m").
func (m *Mandatory) DurationSlice(key string) []time.Duration {
	if !m.viper.IsSet(key) {
		m.missingKeys = append(m.missingKeys, key)
		return nil
	}
	return cast.ToDurationSlice(m.viper.GetStringSlice(key))
}
