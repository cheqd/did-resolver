package types

import (
	resource "github.com/cheqd/cheqd-node/x/resource/types"
)

type DereferencedResource struct {
	Context []string
	resource.ResourceHeader
	Data []byte
}
