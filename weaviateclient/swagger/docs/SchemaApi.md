# \SchemaApi

All URIs are relative to *https://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SchemaActionsCreate**](SchemaApi.md#SchemaActionsCreate) | **Post** /schema/actions | Create a new Action class in the schema.
[**SchemaActionsDelete**](SchemaApi.md#SchemaActionsDelete) | **Delete** /schema/actions/{className} | Remove an Action class (and all data in the instances) from the schema.
[**SchemaActionsPropertiesAdd**](SchemaApi.md#SchemaActionsPropertiesAdd) | **Post** /schema/actions/{className}/properties | Add a property to an Action class.
[**SchemaDump**](SchemaApi.md#SchemaDump) | **Get** /schema | Dump the current the database schema.
[**SchemaThingsCreate**](SchemaApi.md#SchemaThingsCreate) | **Post** /schema/things | Create a new Thing class in the schema.
[**SchemaThingsDelete**](SchemaApi.md#SchemaThingsDelete) | **Delete** /schema/things/{className} | Remove a Thing class (and all data in the instances) from the schema.
[**SchemaThingsPropertiesAdd**](SchemaApi.md#SchemaThingsPropertiesAdd) | **Post** /schema/things/{className}/properties | Add a property to a Thing class.


# **SchemaActionsCreate**
> Class SchemaActionsCreate(ctx, actionClass)
Create a new Action class in the schema.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **actionClass** | [**Class**](Class.md)|  | 

### Return type

[**Class**](Class.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SchemaActionsDelete**
> SchemaActionsDelete(ctx, className)
Remove an Action class (and all data in the instances) from the schema.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **className** | **string**|  | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SchemaActionsPropertiesAdd**
> Property SchemaActionsPropertiesAdd(ctx, className, body)
Add a property to an Action class.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **className** | **string**|  | 
  **body** | [**Property**](Property.md)|  | 

### Return type

[**Property**](Property.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SchemaDump**
> InlineResponse2001 SchemaDump(ctx, )
Dump the current the database schema.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse2001**](inline_response_200_1.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SchemaThingsCreate**
> Class SchemaThingsCreate(ctx, thingClass)
Create a new Thing class in the schema.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **thingClass** | [**Class**](Class.md)|  | 

### Return type

[**Class**](Class.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SchemaThingsDelete**
> SchemaThingsDelete(ctx, className)
Remove a Thing class (and all data in the instances) from the schema.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **className** | **string**|  | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SchemaThingsPropertiesAdd**
> Property SchemaThingsPropertiesAdd(ctx, className, body)
Add a property to a Thing class.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **className** | **string**|  | 
  **body** | [**Property**](Property.md)|  | 

### Return type

[**Property**](Property.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

