package services

import (
	"context"
	"crypto/tls"
	"errors"
	"strings"

	"google.golang.org/grpc/credentials"

	didTypes "github.com/cheqd/cheqd-node/x/did/types"
	didUtils "github.com/cheqd/cheqd-node/x/did/utils"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	DELIMITER = ":"
)

type LedgerServiceI interface {
	QueryDIDDoc(did string, version string) (*didTypes.DidDocWithMetadata, *types.IdentityError)
	QueryAllDidDocVersionsMetadata(did string) ([]*didTypes.Metadata, *types.IdentityError)
	QueryResource(collectionDid string, resourceId string) (*resource.ResourceWithMetadata, *types.IdentityError)
	QueryCollectionResources(did string) ([]*resource.Metadata, *types.IdentityError)
	GetNamespaces() []string
}

type LedgerService struct {
	ledgers map[string]types.Network // namespace -> endpoint with configs
}

func NewLedgerService() LedgerService {
	ls := LedgerService{}
	ls.ledgers = make(map[string]types.Network)
	return ls
}

// QueryDIDDoc godoc
//
//	@Summary		Resolve DID Document on did:cheqd
//	@Description	Fetch DID Document ("DIDDoc") from cheqd network
//	@Tags			DID Resolution
//	@Accept			application/did+ld+json,application/ld+json,application/did+json
//	@Produce		application/did+ld+json,application/ld+json,application/did+json
//	@Param			did			path		string	true	"Full DID with unique identifier"
//	@Param			service		query		string	false	"Service Type"
//	@Param			fragmentId	query		string	false	"#Fragment"
//	@Param			versionId	query		string	false	"Version"
//	@Success		200			{object}	types.DidResolution
//	@Failure		400			{object}	types.IdentityError
//	@Failure		404			{object}	types.IdentityError
//	@Failure		406			{object}	types.IdentityError
//	@Failure		500			{object}	types.IdentityError
//	@Router			/{did} [get]
func (ls LedgerService) QueryDIDDoc(did string, version string) (*didTypes.DidDocWithMetadata, *types.IdentityError) {
	method, namespace, _, _ := didUtils.TrySplitDID(did)
	serverAddr, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDIDError(did, types.JSON, nil, false)
	}

	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryDIDDoc: failed connection")
		return nil, types.NewInternalError(did, types.JSON, err, false)
	}

	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying did doc: %s", did)
	client := didTypes.NewQueryClient(conn)

	if version == "" {
		didDocResponse, err := client.DidDoc(context.Background(), &didTypes.QueryDidDocRequest{Id: did})
		if err != nil {
			return nil, types.NewNotFoundError(did, types.JSON, err, false)
		}

		return didDocResponse.Value, nil
	} else {
		didDocResponse, err := client.DidDocVersion(context.Background(), &didTypes.QueryDidDocVersionRequest{Id: did, Version: version})
		if err != nil {
			return nil, types.NewNotFoundError(did, types.JSON, err, false)
		}

		return didDocResponse.Value, nil
	}
}

// QueryDidDocVersionMetadata godoc
//
//	@Summary		Resolve DID Document Version on did:cheqd
//	@Description	Fetch specific all version of a DID Document ("DIDDoc") for a given DID
//	@Tags			DID Resolution
//	@Accept			application/did+ld+json,application/ld+json,application/did+json
//	@Produce		application/did+ld+json,application/ld+json,application/did+json
//	@Param			did			path		string	true	"Full DID with unique identifier"
//	@Param			versionId	path		string	true	"version of a DID document"
//	@Success		200			{object}	types.DidResolution
//	@Failure		400			{object}	types.IdentityError
//	@Failure		404			{object}	types.IdentityError
//	@Failure		406			{object}	types.IdentityError
//	@Failure		500			{object}	types.IdentityError
//	@Router			/{did}/version/{versionId} [get]
//
// QueryAllDidDocVersionsMetadata godoc
//
//	@Summary		Resolve DID Document Versions on did:cheqd
//	@Description	Fetch specific all versions of a DID Document ("DIDDoc") for a given DID
//	@Tags			DID Resolution
//	@Accept			application/did+ld+json,application/ld+json,application/did+json
//	@Produce		application/did+ld+json,application/ld+json,application/did+json
//	@Param			did	path		string	true	"Full DID with unique identifier"
//	@Success		200	{object}	types.DidDereferencing{contentStream=types.DereferencedDidVersionsList}
//	@Failure		400	{object}	types.IdentityError
//	@Failure		404	{object}	types.IdentityError
//	@Failure		406	{object}	types.IdentityError
//	@Failure		500	{object}	types.IdentityError
//	@Router			/{did}/versions [get]
func (ls LedgerService) QueryAllDidDocVersionsMetadata(did string) ([]*didTypes.Metadata, *types.IdentityError) {
	method, namespace, _, _ := didUtils.TrySplitDID(did)
	serverAddr, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDIDError(did, types.JSON, nil, false)
	}

	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryAllDidDocVersionsMetadata: failed connection")
		return nil, types.NewInternalError(did, types.JSON, err, false)
	}
	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying all did doc versions metadata: %s", did)
	client := didTypes.NewQueryClient(conn)

	response, err := client.AllDidDocVersionsMetadata(context.Background(), &didTypes.QueryAllDidDocVersionsMetadataRequest{Id: did})
	if err != nil {
		return nil, types.NewNotFoundError(did, types.JSON, err, false)
	}

	return response.Versions, nil
}

