# \ClassificationsApi

All URIs are relative to *https://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ClassificationsGet**](ClassificationsApi.md#ClassificationsGet) | **Get** /classifications/{id} | View previously created classification
[**ClassificationsPost**](ClassificationsApi.md#ClassificationsPost) | **Post** /classifications/ | Starts a classification.


# **ClassificationsGet**
> Classification ClassificationsGet(ctx, id)
View previously created classification

Get status, results and metadata of a previously created classification

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| classification id | 

### Return type

[**Classification**](Classification.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ClassificationsPost**
> Classification ClassificationsPost(ctx, params)
Starts a classification.

Trigger a classification based on the specified params. Classifications will run in the background, use GET /classifications/<id> to retrieve the status of your classificaiton.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **params** | [**Classification**](Classification.md)| parameters to start a classification | 

### Return type

[**Classification**](Classification.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

