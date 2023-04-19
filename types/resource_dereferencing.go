package types

type ResourceDereferencing struct {
	Context               string                     `json:"@context,omitempty" example:"https://w3id.org/did-resolution/v1"`
	DereferencingMetadata DereferencingMetadata      `json:"dereferencingMetadata"`
	ContentStream         ContentStreamI             `json:"contentStream"`
	Metadata              ResolutionResourceMetadata `json:"contentMetadata"`
}

func NewResourceDereferencingFromContent(did string, contentType ContentType, contentStream ContentStreamI) *ResourceDereferencing {
	dereferenceMetadata := NewDereferencingMetadata(did, contentType, "")

	var context string
	if contentType == DIDJSONLD || contentType == JSONLD {
		context = ResolutionSchemaJSONLD
	}

	return &ResourceDereferencing{Context: context, ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}
}

// Interface implementation

func (d ResourceDereferencing) GetContentType() string {
	return string(d.DereferencingMetadata.ContentType)
}

func (d ResourceDereferencing) GetBytes() []byte {
	if d.ContentStream == nil {
		return []byte{}
	}
	return d.ContentStream.GetBytes()
}

func (r ResourceDereferencing) IsRedirect() bool {
	return false
}
