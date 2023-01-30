import re

import pytest
import requests

from helpers import run, TESTNET_DID, MAINNET_DID, TESTNET_FRAGMENT, MAINNET_FRAGMENT, \
    FAKE_TESTNET_DID, FAKE_MAINNET_DID, FAKE_TESTNET_FRAGMENT, FAKE_MAINNET_FRAGMENT, RESOLVER_URL, PATH, \
    LDJSON, DIDJSON, DIDLDJSON, HTML, FAKE_TESTNET_RESOURCE, TESTNET_RESOURCE_METADATA, TESTNET_RESOURCE_NAME, JSON, \
    TESTNET_RESOURCE, RESOURCE_DATA, TESTNET_RESOURCE_LIST, INDY_TESTNET_DID, MIGRATED_INDY_TESTNET_DID, \
    TESTNET_DID_VERSION, TESTNET_DID_VERSION_ID, FAKE_TESTNET_VERSION, TESTNET_DID_VERSIONS, FAKE_TESTNET_DID_VERSIONS


@pytest.mark.parametrize(
    "did_url, expected_output",
    [
        (TESTNET_DID,
         fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\": \"{TESTNET_DID}\"(.*?)didDocumentMetadata(.*?){TESTNET_RESOURCE_NAME}"),
        
        # mainnet DID currently use another format of DID, when mainnet network will be same like testnet network you can run this test too.
        # (MAINNET_DID, fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\": \"{MAINNET_DID}\"(.*?)didDocumentMetadata"),
        
        (FAKE_TESTNET_DID, r"\"didResolutionMetadata(.*?)\"error\": \"notFound\"(.*?)"
                           r"didDocument\": null,(.*?)\"didDocumentMetadata\": \{\}"),
        (FAKE_MAINNET_DID, r"\"didResolutionMetadata(.*?)\"error\": \"notFound\"(.*?)"
                           r"didDocument\": null,(.*?)\"didDocumentMetadata\": \{\}"),
        ("did:wrong_method:MTMxDQKMTMxDQKMT", r"\"didResolutionMetadata(.*?)\"error\": \"methodNotSupported\"(.*?)"
                                              r"didDocument\": null,(.*?)\"didDocumentMetadata\": \{\}"),

        (TESTNET_FRAGMENT, r"(.*?)dereferencingMetadata\"(.*?)"
                           fr"\"contentStream\":(.*?)\"id\": \"{TESTNET_FRAGMENT}\"(.*?)contentMetadata"),
        
        # mainnet DID currently use another format of DID, when mainnet network will be same like testnet network you can run this test too.
        # (MAINNET_FRAGMENT, r"(.*?)dereferencingMetadata\"(.*?)"
        #                    fr"\"contentStream\":(.*?)\"id\": \"{MAINNET_FRAGMENT}\"(.*?)contentMetadata"),
        
        (FAKE_TESTNET_FRAGMENT, r"\"dereferencingMetadata(.*?)\"error\": \"notFound\"(.*?)"
                                r"\"contentStream\": null,(.*?)\"contentMetadata\": \{\}"),
        (FAKE_MAINNET_FRAGMENT, r"\"dereferencingMetadata(.*?)\"error\": \"notFound\"(.*?)"
                                r"\"contentStream\": null,(.*?)\"contentMetadata\": \{\}"),

        (TESTNET_RESOURCE_METADATA, r"\"dereferencingMetadata(.*?)\"contentStream\":(.*?)linkedResourceMetadata(.*?)"
                                    "resourceCollectionId(.*?)\"contentMetadata\":(.*?)"),
        (TESTNET_RESOURCE_LIST, r"\"dereferencingMetadata(.*?)\"contentStream\":(.*?)linkedResourceMetadata(.*?)"
                                "resourceCollectionId(.*?)\"contentMetadata\":(.*?)"),
        (TESTNET_RESOURCE, RESOURCE_DATA),
        (FAKE_TESTNET_RESOURCE, r"\"dereferencingMetadata(.*?)\"error\": \"notFound\"(.*?)"
                                r"\"contentStream\": null,(.*?)\"contentMetadata\": \{\}"),
        (INDY_TESTNET_DID, fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\": \"{MIGRATED_INDY_TESTNET_DID}\""),
        (MIGRATED_INDY_TESTNET_DID, fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\": \"{MIGRATED_INDY_TESTNET_DID}\""),
        (TESTNET_DID_VERSION, 
            fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\": \"{TESTNET_DID}\"(.*?)didDocumentMetadata(.*?)\"versionId\": \"{TESTNET_DID_VERSION_ID}\""),
        (FAKE_TESTNET_VERSION, r"\"didResolutionMetadata(.*?)\"error\": \"notFound\"(.*?)"
                           r"didDocument\": null,(.*?)\"didDocumentMetadata\": \{\}"),
        (TESTNET_DID_VERSIONS, r"\"dereferencingMetadata(.*?)\"contentStream\":(.*?)\"contentMetadata\":(.*?)"),
        (FAKE_TESTNET_DID_VERSIONS, r"\"didResolutionMetadata(.*?)\"error\": \"notFound\"(.*?)"
                           r"didDocument\": null,(.*?)\"didDocumentMetadata\": \{\}"),
    ]
)
def test_resolution(did_url, expected_output):
    run("curl", RESOLVER_URL + PATH + did_url.replace("#", "%23"), expected_output)