// QueryResource godoc
//
//	@Summary		Fetch specific Resource
//	@Description	Get specific Resource within a DID Resource Collection
//	@Tags			Resource Resolution
//	@Accept			*/*
//	@Produce		*/*
//	@Param			did			path		string	true	"Full DID with unique identifier"
//	@Param			resourceId	path		string	true	"Resource-specific unique-identifier"
//	@Success		200			{object}	[]byte
//	@Failure		400			{object}	types.IdentityError
//	@Failure		404			{object}	types.IdentityError
//	@Failure		406			{object}	types.IdentityError
//	@Failure		500			{object}	types.IdentityError
//	@Router			/{did}/resources/{resourceId} [get]
func (ls LedgerService) QueryResource(did string, resourceId string) (*resource.ResourceWithMetadata, *types.IdentityError) {
	method, namespace, collectionId, _ := didUtils.TrySplitDID(did)
	serverAddr, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDIDError(did, types.JSON, nil, true)
	}

	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryResource: failed connection")
		return nil, types.NewInternalError(did, types.JSON, err, true)
	}

	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying did resource: %s, %s", collectionId, resourceId)

	client := resource.NewQueryClient(conn)
	resourceResponse, err := client.Resource(context.Background(), &resource.QueryResourceRequest{CollectionId: collectionId, Id: resourceId})
	if err != nil {
		log.Info().Msgf("Resource not found %s", err.Error())
		return nil, types.NewNotFoundError(did, types.JSON, err, true)
	}

	return resourceResponse.Resource, nil
}

// QueryCollectionResources godoc
//
//	@Summary		Fetch metadata for all Resources
//	@Description	Get metadata for all Resources within a DID Resource Collection
//	@Tags			Resource Resolution
//	@Accept			application/did+ld+json,application/ld+json,application/did+json
//	@Produce		application/did+ld+json,application/ld+json,application/did+json
//	@Param			did	path		string	true	"Full DID with unique identifier"
//	@Success		200	{object}	types.DidDereferencing
//	@Failure		400	{object}	types.IdentityError
//	@Failure		404	{object}	types.IdentityError
//	@Failure		406	{object}	types.IdentityError
//	@Failure		500	{object}	types.IdentityError
//	@Router			/{did}/metadata [get]
func (ls LedgerService) QueryCollectionResources(did string) ([]*resource.Metadata, *types.IdentityError) {
	method, namespace, collectionId, _ := didUtils.TrySplitDID(did)
	serverAddr, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDIDError(did, types.JSON, nil, false)
	}

	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryResource: failed connection")
		return nil, types.NewInternalError(did, types.JSON, err, false)
	}

	log.Info().Msgf("Querying did resources: %s", did)

	client := resource.NewQueryClient(conn)
	resourceResponse, err := client.CollectionResources(context.Background(), &resource.QueryCollectionResourcesRequest{CollectionId: collectionId})
	if err != nil {
		return nil, types.NewNotFoundError(did, types.JSON, err, false)
	}

	return resourceResponse.Resources, nil
}

func (ls *LedgerService) RegisterLedger(method string, endpoint types.Network) error {
	if endpoint.Namespace == "" || method == "" {
		err := errors.New("namespace and method cannot be empty")
		log.Error().Err(err).Msg("RegisterLedger: failed")
		return err
	}

	if endpoint.Endpoint == "" {
		return errors.New("ledger node url cannot be empty")
	}

	ls.ledgers[method+DELIMITER+endpoint.Namespace] = endpoint
	return nil
}

func (ls LedgerService) openGRPCConnection(endpoint types.Network) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	if endpoint.UseTls {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), endpoint.Timeout)
	defer cancel()

	conn, err = grpc.DialContext(ctx, endpoint.Endpoint, opts...)

	if err != nil {
		log.Error().Err(err).Msgf("openGRPCConnection: context failed")
		return nil, err
	}

	log.Info().Msg("openGRPCConnection: opened")
	return conn, nil
}

func mustCloseGRPCConnection(conn *grpc.ClientConn) {
	if conn == nil {
		return
	}
	err := conn.Close()
	if err != nil {
		log.Panic().Err(err).Msg("QueryDIDDoc: failed to close connection")
		panic(err)
	}
}

func (ls LedgerService) GetNamespaces() []string {
	keys := make([]string, 0, len(ls.ledgers))
	for k := range ls.ledgers {
		namespace := strings.Split(k, DELIMITER)[1]
		keys = append(keys, namespace)
	}
	return keys
}
