import re

import pytest
import requests
import helpers


@pytest.mark.parametrize(
    "did_url, expected_output",
    [
        (helpers.TESTNET_DID,
         fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\": \"{helpers.TESTNET_DID}\"(.*?)didDocumentMetadata(.*?){helpers.TESTNET_RESOURCE_NAME}"),

        # mainnet DID currently use another format of DID, when mainnet network will be same like testnet network you can run this test too.
        # (MAINNET_DID, fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\": \"{MAINNET_DID}\"(.*?)didDocumentMetadata"),

        (helpers.FAKE_TESTNET_DID, r"\"didResolutionMetadata(.*?)\"error\": \"notFound\"(.*?)"
         r"didDocument\": null,(.*?)\"didDocumentMetadata\": \{\}"),
        (helpers.FAKE_MAINNET_DID, r"\"didResolutionMetadata(.*?)\"error\": \"notFound\"(.*?)"
         r"didDocument\": null,(.*?)\"didDocumentMetadata\": \{\}"),
        ("did:wrong_method:MTMxDQKMTMxDQKMT", r"\"didResolutionMetadata(.*?)\"error\": \"methodNotSupported\"(.*?)"
                                              r"didDocument\": null,(.*?)\"didDocumentMetadata\": \{\}"),

        (helpers.TESTNET_FRAGMENT_1, r"(.*?)dereferencingMetadata\"(.*?)"
         fr"\"contentStream\":(.*?)\"id\": \"{helpers.TESTNET_FRAGMENT_1}\"(.*?)contentMetadata"),
        (helpers.TESTNET_FRAGMENT_2, r"(.*?)dereferencingMetadata\"(.*?)"
         fr"\"contentStream\":(.*?)\"id\": \"{helpers.MIGRATED_TESTNET_FRAGMENT_2}\"(.*?)contentMetadata"),

        # mainnet DID currently use another format of DID, when mainnet network will be same like testnet network you can run this test too.
        # (MAINNET_FRAGMENT_1, r"(.*?)dereferencingMetadata\"(.*?)"
        #                    fr"\"contentStream\":(.*?)\"id\": \"{MAINNET_FRAGMENT_1}\"(.*?)contentMetadata"),

        (helpers.FAKE_TESTNET_FRAGMENT_1, r"\"dereferencingMetadata(.*?)\"error\": \"notFound\"(.*?)"
         r"\"contentStream\": null,(.*?)\"contentMetadata\": \{\}"),
        (helpers.FAKE_TESTNET_FRAGMENT_2, r"\"dereferencingMetadata(.*?)\"error\": \"notFound\"(.*?)"
         r"\"contentStream\": null,(.*?)\"contentMetadata\": \{\}"),
        (helpers.FAKE_MAINNET_FRAGMENT_1, r"\"dereferencingMetadata(.*?)\"error\": \"notFound\"(.*?)"
         r"\"contentStream\": null,(.*?)\"contentMetadata\": \{\}"),
        (helpers.TESTNET_RESOURCE_METADATA_1, r"\"dereferencingMetadata(.*?)\"contentStream\":(.*?)linkedResourceMetadata(.*?)"
         "resourceCollectionId(.*?)\"contentMetadata\":(.*?)"),
        (helpers.TESTNET_RESOURCE_LIST, r"\"dereferencingMetadata(.*?)\"contentStream\":(.*?)linkedResourceMetadata(.*?)"
         "resourceCollectionId(.*?)\"contentMetadata\":(.*?)"),
        (helpers.TESTNET_RESOURCE_1, helpers.RESOURCE_DATA_1),
        (helpers.FAKE_TESTNET_RESOURCE, r"\"dereferencingMetadata(.*?)\"error\": \"notFound\"(.*?)"
         r"\"contentStream\": null,(.*?)\"contentMetadata\": \{\}"),
        (helpers.INDY_TESTNET_DID_1,
         fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\": \"{helpers.MIGRATED_INDY_TESTNET_DID_1}\""),
        (helpers.MIGRATED_INDY_TESTNET_DID_1,
         fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\": \"{helpers.MIGRATED_INDY_TESTNET_DID_1}\""),
        (helpers.INDY_TESTNET_DID_2,
         fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\": \"{helpers.MIGRATED_INDY_TESTNET_DID_2}\""),
        (helpers.MIGRATED_INDY_TESTNET_DID_2,
         fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\": \"{helpers.MIGRATED_INDY_TESTNET_DID_2}\""),
        (helpers.TESTNET_DID_VERSION,
            fr"didResolutionMetadata(.*?)didDocument(.*?)\"id\": \"{helpers.TESTNET_DID}\"(.*?)didDocumentMetadata(.*?)\"versionId\": \"{helpers.TESTNET_DID_VERSION_ID}\""),
        (helpers.FAKE_TESTNET_VERSION, r"\"didResolutionMetadata(.*?)\"error\": \"notFound\"(.*?)"
         r"didDocument\": null,(.*?)\"didDocumentMetadata\": \{\}"),
        (helpers.TESTNET_DID_VERSIONS,
         r"\"dereferencingMetadata(.*?)\"contentStream\":(.*?)\"contentMetadata\":(.*?)"),
        (helpers.FAKE_TESTNET_DID_VERSIONS, r"\"dereferencingMetadata(.*?)\"contentStream\":(.*?)\"contentMetadata\":(.*?)"),
        (helpers.TESTNET_DID_VERSION_METADATA, r"\"dereferencingMetadata(.*?)\"contentStream\":(.*?)linkedResourceMetadata(.*?)"
         "resourceCollectionId(.*?)\"contentMetadata\":(.*?)"),
    ]
)
def test_resolution(did_url, expected_output):
    helpers.run("curl -L", helpers.RESOLVER_URL + helpers.PATH +
                did_url.replace("#", "%23"), expected_output)


