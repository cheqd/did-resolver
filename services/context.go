package services

import (
	"github.com/labstack/echo/v4"
)

type ResolverContext struct {
	echo.Context
	LedgerService   LedgerServiceI
	DidDocService   DIDDocService
	ResourceService ResourceService
}
