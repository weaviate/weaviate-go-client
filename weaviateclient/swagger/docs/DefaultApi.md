# \DefaultApi

All URIs are relative to *https://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**WeaviateRoot**](DefaultApi.md#WeaviateRoot) | **Get** / | 
[**WeaviateWellknownLiveness**](DefaultApi.md#WeaviateWellknownLiveness) | **Get** /.well-known/live | 
[**WeaviateWellknownReadiness**](DefaultApi.md#WeaviateWellknownReadiness) | **Get** /.well-known/ready | 


# **WeaviateRoot**
> InlineResponse200 WeaviateRoot(ctx, )


Home. Discover the REST API

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse200**](inline_response_200.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **WeaviateWellknownLiveness**
> WeaviateWellknownLiveness(ctx, )


Determines whether the application is alive. Can be used for kubernetes liveness probe

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **WeaviateWellknownReadiness**
> WeaviateWellknownReadiness(ctx, )


Determines whether the application is ready to receive traffic. Can be used for kubernetes readiness probe.

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