@pytest.mark.parametrize(
    "accept, expected_header, has_context, expected_status_code, expected_body",
    [
        (helpers.LDJSON, helpers.DIDLDJSON, True, 200,
         r"(.*?)context(.*?)didResolutionMetadata"),
        (helpers.DIDLDJSON, helpers.DIDLDJSON, True, 200,
         "(.*?)didResolutionMetadata(.*?)application/did\+ld\+json"
         "(.*?)didDocument(.*?)@context(.*?)didDocumentMetadata(.*?)"),
        ("*/*", helpers.DIDLDJSON, True, 200,
         "(.*?)didResolutionMetadata(.*?)application/did\+ld\+json"
         "(.*?)didDocument(.*?)@context(.*?)didDocumentMetadata(.*?)"),
        (helpers.DIDJSON, helpers.DIDJSON, False, 200,
         r"(.*?)didResolutionMetadata(.*?)application/did\+json"
         r"(.*?)didDocument(.*?)(?!`@context`)(.*?)didDocumentMetadata(.*?)"),
        (helpers.HTML + ",application/xhtml+xml", helpers.JSON, False, 406,
         "(.*?)didResolutionMetadata(.*?)\"error\": \"representationNotSupported\""
         "(.*?)\"didDocument\": null,(.*?)\"didDocumentMetadata\": \{\}"),
    ]
)
def test_resolution_content_type(accept, expected_header, expected_body, has_context, expected_status_code):
    url = helpers.RESOLVER_URL + helpers.PATH + helpers.TESTNET_DID
    header = {"Accept": accept} if accept else {}

    r = requests.get(url, headers=header)
    print(r.text.replace("\n", "\\n"))
    assert r.headers["Content-Type"] == expected_header
    assert r.status_code == expected_status_code
    assert re.match(expected_body, r.text.replace(
        "\n", "\\n").replace("\n", "\\n"))
    if has_context:
        assert re.findall(r"context", r.text.replace(
            "\n", "\\n").replace("\n", "\\n"))
    else:
        assert not re.findall(r"context", r.text.replace(
            "\n", "\\n").replace("\n", "\\n"))


