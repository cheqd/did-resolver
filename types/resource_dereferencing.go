package types

type ResourceDereferencing struct {
	Context               string                `json:"@context,omitempty" example:"https://w3id.org/did-resolution/v1"`
	DereferencingMetadata DereferencingMetadata `json:"dereferencingMetadata"`
	ContentStream         ContentStreamI        `json:"contentStream"`
	Metadata              *DereferencedResource `json:"contentMetadata"`
}

func NewResourceDereferencingFromContent(did string, contentType ContentType, contentStream ContentStreamI) *ResourceDereferencing {
	dereferenceMetadata := NewDereferencingMetadata(did, contentType, "")

	var context string
	if contentType == DIDJSONLD || contentType == JSONLD {
		context = ResolutionSchemaJSONLD
	}

	return &ResourceDereferencing{Context: context, ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}
}

func NewResourceDereferencingFromDidDocMetadata(did string, contentType ContentType, didDocMetadata *ResolutionDidDocMetadata) *ResourceDereferencing {
	dereferenceMetadata := NewDereferencingMetadata(did, contentType, "")

	var context string
	if contentType == DIDJSONLD || contentType == JSONLD {
		context = ResolutionSchemaJSONLD
	}

	resources := didDocMetadata.Resources
	metadata := resources[0]
	return &ResourceDereferencing{Context: context, Metadata: &metadata, DereferencingMetadata: dereferenceMetadata}
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
