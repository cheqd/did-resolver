package services

import (
	echo "github.com/labstack/echo/v4"
)

type RequestServiceI interface {
	// Checks
	IsRedirectNeeded(c ResolverContext) bool

	// Methods
	// Preparations
	BasicPrepare(c ResolverContext) error
	SpecificPrepare(c ResolverContext) error

	// Validation
	BasicValidation(c ResolverContext) error
	SpecificValidation(c ResolverContext) error

	// Redirect if needed
	Redirect(c ResolverContext) error

	// Ask ledger for data
	Query(c ResolverContext) error

	// Some kind of postprocessing for response
	MakeResponse(c ResolverContext) error

	Respond(c ResolverContext) error
}

// The main flow for all the requests
func EchoWrapHandler(controller RequestServiceI) echo.HandlerFunc {
	return func(c echo.Context) error {
		rc := c.(*ResolverContext)
		if err := controller.BasicPrepare(*rc); err != nil {
			return err
		}
		if err := controller.SpecificPrepare(*rc); err != nil {
			return err
		}
		if controller.IsRedirectNeeded(*rc) {
			return controller.Redirect(*rc)
		}
		if err := controller.BasicValidation(*rc); err != nil {
			return err
		}
		if err := controller.SpecificValidation(*rc); err != nil {
			return err
		}
		if err := controller.Query(*rc); err != nil {
			return err
		}
		if err := controller.MakeResponse(*rc); err != nil {
			return err
		}
		return controller.Respond(*rc)
	}
}
