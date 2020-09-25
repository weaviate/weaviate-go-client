# \ContextionaryAPIApi

All URIs are relative to *https://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**C11yConcepts**](ContextionaryAPIApi.md#C11yConcepts) | **Get** /c11y/concepts/{concept} | Checks if a concept is part of the contextionary.
[**C11yCorpusGet**](ContextionaryAPIApi.md#C11yCorpusGet) | **Post** /c11y/corpus | Checks if a word or wordString is part of the contextionary.
[**C11yExtensions**](ContextionaryAPIApi.md#C11yExtensions) | **Post** /c11y/extensions/ | Extend the contextionary with custom concepts
[**C11yWords**](ContextionaryAPIApi.md#C11yWords) | **Get** /c11y/words/{words} | Checks if a word or wordString is part of the contextionary.


# **C11yConcepts**
> C11yWordsResponse C11yConcepts(ctx, concept)
Checks if a concept is part of the contextionary.

Checks if a concept is part of the contextionary. Concepts should be concatenated as described here: https://github.com/semi-technologies/weaviate/blob/master/docs/en/use/schema-schema.md#camelcase

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **concept** | **string**| CamelCase list of words to validate. | 

### Return type

[**C11yWordsResponse**](C11yWordsResponse.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **C11yCorpusGet**
> C11yCorpusGet(ctx, corpus)
Checks if a word or wordString is part of the contextionary.

Analyzes a sentence based on the contextionary

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **corpus** | [**Corpus**](Corpus.md)| A text corpus | 

### Return type

 (empty response body)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **C11yExtensions**
> C11yExtension C11yExtensions(ctx, extension)
Extend the contextionary with custom concepts

Extend the contextionary with your own custom concepts

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **extension** | [**C11yExtension**](C11yExtension.md)| Description and definition of the concept to extend the contextionary with | 

### Return type

[**C11yExtension**](C11yExtension.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **C11yWords**
> C11yWordsResponse C11yWords(ctx, words)
Checks if a word or wordString is part of the contextionary.

Checks if a word or wordString is part of the contextionary. Words should be concatenated as described here: https://github.com/semi-technologies/weaviate/blob/master/docs/en/use/schema-schema.md#camelcase

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **words** | **string**| CamelCase list of words to validate. | 

### Return type

[**C11yWordsResponse**](C11yWordsResponse.md)

### Authorization

[oidc](../README.md#oidc)

### HTTP request headers

 - **Content-Type**: application/yaml, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

