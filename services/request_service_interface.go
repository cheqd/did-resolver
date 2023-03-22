package services

import (
	echo "github.com/labstack/echo/v4"
)

type RequestServiceI interface {
	// Checks
	IsRedirectNeeded(c ResolverContext) bool

	// Methods
	// Setup
	Setup(c ResolverContext) error
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
	SetupResponse(c ResolverContext) error

	Respond(c ResolverContext) error
}

// The main flow for all the requests
func EchoWrapHandler(controller RequestServiceI) echo.HandlerFunc {
	return func(c echo.Context) error {
		rc := c.(ResolverContext)
		// Setup
		if err := controller.Setup(rc); err != nil {
			return err
		}
		// Preparations, like get parameters from context and others
		if err := controller.BasicPrepare(rc); err != nil {
			return err
		}
		if err := controller.SpecificPrepare(rc); err != nil {
			return err
		}
		// Redirect if needed
		if controller.IsRedirectNeeded(rc) {
			return controller.Redirect(rc)
		}
		// Validation
		if err := controller.BasicValidation(rc); err != nil {
			return err
		}
		if err := controller.SpecificValidation(rc); err != nil {
			return err
		}
		// Query
		if err := controller.Query(rc); err != nil {
			return err
		}
		// Make response. Set specific headers, etc.
		if err := controller.SetupResponse(rc); err != nil {
			return err
		}
		return controller.Respond(rc)
	}
}
