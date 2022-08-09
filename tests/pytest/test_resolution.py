import re

import pytest
import requests

from helpers import run, TESTNET_DID, MAINNET_DID, TESTNET_FRAGMENT, MAINNET_FRAGMENT, \
    FAKE_TESTNET_DID, FAKE_MAINNET_DID, FAKE_TESTNET_FRAGMENT, FAKE_MAINNET_FRAGMENT, RESOLVER_URL, PATH, \
    LDJSON, DIDJSON, DIDLDJSON, HTML, FAKE_TESTNET_RESOURCE, TESTNET_RESOURCE_METADATA, TESTNET_RESOURCE_NAME, JSON, \
    TESTNET_RESOURCE, RESOURCE_DATA, TESTNET_RESOURCE_LIST, TESTNET_RESOURCE_LIST_REDIRECT


@pytest.mark.parametrize(
    "did_url, expected_output",
    [
        (TESTNET_DID,
         fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\":\"{TESTNET_DID}\"(.*?)didDocumentMetadata(.*?){TESTNET_RESOURCE_NAME}"),
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

        (TESTNET_RESOURCE_METADATA, fr"\"contentStream\":(.*?)collectionId(.*?),\"contentMetadata\":(.*?),"
                                    r"\"dereferencingMetadata(.*?)"),
        (TESTNET_RESOURCE_LIST, fr"\"contentStream\":\[(.*?)collectionId(.*?),\"contentMetadata\":(.*?),"
                                r"\"dereferencingMetadata(.*?)"),
        (TESTNET_RESOURCE, RESOURCE_DATA),
        (FAKE_TESTNET_RESOURCE, r"\"contentStream\":null,\"contentMetadata\":\[\],"
                                r"\"dereferencingMetadata(.*?)\"error\":\"notFound\""),
    ]
)
def test_resolution(did_url, expected_output):
    run("curl", RESOLVER_URL + PATH + did_url.replace("#", "%23"), expected_output)


@pytest.mark.parametrize(
    "accept, expected_header, has_context, expected_status_code, expected_body",
    [
        (LDJSON, LDJSON, True, 200,
         r"(.*?)didResolutionMetadata(.*?)application/ld\+json"
         r"(.*?)didDocument(.*?)@context(.*?)didDocumentMetadata"),
        (DIDLDJSON, DIDLDJSON, True, 200,
         "(.*?)didResolutionMetadata(.*?)application/did\+ld\+json"
         "(.*?)didDocument(.*?)@context(.*?)didDocumentMetadata"),
        ("*/*", DIDLDJSON, True, 200,
         "(.*?)didResolutionMetadata(.*?)application/did\+ld\+json"
         "(.*?)didDocument(.*?)@context(.*?)didDocumentMetadata"),
        (DIDJSON, DIDJSON, False, 200,
         r"(.*?)didResolutionMetadata(.*?)application/did\+json"
         r"(.*?)didDocument(.*?)(?!`@context`)(.*?)didDocumentMetadata"),
        (HTML + ",application/xhtml+xml", JSON, False, 406,
         "(.*?)didResolutionMetadata(.*?)\"error\":\"representationNotSupported\""
         "(.*?)\"didDocument\":null,\"didDocumentMetadata\":\[\]"),
    ]
)
def test_resolution_content_type(accept, expected_header, expected_body, has_context, expected_status_code):
    url = RESOLVER_URL + PATH + TESTNET_DID
    header = {"Accept": accept} if accept else {}

    r = requests.get(url, headers=header)

    assert r.headers["Content-Type"] == expected_header
    assert r.status_code == expected_status_code
    assert re.match(expected_body, r.text)
    if has_context:
        assert re.findall(r"context", r.text)
    else:
        assert not re.findall(r"context", r.text)


