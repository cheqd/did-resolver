import pytest
from helpers import run, TESTNET_DID, MAINNET_DID, TESTNET_FRAGMENT, MAINNET_FRAGMENT, \
    FAKE_TESTNET_DID, FAKE_MAINNET_DID, FAKE_TESTNET_FRAGMENT, FAKE_MAINNET_FRAGMENT, RESOLVER_URL, PATH, \
    FAKE_TESTNET_RESOURCE


@pytest.mark.parametrize(
    "did_url, expected_output",
    [
        (TESTNET_DID, fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\":\"{TESTNET_DID}\"(.*?)didDocumentMetadata"),
        (MAINNET_DID, fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\":\"{MAINNET_DID}\"(.*?)didDocumentMetadata"),
        (FAKE_TESTNET_DID, r"\"didResolutionMetadata(.*?)\"error\":\"notFound\"(.*?)"
                           r"didDocument\":null,\"didDocumentMetadata\":\[\]"),
        (FAKE_MAINNET_DID, r"\"didResolutionMetadata(.*?)\"error\":\"notFound\"(.*?)"
                           r"didDocument\":null,\"didDocumentMetadata\":\[\]"),
        ("did:wrong_method:MTMxDQKMTMxDQKMT", r"\"didResolutionMetadata(.*?)\"error\":\"methodNotSupported\"(.*?)"
                                              r"didDocument\":null,\"didDocumentMetadata\":\[\]"),

        (TESTNET_FRAGMENT, fr"\"contentStream\":(.*?)\"id\":\"{TESTNET_FRAGMENT}\"(.*?)contentMetadata"
                           r"(.*?)dereferencingMetadata\""),
        (MAINNET_FRAGMENT, fr"\"contentStream\":(.*?)\"id\":\"{MAINNET_FRAGMENT}\"(.*?)contentMetadata"
                           r"(.*?)dereferencingMetadata\""),
        (FAKE_TESTNET_FRAGMENT, r"\"contentStream\":null,\"contentMetadata\":\[\],"
                                r"\"dereferencingMetadata(.*?)\"error\":\"notFound\""),
        (FAKE_MAINNET_FRAGMENT, r"\"contentStream\":null,\"contentMetadata\":\[\],"
                                r"\"dereferencingMetadata(.*?)\"error\":\"notFound\""),

        (FAKE_TESTNET_RESOURCE, r"\"contentStream\":null,\"contentMetadata\":\[\],"
                                r"\"dereferencingMetadata(.*?)\"error\":\"notFound\""),
    ]
)
def test_resolution(did_url, expected_output):
    run("curl", RESOLVER_URL + PATH + did_url.replace("#", "%23"), expected_output)
