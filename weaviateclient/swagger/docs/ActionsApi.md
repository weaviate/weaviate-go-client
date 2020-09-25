# \ActionsApi

All URIs are relative to *https://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActionsCreate**](ActionsApi.md#ActionsCreate) | **Post** /actions | Create Actions between two Things (object and subject).
[**ActionsDelete**](ActionsApi.md#ActionsDelete) | **Delete** /actions/{id} | Delete an Action based on its UUID.
[**ActionsGet**](ActionsApi.md#ActionsGet) | **Get** /actions/{id} | Get a specific Action based on its UUID and a Thing UUID. Also available as Websocket bus.
[**ActionsList**](ActionsApi.md#ActionsList) | **Get** /actions | Get a list of Actions.
[**ActionsPatch**](ActionsApi.md#ActionsPatch) | **Patch** /actions/{id} | Update an Action based on its UUID (using patch semantics).
[**ActionsReferencesCreate**](ActionsApi.md#ActionsReferencesCreate) | **Post** /actions/{id}/references/{propertyName} | Add a single reference to a class-property.
[**ActionsReferencesDelete**](ActionsApi.md#ActionsReferencesDelete) | **Delete** /actions/{id}/references/{propertyName} | Delete the single reference that is given in the body from the list of references that this property has.
[**ActionsReferencesUpdate**](ActionsApi.md#ActionsReferencesUpdate) | **Put** /actions/{id}/references/{propertyName} | Replace all references to a class-property.
[**ActionsUpdate**](ActionsApi.md#ActionsUpdate) | **Put** /actions/{id} | Update an Action based on its UUID.
[**ActionsValidate**](ActionsApi.md#ActionsValidate) | **Post** /actions/validate | Validate an Action based on a schema.
[**BatchingActionsCreate**](ActionsApi.md#BatchingActionsCreate) | **Post** /batching/actions | Creates new Actions based on an Action template as a batch.


# **ActionsCreate**
> Action ActionsCreate(ctx, body)
Create Actions between two Things (object and subject).

Registers a new Action. Provided meta-data and schema values are validated.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Action**](Action.md)|  | 

### Return type

[**Action**](Action.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ActionsDelete**
> ActionsDelete(ctx, id)
Delete an Action based on its UUID.

Deletes an Action from the system.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Thing. | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ActionsGet**
> Action ActionsGet(ctx, id, optional)
Get a specific Action based on its UUID and a Thing UUID. Also available as Websocket bus.

Lists Actions.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Action. | 
 **optional** | ***ActionsApiActionsGetOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ActionsApiActionsGetOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **meta** | **optional.Bool**| Should additional meta information (e.g. about classified properties) be included? Defaults to false. | 
 **include** | **optional.String**| Include additional information, such as classification infos. Allowed values include: classification, _classification, vector, _vector, interpretation, _interpretation | 

### Return type

[**Action**](Action.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ActionsList**
> ActionsListResponse ActionsList(ctx, optional)
Get a list of Actions.

Lists all Actions in reverse order of creation, owned by the user that belongs to the used token.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ActionsApiActionsListOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ActionsApiActionsListOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **optional.Int64**| The maximum number of items to be returned per page. Default value is set in Weaviate config. | 
 **meta** | **optional.Bool**| Should additional meta information (e.g. about classified properties) be included? Defaults to false. | 
 **include** | **optional.String**| Include additional information, such as classification infos. Allowed values include: classification, _classification, vector, _vector, interpretation, _interpretation | 

### Return type

[**ActionsListResponse**](ActionsListResponse.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ActionsPatch**
> ActionsPatch(ctx, id, optional)
Update an Action based on its UUID (using patch semantics).

Updates an Action. This method supports json-merge style patch semantics (RFC 7396). Provided meta-data and schema values are validated. LastUpdateTime is set to the time this function is called.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Action. | 
 **optional** | ***ActionsApiActionsPatchOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ActionsApiActionsPatchOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of Action**](Action.md)| RFC 7396-style patch, the body contains the action object to merge into the existing action object. | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ActionsReferencesCreate**
> ActionsReferencesCreate(ctx, id, propertyName, body)
Add a single reference to a class-property.

Add a single reference to a class-property.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Action. | 
  **propertyName** | **string**| Unique name of the property related to the Action. | 
  **body** | [**SingleRef**](SingleRef.md)|  | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ActionsReferencesDelete**
> ActionsReferencesDelete(ctx, id, propertyName, body)
Delete the single reference that is given in the body from the list of references that this property has.

Delete the single reference that is given in the body from the list of references that this property has.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Action. | 
  **propertyName** | **string**| Unique name of the property related to the Action. | 
  **body** | [**SingleRef**](SingleRef.md)|  | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ActionsReferencesUpdate**
> ActionsReferencesUpdate(ctx, id, propertyName, body)
Replace all references to a class-property.

Replace all references to a class-property.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Action. | 
  **propertyName** | **string**| Unique name of the property related to the Action. | 
  **body** | [**MultipleRef**](MultipleRef.md)|  | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ActionsUpdate**
> Action ActionsUpdate(ctx, id, body)
Update an Action based on its UUID.

Updates an Action's data. Given meta-data and schema values are validated. LastUpdateTime is set to the time this function is called.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Action. | 
  **body** | [**Action**](Action.md)|  | 

### Return type

[**Action**](Action.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ActionsValidate**
> ActionsValidate(ctx, body)
Validate an Action based on a schema.

Validate an Action's schema and meta-data. It has to be based on a schema, which is related to the given Action to be accepted by this validation.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Action**](Action.md)|  | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **BatchingActionsCreate**
> []ActionsGetResponse BatchingActionsCreate(ctx, body)
Creates new Actions based on an Action template as a batch.

Register new Actions in bulk. Given meta-data and schema values are validated.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body1**](Body1.md)|  | 

### Return type

[**[]ActionsGetResponse**](ActionsGetResponse.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