dereferencing_content_type_test_set = [
    (LDJSON, LDJSON, True, 200,
     r"(.*?)contentStream(.*?)@context(.*?)contentMetadata"
     r"(.*?)dereferencingMetadata(.*?)application/ld\+json"),
    (DIDLDJSON, DIDLDJSON, True, 200,
     "(.*?)contentStream(.*?)@context(.*?)contentMetadata"
     "(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"),
    ("*/*", DIDLDJSON, True, 200,
     "(.*?)contentStream(.*?)@context(.*?)contentMetadata"
     "(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"),
    (DIDJSON, DIDJSON, False, 200,
     r"(.*?)contentStream(.*?)contentMetadata"
     r"(.*?)dereferencingMetadata(.*?)application/did\+json"),
    (HTML + ",application/xhtml+xml", JSON, False, 406,
     "(.*?)\"contentStream\":null,\"contentMetadata\":\[\],"
     "\"dereferencingMetadata(.*?)\"error\":\"representationNotSupported\""),
]


@pytest.mark.parametrize(
    "accept, expected_header, has_context, expected_status_code, expected_body",
    dereferencing_content_type_test_set
)
def test_dereferencing_content_type_fragment(accept, expected_header, expected_body, has_context, expected_status_code):
    url = RESOLVER_URL + PATH + TESTNET_FRAGMENT.replace("#", "%23")
    header = {"Accept": accept} if accept else {}

    r = requests.get(url, headers=header)

    assert r.headers["Content-Type"] == expected_header
    assert r.status_code == expected_status_code
    assert re.match(expected_body, r.text)
    if has_context:
        assert re.findall(r"context", r.text)
    else:
        assert not re.findall(r"context", r.text)


@pytest.mark.parametrize(
    "accept, expected_header, has_context, expected_status_code, expected_body",
    dereferencing_content_type_test_set
)
def test_dereferencing_content_type_resource_metadata(accept, expected_header, expected_body, has_context,
                                                      expected_status_code):
    url = RESOLVER_URL + PATH + TESTNET_RESOURCE_METADATA
    header = {"Accept": accept} if accept else {}

    r = requests.get(url, headers=header)

    assert r.headers["Content-Type"] == expected_header
    assert re.match(expected_body, r.text)
    if has_context:
        assert re.findall(r"context", r.text)
    else:
        assert not re.findall(r"context", r.text)


@pytest.mark.parametrize(
    "accept, expected_header, expected_status_code",
    [(LDJSON, JSON, 200), ]
)
def test_dereferencing_content_type_resource(accept, expected_header, expected_status_code):
    url = RESOLVER_URL + PATH + TESTNET_RESOURCE
    header = {"Accept": accept} if accept else {}
    r = requests.get(url, headers=header)
    assert r.headers["Content-Type"] == expected_header

@pytest.mark.parametrize(
    "accept, expected_header, expected_status_code, expected_body",
    [(LDJSON, LDJSON, 301,
      r"(.*?)\"contentStream\":\[(.*?)collectionId(.*?),\"contentMetadata\":(.*?),\"dereferencingMetadata(.*?)"),]
)
def test_dereferencing_content_type_resource_redirect(accept, expected_header, expected_status_code, expected_body):
    url = RESOLVER_URL + PATH + TESTNET_RESOURCE_LIST_REDIRECT
    header = {"Accept": accept} if accept else {}
    r = requests.get(url, headers=header)
    assert r.headers["Content-Type"] == expected_header
    assert re.match(expected_body, r.text)


@pytest.mark.parametrize(
    "did_url, expected_status_code",
    [
        (TESTNET_DID, 200),
        (TESTNET_FRAGMENT, 200),
        (TESTNET_RESOURCE_METADATA, 200),
        (FAKE_TESTNET_DID, 404),
        (FAKE_TESTNET_FRAGMENT, 404),
        (FAKE_TESTNET_RESOURCE, 404),
        ("did:wrong_method:MTMxDQKMTMxDQKMT", 406),
        (TESTNET_DID + "/", 406),
        (TESTNET_DID + "invalidDID", 400),
    ]
)
def test_resolution_status_code(did_url, expected_status_code):
    url = RESOLVER_URL + PATH + did_url.replace("#", "%23")
    r = requests.get(url)

    assert r.status_code == expected_status_code
