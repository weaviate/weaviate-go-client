# \MetaApi

All URIs are relative to *https://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**MetaGet**](MetaApi.md#MetaGet) | **Get** /meta | Returns meta information of the current Weaviate instance.


# **MetaGet**
> Meta MetaGet(ctx, )
Returns meta information of the current Weaviate instance.

Gives meta information about the server and can be used to provide information to another Weaviate instance that wants to interact with the current instance.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**Meta**](Meta.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

