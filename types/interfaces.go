package types

type ContentStreamI interface {
	AddContext(newProtocol string)
	RemoveContext()
	GetBytes() []byte
}

type ResolutionResultI interface {
	GetContentType() string
	GetBytes() []byte
}
