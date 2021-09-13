package configuration

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func NewReader() *Reader { // TODO: extract to external library
	v := viper.New()
	v.AutomaticEnv()
	v.AllowEmptyEnv(true)
	return &Reader{
		Mandatory: Mandatory{
			viper: v,
		},
		Optional: Optional{
			viper: v,
		},
	}
}

type Reader struct {
	Mandatory Mandatory
	Optional  Optional
}

func (r *Reader) Error() error {
	if len(r.Mandatory.missingKeys) > 0 {
		return errors.Errorf("missing keys: %s", strings.Join(r.Mandatory.missingKeys, ", "))
	}
	return nil
}
