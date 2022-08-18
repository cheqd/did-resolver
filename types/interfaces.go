package types

type ContentStreamI interface {
	AddContext(newProtocol string)
	RemoveContext()
	GetBytes() []byte
}

type ResolutionResultI interface {
	GetStatus() int
	GetContentType() string
	GetBytes() []byte
}