secondary_dereferencing_content_type_test_set = [
    (helpers.LDJSON, helpers.DIDLDJSON, True, 200,
     r"(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"
     r"(.*?)contentStream(.*?)@context(.*?)contentMetadata"),
    (helpers.DIDLDJSON, helpers. DIDLDJSON, True, 200,
     "(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"
     "(.*?)contentStream(.*?)@context(.*?)contentMetadata"),
    ("*/*", helpers.DIDLDJSON, True, 200,
     "(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"
     "(.*?)contentStream(.*?)@context(.*?)contentMetadata"),
    (helpers.DIDJSON, helpers.DIDJSON, False, 200,
     r"(.*?)dereferencingMetadata(.*?)application/did\+json"
     r"(.*?)contentStream(.*?)contentMetadata"),
    (helpers.HTML, helpers.JSON, False, 406,
     "(.*?)dereferencingMetadata(.*?)\"error\": \"representationNotSupported\""
     "(.*?)\"contentStream\": null,(.*?)\"contentMetadata\": \{\}"),
]


@pytest.mark.parametrize(
    "accept, expected_header, has_context, expected_status_code, expected_body",
    secondary_dereferencing_content_type_test_set
)
def test_dereferencing_content_type_fragment(accept, expected_header, expected_body, has_context, expected_status_code):
    url = helpers.RESOLVER_URL + helpers.PATH + \
        helpers.TESTNET_FRAGMENT_1.replace("#", "%23")
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
    (helpers.LDJSON, helpers.DIDLDJSON, True, 200,
     r"(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"
     r"(.*?)contentStream(.*?)contentMetadata"),
    (helpers.DIDLDJSON, helpers.DIDLDJSON, True, 200,
     "(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"
     "(.*?)contentStream(.*?)contentMetadata"),
    ("*/*", helpers.DIDLDJSON, True, 200,
     "(.*?)dereferencingMetadata(.*?)application/did\+ld\+json"
     "(.*?)contentStream(.*?)contentMetadata"),
    (helpers.DIDJSON, helpers.DIDJSON, False, 200,
     r"(.*?)dereferencingMetadata(.*?)application/did\+json"
     r"(.*?)contentStream(.*?)contentMetadata"),
    (helpers.HTML, helpers.JSON, False, 406,
     "(.*?)dereferencingMetadata(.*?)\"error\": \"representationNotSupported\""
     "(.*?)\"contentStream\": null,(.*?)\"contentMetadata\": \{\}"),
]


@pytest.mark.parametrize(
    "accept, expected_header, has_context, expected_status_code, expected_body",
    primary_dereferencing_content_type_test_set
)
def test_dereferencing_content_type_resource_metadata(accept, expected_header, expected_body, has_context,
                                                      expected_status_code):
    url = helpers.RESOLVER_URL + helpers.PATH + helpers.TESTNET_RESOURCE_METADATA_1
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
    [(helpers.LDJSON, helpers.JSON, 200), ]
)
def test_dereferencing_content_type_resource(accept, expected_header, expected_status_code):
    url = helpers.RESOLVER_URL + helpers.PATH + helpers.TESTNET_RESOURCE_1
    header = {"Accept": accept} if accept else {}
    r = requests.get(url, headers=header)
    assert r.headers["Content-Type"] == expected_header


@pytest.mark.parametrize(
    "did_url, expected_status_code",
    [
        (helpers.TESTNET_DID, 200),
        (helpers.INDY_TESTNET_DID_1, 200),
        (helpers.INDY_TESTNET_DID_2, 200),
        (helpers.TESTNET_FRAGMENT_1, 200),
        (helpers.TESTNET_RESOURCE_2, 200),
        (helpers.TESTNET_RESOURCE_METADATA_1, 200),
        (helpers.TESTNET_RESOURCE_METADATA_2, 200),
        (helpers.FAKE_TESTNET_DID, 404),
        (helpers.FAKE_TESTNET_FRAGMENT_1, 404),
        (helpers.FAKE_TESTNET_FRAGMENT_2, 404),
        (helpers.FAKE_TESTNET_RESOURCE, 404),
        ("did:wrong_method:MTMxDQKMTMxDQKMT", 501),
        (helpers.TESTNET_DID + "/", 400),
        (helpers.TESTNET_DID + "invalidDID", 400),
    ]
)
def test_resolution_status_code(did_url, expected_status_code):
    url = helpers.RESOLVER_URL + helpers.PATH + did_url.replace("#", "%23")
    r = requests.get(url)

    assert r.status_code == expected_status_code
