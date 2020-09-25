# \GraphqlApi

All URIs are relative to *https://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GraphqlBatch**](GraphqlApi.md#GraphqlBatch) | **Post** /graphql/batch | Get a response based on GraphQL.
[**GraphqlPost**](GraphqlApi.md#GraphqlPost) | **Post** /graphql | Get a response based on GraphQL


# **GraphqlBatch**
> GraphQlResponses GraphqlBatch(ctx, body)
Get a response based on GraphQL.

Perform a batched GraphQL query

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**GraphQlQueries**](GraphQlQueries.md)| The GraphQL queries. | 

### Return type

[**GraphQlResponses**](GraphQLResponses.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GraphqlPost**
> GraphQlResponse GraphqlPost(ctx, body)
Get a response based on GraphQL

Get an object based on GraphQL

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**GraphQlQuery**](GraphQlQuery.md)| The GraphQL query request parameters. | 

### Return type

[**GraphQlResponse**](GraphQLResponse.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

