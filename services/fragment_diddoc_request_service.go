package services

import "strings"


type FragmentDIDDocRequestService struct {
	BaseDidDocRequestService
}

func (dd *FragmentDIDDocRequestService) SpecificValidation(c ResolverContext) error {
	return nil
}

func (dd *FragmentDIDDocRequestService) Prepare(c ResolverContext) error {
	splitted := strings.Split(c.Param("did"), "#")

	if len(splitted) == 2 {
		dd.fragment = splitted[1]
	}
	return nil
}

func (dd *FragmentDIDDocRequestService) MakeAnswer(c ResolverContext) error {
	return nil
}