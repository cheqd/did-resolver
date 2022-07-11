import re

import pytest
import requests

from helpers import run, TESTNET_DID, MAINNET_DID, TESTNET_FRAGMENT, MAINNET_FRAGMENT, \
    FAKE_TESTNET_DID, FAKE_MAINNET_DID, FAKE_TESTNET_FRAGMENT, FAKE_MAINNET_FRAGMENT, RESOLVER_URL, PATH, \
    LDJSON, DIDJSON, DIDLDJSON, HTML


@pytest.mark.parametrize(
    "did_url, expected_output",
    [
        (TESTNET_DID, fr"didDocument(.*?)\"id\":\"{TESTNET_DID}\"(.*?)didDocumentMetadata"
                      r"(.*?)didResolutionMetadata"),
        (MAINNET_DID, fr"didDocument(.*?)\"id\":\"{MAINNET_DID}\"(.*?)didDocumentMetadata"
                      r"(.*?)didResolutionMetadata"),
        (FAKE_TESTNET_DID, r"didDocument\":null,\"didDocumentMetadata\":\[\],"
                           r"\"didResolutionMetadata(.*?)\"error\":\"notFound\""),
        (FAKE_MAINNET_DID, r"didDocument\":null,\"didDocumentMetadata\":\[\],"
                           r"\"didResolutionMetadata(.*?)\"error\":\"notFound\""),
        ("did:wrong_method:MTMxDQKMTMxDQKMT", r"didDocument\":null,\"didDocumentMetadata\":\[\],"
                                              r"\"didResolutionMetadata(.*?)\"error\":\"methodNotSupported\""),

        (TESTNET_FRAGMENT, fr"\"contentStream\":(.*?)\"id\":\"{TESTNET_FRAGMENT}\"(.*?)contentMetadata"
                           r"(.*?)dereferencingMetadata\""),
        (MAINNET_FRAGMENT, fr"\"contentStream\":(.*?)\"id\":\"{MAINNET_FRAGMENT}\"(.*?)contentMetadata"
                           r"(.*?)dereferencingMetadata\""),
        (FAKE_TESTNET_FRAGMENT, r"\"contentStream\":null,\"contentMetadata\":\[\],"
                                r"\"dereferencingMetadata(.*?)\"error\":\"FragmentNotFound\""),
        (FAKE_MAINNET_FRAGMENT, r"\"contentStream\":null,\"contentMetadata\":\[\],"
                                r"\"dereferencingMetadata(.*?)\"error\":\"FragmentNotFound\""),
    ]
)
def test_resolution(did_url, expected_output):
    run("curl", RESOLVER_URL + PATH + did_url.replace("#", "%23"), expected_output)


@pytest.mark.parametrize(
    "accept, expected_header, expected_body",
    [
        (LDJSON, LDJSON, r"(.*?)didDocument(.*?)@context(.*?)didDocumentMetadata"
                         r"(.*?)didResolutionMetadata(.*?)application/ld\+json"),
        (DIDLDJSON, DIDLDJSON, "(.*?)didDocument(.*?)@context(.*?)didDocumentMetadata"
                               "(.*?)didResolutionMetadata(.*?)application/did\+ld\+json"),
        ("", DIDLDJSON, "(.*?)didDocument(.*?)@context(.*?)didDocumentMetadata"
                        "(.*?)didResolutionMetadata(.*?)application/did\+ld\+json"),
        (DIDJSON, DIDJSON, r"(.*?)didDocument(.*?)(?!`@context`)(.*?)didDocumentMetadata"
                           r"(.*?)didResolutionMetadata(.*?)application/did\+json"),
        (HTML + ",application/xhtml+xml", HTML, fr"(.*?)didDocument(.*?)(?!`@context`)(.*?)didDocumentMetadata"
                                                fr"(.*?)didResolutionMetadata(.*?){HTML}"),
    ]
)
def test_resolution_content_type(accept, expected_header, expected_body):
    url = RESOLVER_URL + PATH + TESTNET_DID
    header = {"Accept": accept} if accept else {}

    r = requests.get(url, headers=header)

    assert r.headers["Content-Type"] == expected_header
    assert re.match(expected_body, r.text)
