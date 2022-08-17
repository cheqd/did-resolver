package types

type DereferencingOption ResolutionOption

type DereferencingMetadata ResolutionMetadata

type DidDereferencing struct {
	DereferencingMetadata DereferencingMetadata    `json:"dereferencingMetadata,omitempty"`
	ContentStream         ContentStreamI           `json:"contentStream,omitempty"`
	Metadata              ResolutionDidDocMetadata `json:"contentMetadata,omitempty"`
}

func NewDereferencingMetadata(did string, contentType ContentType, resolutionError ErrorType) DereferencingMetadata {
	return DereferencingMetadata(NewResolutionMetadata(did, contentType, resolutionError))
}

func (d DidDereferencing) GetStatus() int {
	return d.DereferencingMetadata.ResolutionError.GetStatusCode()
}

func (d DidDereferencing) GetContentType() string {
	return string(d.DereferencingMetadata.ContentType)
}

func (d DidDereferencing) GetBytes() []byte {
	if d.ContentStream == nil {
		return []byte{}
	}
	return d.ContentStream.GetBytes()
}
