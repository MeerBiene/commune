package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type Client struct {
	Domain          string `toml:"domain"`
	ShortlinkDomain string `toml:"shortlink_domain"`
	Port            int    `toml:"port"`
	CookieName      string `toml:"cookie_name"`
	SecureCookie    string `toml:"secure_cookie"`
}

type Matrix struct {
	Server           string `toml:"server"`
	FederationServer string `toml:"federation_server"`
	PublicServer     string `toml:"public_server"`
	Port             int    `toml:"port"`
	Username         string `toml:"username"`
	Password         string `toml:"password"`
}

type DB struct {
	Client string `toml:"client"`
	Matrix string `toml:"matrix"`
}

type Redis struct {
	Address  string `toml:"address"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type Spaces struct {
	Prefix string `toml:"prefix"`
}

type Auth struct {
	VerifyEmail            bool   `toml:"verify_email"`
	DisableRegistration    bool   `toml:"disable_registration"`
	DisableFederatedLogin  bool   `toml:"disable_federated_login"`
	DisableProfileCreation bool   `toml:"disable_profile_creation"`
	SharedSecret           string `toml:"shared_secret"`
}

type Privacy struct {
	DisablePublic bool `toml:"disable_public"`
}

type Email struct {
	Account   string `toml:"account"`
	Server    string `toml:"server"`
	Port      int    `toml:"port"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	Templates struct {
		VerifyEmail   int64 `toml:"verify_email"`
		PasswordReset int64 `toml:"password_reset"`
	} `toml:"templates"`
}

type Tenor struct {
	Key string `toml:"key"`
}

type Github struct {
	Secret string `toml:"secret"`
}

type Oauth2 struct {
	Github struct {
		ClientID     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
	} `toml:"github"`
	Google struct {
		ClientID     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
	} `toml:"google"`
	Discord struct {
		ClientID     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
	} `toml:"discord"`
	Reddit struct {
		ClientID     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
	} `toml:"reddit"`
}

type Config struct {
	Name       string  `toml:"name"`
	Mode       string  `toml:"mode"`
	Client     Client  `toml:"client"`
	Matrix     Matrix  `toml:"matrix"`
	DB         DB      `toml:"db"`
	Redis      Redis   `toml:"redis"`
	YoutubeKey string  `toml:"youtube_key"`
	Spaces     Spaces  `toml:"spaces"`
	Auth       Auth    `toml:"auth"`
	Privacy    Privacy `toml:"privacy"`
	Email      Email   `toml:"email"`
	Tenor      Tenor   `toml:"tenor"`
	Github     Github  `toml:"github"`
	Oauth2     Oauth2  `toml:"oauth_2"`
}

var conf Config

// Read reads the config file and returns the Values struct
func Read(s string) (*Config, error) {
	file, err := os.Open(s)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if _, err := toml.Decode(string(b), &conf); err != nil {
		panic(err)
	}

	return &conf, err
}
