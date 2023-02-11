package main

import (
	"github.com/cheqd/did-resolver/cmd"
	_ "github.com/cheqd/did-resolver/docs"
)

// @title DID Resolver for did:cheqd method
// @version 1.0
// @description	Universal Resolver driver for did:cheqd method
// @contact.name Cheqd Foundation Limited
// @contact.url	https://cheqd.io
// @license.name Apache 2.0
// @license.url	https://github.com/cheqd/did-resolver/blob/main/LICENSE
// @host resolver.cheqd.net
// @BasePath /1.0/identifiers
// @schemes	https http
func main() {
	rootCmd := cmd.GetRootCmd()

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
