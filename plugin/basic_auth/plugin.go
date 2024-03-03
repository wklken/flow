package basic_auth

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

const (
	// version  = "0.1"
	priority = 102
	name     = "basic_auth"
)

type Plugin struct {
	config Config
}

// FIXME: use jsonschema to unmarshal the config dynamic

type Config struct {
	Credentials map[string]string `mapstructure:"credentials"`
	Realm       string            `mapstructure:"realm"`
}

func (p *Plugin) Name() string {
	return name
}

func (p *Plugin) Priority() int {
	return priority
}

func (p *Plugin) Init(config string) error {
	fmt.Println("init the basic_auth plugin", config)
	v := viper.New()
	v.SetConfigType("json")

	// TODO: how to make the default value
	// v.SetDefault("header_name", "X-Request-ID")
	// v.SetDefault("set_in_response", true)

	v.ReadConfig(bytes.NewBuffer([]byte(config)))

	fmt.Println("config: ", v.AllSettings())

	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		return err
	}
	fmt.Printf("config: %+v\n", c)
	p.config = c

	return nil
}

func (p *Plugin) Handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			basicAuthFailed(w, p.config.Realm)
			return
		}

		credPass, credUserOk := p.config.Credentials[user]
		if !credUserOk || pass != credPass {
			basicAuthFailed(w, p.config.Realm)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func basicAuthFailed(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
}
