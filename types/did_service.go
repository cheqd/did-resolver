package types

// Implements ResolutionResult interface
type ServiceResult struct {
	endpoint string
}

// Interface implementation
func (s ServiceResult) GetContentType() string {
	return ""
}

func (s ServiceResult) GetBytes() []byte {
	return []byte(s.endpoint)
}

func (s ServiceResult) GetServiceEndpoint() string {
	return s.endpoint
}

func (r ServiceResult) IsRedirect() bool {
	return true
}

// end of Interface implementation

func NewServiceResult(endpoint string) *ServiceResult {
	return &ServiceResult{endpoint: endpoint}
}
