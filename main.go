package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/armon/go-socks5"
	"github.com/spf13/viper"
)

var (
	port       = flag.String("port", "8086", "The port to listen.")
	configFile = flag.String("configFile", "config.yml", "configuration file")
	cred       = socks5.StaticCredentials{}
)

func main() {
	flag.Parse()

	viper.SetConfigType("yaml")
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %s \n\n", err)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	credentials := viper.Get("credentials").([]interface{})
	for i := range credentials {
		c := credentials[i].(map[interface{}]interface{})
		cred[c["user"].(string)] = c["password"].(string)
	}

	auth := socks5.UserPassAuthenticator{Credentials: cred}
	conf := &socks5.Config{AuthMethods: []socks5.Authenticator{auth}}

	s, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy
	if err := s.ListenAndServe("tcp", ":"+*port); err != nil {
		panic(err)
	}

}
