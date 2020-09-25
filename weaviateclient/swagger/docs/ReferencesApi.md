# \ReferencesApi

All URIs are relative to *https://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BatchingReferencesCreate**](ReferencesApi.md#BatchingReferencesCreate) | **Post** /batching/references | Creates new Cross-References between arbitrary classes in bulk.


# **BatchingReferencesCreate**
> []BatchReferenceResponse BatchingReferencesCreate(ctx, body)
Creates new Cross-References between arbitrary classes in bulk.

Register cross-references between any class items (things or actions) in bulk.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**[]BatchReference**](BatchReference.md)| A list of references to be batched. The ideal size depends on the used database connector. Please see the documentation of the used connector for help | 

### Return type

[**[]BatchReferenceResponse**](BatchReferenceResponse.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

