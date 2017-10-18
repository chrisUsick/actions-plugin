package main

import (
	"log"
	"os"
	"syscall"

	"github.com/chrisUsick/actions-plugin/actions"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical/plugin"
)

func main() {
	fErr, openErr := os.OpenFile("/vault/logs/action-plugin.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if openErr != nil {
		log.Fatal(openErr)
		os.Exit(1)
	}
	defer fErr.Close()
	syscall.Dup2(int(fErr.Fd()), 1) /* -- stdout */
	syscall.Dup2(int(fErr.Fd()), 2) /* -- stderr */

	log.Println("Hello actions-plugin")
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	log.Println("Parsed flags")
	log.Printf("%s", flags.Args())
	log.Printf("%s", os.Args)

	tlsConfig := apiClientMeta.GetTLSConfig()
	if tlsConfig != nil {
		log.Printf("Should skip tls verify? %t", tlsConfig.Insecure)
	} else {
		log.Println("tlsConfig is nil")
	}
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

	log.Println("calling plugin.serve")
	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: actions.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
