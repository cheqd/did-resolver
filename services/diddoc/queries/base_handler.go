package queries

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type BaseQueryHandlerI interface {
	SetNext(c services.ResolverContext, next BaseQueryHandlerI) error
	// ToDo too many parameters, need to increase
	Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error)
	Continue(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error)
}

type BaseQueryHandler struct {
	IsDereferencing bool
	next            BaseQueryHandlerI
}

func (b *BaseQueryHandler) SetNext(c services.ResolverContext, next BaseQueryHandlerI) error {
	// All the query handlers are dereferencing by default
	b.IsDereferencing = true
	if next == nil {
		return types.NewInternalError("next handler is nil", types.DIDJSONLD, nil, b.IsDereferencing)
	}
	b.next = next
	return nil
}

func (b *BaseQueryHandler) Continue(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	if b.next == nil {
		return nil, types.NewInternalError("next handler is nil", types.DIDJSONLD, nil, b.IsDereferencing)
	}
	return b.next.Handle(c, service, response)
}
