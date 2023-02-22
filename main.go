package main

import (
	"github.com/cheqd/did-resolver/cmd"
)

//	@title			DID Resolver for cheqd DID method
//	@version		v3.0
//	@description	Universal Resolver driver for cheqd DID method
//	@contact.name	Cheqd Foundation Limited
//	@contact.url	https://cheqd.io
//	@license.name	Apache 2.0
//	@license.url	https://github.com/cheqd/did-resolver/blob/main/LICENSE
//	@host			resolver.cheqd.net
//	@BasePath		/1.0/identifiers
//	@schemes		https http

func main() {
	rootCmd := cmd.GetRootCmd()

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
