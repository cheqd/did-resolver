package types

import (
	"encoding/json"
	"net/url"
	"strings"

	did "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
)

type DidDoc struct {
	Context              []string             `json:"@context,omitempty" example:"https://www.w3.org/ns/did/v1"`
	Id                   string               `json:"id,omitempty" example:"did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47"`
	Controller           []string             `json:"controller,omitempty" example:"did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47"`
	VerificationMethod   []VerificationMethod `json:"verificationMethod,omitempty"`
	Authentication       []string             `json:"authentication,omitempty" example:"did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47#key-1"`
	AssertionMethod      []AssertionMethod    `json:"assertionMethod,omitempty"`
	CapabilityInvocation []string             `json:"capabilityInvocation,omitempty"`
	CapabilityDelegation []string             `json:"capability_delegation,omitempty"`
	KeyAgreement         []string             `json:"keyAgreement,omitempty"`
	Service              []Service            `json:"service,omitempty"`
	AlsoKnownAs          []string             `json:"alsoKnownAs,omitempty"`
}

type VerificationMethod struct {
	Context            []string    `json:"@context,omitempty"`
	Id                 string      `json:"id,omitempty"`
	Type               string      `json:"type,omitempty"`
	Controller         string      `json:"controller,omitempty"`
	PublicKeyJwk       interface{} `json:"publicKeyJwk,omitempty"`
	PublicKeyMultibase string      `json:"publicKeyMultibase,omitempty"`
	PublicKeyBase58    string      `json:"publicKeyBase58,omitempty"`
}

type VerificationMaterial interface{}

type Service struct {
	Context         []string `json:"@context,omitempty"`
	Id              string   `json:"id,omitempty" example:"did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47#service-1"`
	Type            string   `json:"type,omitempty" example:"did-communication"`
	ServiceEndpoint []string `json:"serviceEndpoint,omitempty" example:"https://example.com/endpoint/8377464"`
}

type AssertionMethod struct {
	Id                  *string             `json:"id,omitempty"`
	AssertionMethodJSON *VerificationMethod `json:"assertionMethodJSON,omitempty"`
}

func (e *AssertionMethod) MarshalJSON() ([]byte, error) {
	// If Id is present, use it
	if e.Id != nil {
		return json.Marshal(e.Id)
	} else {
		// Otherwise use the VerificationMethod
		return json.Marshal(e.AssertionMethodJSON)
	}
}

func (e *AssertionMethod) UnmarshalJSON(data []byte) error {
	// Check for null or empty value
	if string(data) == "null" || len(data) == 0 {
		e.Id = nil
		e.AssertionMethodJSON = nil
		return nil
	}

	// First attempt: Try to unmarshal as a string
	var strValue string
	if err := json.Unmarshal(data, &strValue); err == nil {
		// If successfully parsed as string and it starts with "did:cheqd"
		if strings.HasPrefix(strValue, "did:cheqd") {
			e.Id = &strValue
			e.AssertionMethodJSON = nil
			return nil
		}

		// If it's a string but not a "did:cheqd" string, it might be escaped JSON
		// Try to parse the string as VerificationMethod
		var verMethod VerificationMethod
		if jsonErr := json.Unmarshal([]byte(strValue), &verMethod); jsonErr == nil {
			e.Id = nil
			e.AssertionMethodJSON = &verMethod
			return nil
		}
	}

	return nil
}

func NewDidDoc(protoDidDoc *did.DidDoc) DidDoc {
	verificationMethods := []VerificationMethod{}
	for _, vm := range protoDidDoc.VerificationMethod {
		verificationMethods = append(verificationMethods, *NewVerificationMethod(vm))
	}

	services := []Service{}
	for _, s := range protoDidDoc.Service {
		services = append(services, *NewService(s))
	}

	assertionMethods := []AssertionMethod{}
	for _, am := range protoDidDoc.AssertionMethod {
		assertionMethods = append(assertionMethods, *NewAssertionMethod(am))
	}

	return DidDoc{
		Id:                   protoDidDoc.Id,
		Controller:           protoDidDoc.Controller,
		VerificationMethod:   verificationMethods,
		Authentication:       protoDidDoc.Authentication,
		AssertionMethod:      assertionMethods,
		CapabilityInvocation: protoDidDoc.CapabilityInvocation,
		CapabilityDelegation: protoDidDoc.CapabilityDelegation,
		KeyAgreement:         protoDidDoc.KeyAgreement,
		Service:              services,
		AlsoKnownAs:          protoDidDoc.AlsoKnownAs,
	}
}

