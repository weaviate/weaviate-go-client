# \DiscoveryApi

All URIs are relative to *https://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**WellKnownOpenidConfigurationGet**](DiscoveryApi.md#WellKnownOpenidConfigurationGet) | **Get** /.well-known/openid-configuration | OIDC discovery information if OIDC auth is enabled


# **WellKnownOpenidConfigurationGet**
> InlineResponse2002 WellKnownOpenidConfigurationGet(ctx, )
OIDC discovery information if OIDC auth is enabled

OIDC Discovery page, redirects to the token issuer if one is configured

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse2002**](inline_response_200_2.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

