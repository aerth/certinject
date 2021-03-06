// Copyright 2020 Namecoin Developers GPLv3+

// Command certinject injects certificates into all configured trust stores
package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/namecoin/ncdns/certinject"
	//	flag "github.com/ogier/pflag"
	"flag"

	"gopkg.in/hlandau/easyconfig.v1/adaptflag"
	"gopkg.in/hlandau/easyconfig.v1/cflag"
)

func main() {
	log.SetFlags(0)
	log.Println("Namecoin certinject tool")

	var (
		certs      []string
		flagGroup  = cflag.NewGroup(nil, "certinject")
		certflag   = cflag.String(flagGroup, "certs", "", "path to certificate (separate by comma. if set, skips config)")
		configflag = cflag.String(flagGroup, "conf", filepath.Join(getConfigDir(), "certinject.toml"), "path to config")
	)

	adaptflag.AdaptWithFunc(func(info adaptflag.Info) {
		flag.Var(info.Value, info.Name, info.Usage)
	})
	flag.Parse()
	certs = listCerts(certflag.Value(), configflag.Value())
	if len(certs) == 0 {
		log.Fatalln("no certificates to add")
	}
	log.Printf("injecting %v certificates", len(certs))
	for _, cert := range certs {
		b, err := ioutil.ReadFile(cert)
		if err != nil {
			log.Fatalf("fatal error while injecting %q certificate: \n\t%v", cert, err)
		}
		certinject.InjectCert(b)
	}
	log.Printf("injected %v certificates", len(certs))
}

// configt config type
type configt struct {
	Certs []string
}

func listCerts(certflag, configflag string) (certs []string) {
	if certflag == "" {
		// no -certs flag, parse TOML config
		config, err := readConfigCerts(configflag)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("fatal: there is no toml config file at %q, "+
					"and none specified with %q flag",
					configflag, "certs")
			} else {
				log.Fatalln("error parsing config:", err)
			}
		}
		certs = config.Certs
	} else {
		// read -certs flag for comma separated certificate paths
		certs = strings.Split(certflag, ",")
	}
	return certs
}

// readConfigCerts parses a toml file and returns a list of paths to certificate files
func readConfigCerts(path string) (configt, error) {
	if path == "" {
		return configt{}, errors.New("empty config path")
	}
	var config = configt{}
	_, err := toml.DecodeFile(path, &config)
	return config, err
}

// getConfigDir always will return a valid directory to look for default config file
func getConfigDir() string {
	configdir, err := os.UserConfigDir()
	if err != nil {
		log.Println("error looking for user config dir, using current working dir:", err)
		return ""
	}

	subpath := "Namecoin"
	if runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
		subpath = "namecoin"
	}

	return filepath.Join(configdir, subpath)
}