func NewVerificationMethod(protoVerificationMethod *did.VerificationMethod) *VerificationMethod {
	verificationMethod := &VerificationMethod{
		Id:         protoVerificationMethod.Id,
		Type:       protoVerificationMethod.VerificationMethodType,
		Controller: protoVerificationMethod.Controller,
	}

	switch protoVerificationMethod.VerificationMethodType {
	case "Ed25519VerificationKey2020":
		verificationMethod.PublicKeyMultibase = protoVerificationMethod.VerificationMaterial
	case "Ed25519VerificationKey2018":
		verificationMethod.PublicKeyBase58 = protoVerificationMethod.VerificationMaterial
	case "JsonWebKey2020":
		var publicKeyJwk interface{}
		err := json.Unmarshal([]byte(protoVerificationMethod.VerificationMaterial), &publicKeyJwk)
		if err != nil {
			println("Invalid verification material !!!")
			panic(err)
		}
		verificationMethod.PublicKeyJwk = publicKeyJwk
	}

	return verificationMethod
}

func NewService(protoService *did.Service) *Service {
	return &Service{
		Id:              protoService.Id,
		Type:            protoService.ServiceType,
		ServiceEndpoint: protoService.ServiceEndpoint,
	}
}

func NewAssertionMethod(protoAssertionMethod string) *AssertionMethod {
	// Check if the string starts with "did:cheqd"
	if strings.HasPrefix(protoAssertionMethod, "did:cheqd") {
		return &AssertionMethod{
			Id:                  &protoAssertionMethod,
			AssertionMethodJSON: nil,
		}
	} else {
		// Try to parse it as VerificationMethod
		var verMethodString string
		err := json.Unmarshal([]byte(protoAssertionMethod), &verMethodString)
		// If parsing failed, return nil
		if err != nil {
			return &AssertionMethod{
				Id:                  nil,
				AssertionMethodJSON: nil,
			}
		}

		var verMethod VerificationMethod
		err = json.Unmarshal([]byte(verMethodString), &verMethod)
		// If parsing failed, return nil
		if err != nil {
			return &AssertionMethod{
				Id:                  nil,
				AssertionMethodJSON: nil,
			}
		}

		// Successfully parsed as VerificationMethod
		return &AssertionMethod{
			Id:                  nil,
			AssertionMethodJSON: &verMethod,
		}
	}
}

func (e *DidDoc) AddContext(newProtocol string) { e.Context = AddElemToSet(e.Context, newProtocol) }
func (e *DidDoc) RemoveContext()                { e.Context = nil }
func (e *DidDoc) GetBytes() []byte              { return []byte{} }

func (e *Service) AddContext(newProtocol string) { e.Context = AddElemToSet(e.Context, newProtocol) }
func (e *Service) RemoveContext()                { e.Context = nil }
func (e *Service) GetBytes() []byte              { return []byte{} }

func (e *VerificationMethod) AddContext(newProtocol string) {
	e.Context = AddElemToSet(e.Context, newProtocol)
}
func (e *VerificationMethod) RemoveContext()   { e.Context = nil }
func (e *VerificationMethod) GetBytes() []byte { return []byte{} }

func (d DidDoc) GetServiceByName(serviceId string) (string, error) {
	for _, s := range d.Service {
		_url, err := url.Parse(s.Id)
		if err != nil {
			return "", err
		}
		if _url.Fragment == serviceId {
			return s.ServiceEndpoint[0], nil
		}
	}
	return "", nil
}
