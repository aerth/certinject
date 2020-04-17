// Copyright 2020 Namecoin Developers GPLv3+

// Command certinject injects certificates into all configured trust stores
package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hlandau/xlog"
	"github.com/namecoin/ncdns/certinject"
	easyconfig "gopkg.in/hlandau/easyconfig.v1"
	"gopkg.in/hlandau/easyconfig.v1/cflag"
)

func main() {
	var (
		log, _    = xlog.New("certinject")
		flagGroup = cflag.NewGroup(nil, "certinject")
		certflag  = cflag.String(flagGroup, "certs", "", "path to certificate (separate by comma. if set, skips config)")
	)

	// read config
	config := easyconfig.Configurator{
		ProgramName: "certinject",
	}
	config.ParseFatal(nil)

	certs := listCerts(certflag.Value())
	if len(certs) == 0 {
		log.Fatal("no certificates to add")
	}

	log.Debugf("injecting %v certificates", len(certs))
	for _, cert := range certs {
		log.Debugf("reading certificate: %q", cert)
		b, err := ioutil.ReadFile(cert)
		if err != nil {
			log.Fatalf("fatal error while injecting %q certificate: \n\t%v", cert, err)
		}
		certinject.InjectCert(b)
		log.Debugf("injected certificate: %q", cert)
	}
	log.Debugf("injected %v certificates", len(certs))
}

// configt config type
type configt struct {
	Certs []string
}

func listCerts(certificates string) (certs []string) {
	certificates = strings.TrimSpace(certificates)
	if certificates == "" {
		return nil
	}
	return strings.Split(certificates, ",")
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
