// Copyright 2020 Namecoin Developers GPLv3+

// Command certinject injects certificates into all configured trust stores
package main

import (
	"encoding/pem"
	"io/ioutil"
	"os"

	"github.com/hlandau/xlog"
	"github.com/namecoin/ncdns/certinject"
	easyconfig "gopkg.in/hlandau/easyconfig.v1"
	"gopkg.in/hlandau/easyconfig.v1/cflag"
)

var log, _ = xlog.New("certinject")

func main() {
	var (
		flagGroup = cflag.NewGroup(nil, "certinject")
		certflag  = cflag.String(flagGroup, "cert", "", "path to certificate to inject into trust store")
	)

	// read config
	config := easyconfig.Configurator{
		ProgramName: "certinject",
	}
	config.ParseFatal(nil)

	cert := os.ExpandEnv(certflag.Value())
	log.Debugf("reading certificate: %q", cert)
	b, err := ioutil.ReadFile(cert)
	if err != nil {
		log.Fatalf("fatal error while injecting %q certificate: \n\t%v", cert, err)
	}
	if p, err := pem.Decode(b); err == nil {
		log.Debugf("user provided PEM encoded certificate, extracting DER bytes")
		b = p.Bytes
	}
	certinject.InjectCert(b)
	log.Debugf("injected certificate: %q", cert)
}