@pytest.mark.parametrize(
    "accept, expected_header, has_context, expected_status_code, expected_body",
    [
        (LDJSON, DIDLDJSON, True, 200,
         r"(.*?)context(.*?)didResolutionMetadata"),
        (DIDLDJSON, DIDLDJSON, True, 200,
         "(.*?)didResolutionMetadata(.*?)application/did\+ld\+json"
         "(.*?)didDocument(.*?)@context(.*?)didDocumentMetadata(.*?)"),
        ("*/*", DIDLDJSON, True, 200,
         "(.*?)didResolutionMetadata(.*?)application/did\+ld\+json"
         "(.*?)didDocument(.*?)@context(.*?)didDocumentMetadata(.*?)"),
        (DIDJSON, DIDJSON, False, 200,
         r"(.*?)didResolutionMetadata(.*?)application/did\+json"
         r"(.*?)didDocument(.*?)(?!`@context`)(.*?)didDocumentMetadata(.*?)"),
        (HTML + ",application/xhtml+xml", JSON, False, 406,
         "(.*?)didResolutionMetadata(.*?)\"error\": \"representationNotSupported\""
         "(.*?)\"didDocument\": null,(.*?)\"didDocumentMetadata\": \{\}"),
    ]
)
def test_resolution_content_type(accept, expected_header, expected_body, has_context, expected_status_code):
    url = RESOLVER_URL + PATH + TESTNET_DID
    header = {"Accept": accept} if accept else {}

    r = requests.get(url, headers=header)
    print(r.text.replace("\n", "\\n"))
    assert r.headers["Content-Type"] == expected_header
    assert r.status_code == expected_status_code
    assert re.match(expected_body, r.text.replace("\n", "\\n").replace("\n", "\\n"))
    if has_context:
        assert re.findall(r"context", r.text.replace("\n", "\\n").replace("\n", "\\n"))
    else:
        assert not re.findall(r"context", r.text.replace("\n", "\\n").replace("\n", "\\n"))


secondary_dereferencing_content_type_test_set = [
    (LDJSON, DIDLDJSON, True, 200,
     r"(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"
     r"(.*?)contentStream(.*?)@context(.*?)contentMetadata"),
    (DIDLDJSON, DIDLDJSON, True, 200,
     "(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"
     "(.*?)contentStream(.*?)@context(.*?)contentMetadata"),
    ("*/*", DIDLDJSON, True, 200,
     "(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"
     "(.*?)contentStream(.*?)@context(.*?)contentMetadata"),
    (DIDJSON, DIDJSON, False, 200,
     r"(.*?)dereferencingMetadata(.*?)application/did\+json"
     r"(.*?)contentStream(.*?)contentMetadata"),
    (HTML, JSON, False, 406,
     "(.*?)dereferencingMetadata(.*?)\"error\": \"representationNotSupported\""
     "(.*?)\"contentStream\": null,(.*?)\"contentMetadata\": \{\}"),
]


@pytest.mark.parametrize(
    "accept, expected_header, has_context, expected_status_code, expected_body",
    secondary_dereferencing_content_type_test_set
)
def test_dereferencing_content_type_fragment(accept, expected_header, expected_body, has_context, expected_status_code):
    url = RESOLVER_URL + PATH + TESTNET_FRAGMENT.replace("#", "%23")
    header = {"Accept": accept} if accept else {}

    r = requests.get(url, headers=header)

    assert r.headers["Content-Type"] == expected_header
    assert r.status_code == expected_status_code
    assert re.match(expected_body, r.text.replace("\n", "\\n"))
    if has_context:
        assert re.findall(r"context", r.text.replace("\n", "\\n"))
    else:
        assert not re.findall(r"context", r.text.replace("\n", "\\n"))


primary_dereferencing_content_type_test_set = [
    (LDJSON, DIDLDJSON, True, 200,
     r"(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"
     r"(.*?)contentStream(.*?)contentMetadata"),
    (DIDLDJSON, DIDLDJSON, True, 200,
     "(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"
     "(.*?)contentStream(.*?)contentMetadata"),
    ("*/*", DIDLDJSON, True, 200,
     "(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"
     "(.*?)contentStream(.*?)contentMetadata"),
    (DIDJSON, DIDJSON, False, 200,
     r"(.*?)dereferencingMetadata(.*?)application/did\+json"
     r"(.*?)contentStream(.*?)contentMetadata"),
    (HTML, JSON, False, 406,
     "(.*?)dereferencingMetadata(.*?)\"error\": \"representationNotSupported\""
     "(.*?)\"contentStream\": null,(.*?)\"contentMetadata\": \{\}"),
]


@pytest.mark.parametrize(
    "accept, expected_header, has_context, expected_status_code, expected_body",
    primary_dereferencing_content_type_test_set
)
def test_dereferencing_content_type_resource_metadata(accept, expected_header, expected_body, has_context,
                                                      expected_status_code):
    url = RESOLVER_URL + PATH + TESTNET_RESOURCE_METADATA
    header = {"Accept": accept} if accept else {}

    r = requests.get(url, headers=header)
    print(r.text)
    assert r.headers["Content-Type"] == expected_header
    assert re.match(expected_body, r.text.replace("\n", "\\n"))
    if has_context:
        assert re.findall(r"context", r.text.replace("\n", "\\n"))
    else:
        assert not re.findall(r"context", r.text.replace("\n", "\\n"))


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
    "did_url, expected_status_code",
    [
        (TESTNET_DID, 200),
        (TESTNET_FRAGMENT, 200),
        (TESTNET_RESOURCE_METADATA, 200),
        (FAKE_TESTNET_DID, 404),
        (FAKE_TESTNET_FRAGMENT, 404),
        (FAKE_TESTNET_RESOURCE, 404),
        ("did:wrong_method:MTMxDQKMTMxDQKMT", 406),
        (TESTNET_DID + "/", 400),
        (TESTNET_DID + "invalidDID", 400),
    ]
)
def test_resolution_status_code(did_url, expected_status_code):
    url = RESOLVER_URL + PATH + did_url.replace("#", "%23")
    r = requests.get(url)

    assert r.status_code == expected_status_code
