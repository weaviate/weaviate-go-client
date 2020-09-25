# \ThingsApi

All URIs are relative to *https://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BatchingThingsCreate**](ThingsApi.md#BatchingThingsCreate) | **Post** /batching/things | Creates new Things based on a Thing template as a batch.
[**ThingsCreate**](ThingsApi.md#ThingsCreate) | **Post** /things | Create a new Thing based on a Thing template.
[**ThingsDelete**](ThingsApi.md#ThingsDelete) | **Delete** /things/{id} | Delete a Thing based on its UUID.
[**ThingsGet**](ThingsApi.md#ThingsGet) | **Get** /things/{id} | Get a Thing based on its UUID.
[**ThingsList**](ThingsApi.md#ThingsList) | **Get** /things | Get a list of Things.
[**ThingsPatch**](ThingsApi.md#ThingsPatch) | **Patch** /things/{id} | Update a Thing based on its UUID (using patch semantics).
[**ThingsReferencesCreate**](ThingsApi.md#ThingsReferencesCreate) | **Post** /things/{id}/references/{propertyName} | Add a single reference to a class-property.
[**ThingsReferencesDelete**](ThingsApi.md#ThingsReferencesDelete) | **Delete** /things/{id}/references/{propertyName} | Delete the single reference that is given in the body from the list of references that this property has.
[**ThingsReferencesUpdate**](ThingsApi.md#ThingsReferencesUpdate) | **Put** /things/{id}/references/{propertyName} | Replace all references to a class-property.
[**ThingsUpdate**](ThingsApi.md#ThingsUpdate) | **Put** /things/{id} | Update a Thing based on its UUID.
[**ThingsValidate**](ThingsApi.md#ThingsValidate) | **Post** /things/validate | Validate Things schema.


# **BatchingThingsCreate**
> []ThingsGetResponse BatchingThingsCreate(ctx, body)
Creates new Things based on a Thing template as a batch.

Register new Things in bulk. Provided meta-data and schema values are validated.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body**](Body.md)|  | 

### Return type

[**[]ThingsGetResponse**](ThingsGetResponse.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ThingsCreate**
> Thing ThingsCreate(ctx, body)
Create a new Thing based on a Thing template.

Registers a new Thing. Given meta-data and schema values are validated.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Thing**](Thing.md)|  | 

### Return type

[**Thing**](Thing.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ThingsDelete**
> ThingsDelete(ctx, id)
Delete a Thing based on its UUID.

Deletes a Thing from the system. All Actions pointing to this Thing, where the Thing is the object of the Action, are also being deleted.

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

# **ThingsGet**
> Thing ThingsGet(ctx, id, optional)
Get a Thing based on its UUID.

Returns a particular Thing data.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Thing. | 
 **optional** | ***ThingsApiThingsGetOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ThingsApiThingsGetOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **meta** | **optional.Bool**| Should additional meta information (e.g. about classified properties) be included? Defaults to false. | 
 **include** | **optional.String**| Include additional information, such as classification infos. Allowed values include: classification, _classification, vector, _vector, interpretation, _interpretation | 

### Return type

[**Thing**](Thing.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ThingsList**
> ThingsListResponse ThingsList(ctx, optional)
Get a list of Things.

Lists all Things in reverse order of creation, owned by the user that belongs to the used token.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ThingsApiThingsListOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ThingsApiThingsListOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **optional.Int64**| The maximum number of items to be returned per page. Default value is set in Weaviate config. | 
 **meta** | **optional.Bool**| Should additional meta information (e.g. about classified properties) be included? Defaults to false. | 
 **include** | **optional.String**| Include additional information, such as classification infos. Allowed values include: classification, _classification, vector, _vector, interpretation, _interpretation | 

### Return type

[**ThingsListResponse**](ThingsListResponse.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ThingsPatch**
> ThingsPatch(ctx, id, optional)
Update a Thing based on its UUID (using patch semantics).

Updates a Thing's data. This method supports patch semantics. Given meta-data and schema values are validated. LastUpdateTime is set to the time this function is called.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Thing. | 
 **optional** | ***ThingsApiThingsPatchOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ThingsApiThingsPatchOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of Thing**](Thing.md)| RFC 7396-style patch, the body contains the thing object to merge into the existing thing object. | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ThingsReferencesCreate**
> ThingsReferencesCreate(ctx, id, propertyName, body)
Add a single reference to a class-property.

Add a single reference to a class-property.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Thing. | 
  **propertyName** | **string**| Unique name of the property related to the Thing. | 
  **body** | [**SingleRef**](SingleRef.md)|  | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ThingsReferencesDelete**
> ThingsReferencesDelete(ctx, id, propertyName, body)
Delete the single reference that is given in the body from the list of references that this property has.

Delete the single reference that is given in the body from the list of references that this property has.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Thing. | 
  **propertyName** | **string**| Unique name of the property related to the Thing. | 
  **body** | [**SingleRef**](SingleRef.md)|  | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ThingsReferencesUpdate**
> ThingsReferencesUpdate(ctx, id, propertyName, body)
Replace all references to a class-property.

Replace all references to a class-property.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Thing. | 
  **propertyName** | **string**| Unique name of the property related to the Thing. | 
  **body** | [**MultipleRef**](MultipleRef.md)|  | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ThingsUpdate**
> Thing ThingsUpdate(ctx, id, body)
Update a Thing based on its UUID.

Updates a Thing's data. Given meta-data and schema values are validated. LastUpdateTime is set to the time this function is called.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | [**string**](.md)| Unique ID of the Thing. | 
  **body** | [**Thing**](Thing.md)|  | 

### Return type

[**Thing**](Thing.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ThingsValidate**
> ThingsValidate(ctx, body)
Validate Things schema.

Validate a Thing's schema and meta-data. It has to be based on a schema, which is related to the given Thing to be accepted by this validation.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Thing**](Thing.md)|  | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

