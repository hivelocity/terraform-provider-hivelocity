# \TokenApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteTokenIdResource**](TokenApi.md#DeleteTokenIdResource) | **Delete** /token/{token} | Deletes the Public API Token
[**GetTokenIdResource**](TokenApi.md#GetTokenIdResource) | **Get** /token/{token} | Returns Public API Token
[**GetTokenResource**](TokenApi.md#GetTokenResource) | **Get** /token/ | Returns a list of Public API Tokens for the current user
[**PostTokenResource**](TokenApi.md#PostTokenResource) | **Post** /token/ | Create a new Public API Token for the current user
[**PutTokenIdResource**](TokenApi.md#PutTokenIdResource) | **Put** /token/{token} | Updates the Public API Token


# **DeleteTokenIdResource**
> DeleteTokenIdResource(ctx, token)
Deletes the Public API Token

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **token** | **string**| Public API Token | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetTokenIdResource**
> PublicApiTokenDump GetTokenIdResource(ctx, token, optional)
Returns Public API Token

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **token** | **string**| Public API Token | 
 **optional** | ***TokenApiGetTokenIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a TokenApiGetTokenIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**PublicApiTokenDump**](PublicApiToken-dump.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetTokenResource**
> []PublicApiTokenDump GetTokenResource(ctx, optional)
Returns a list of Public API Tokens for the current user

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***TokenApiGetTokenResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a TokenApiGetTokenResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]PublicApiTokenDump**](PublicApiToken-dump.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostTokenResource**
> PublicApiTokenDump PostTokenResource(ctx, payload, optional)
Create a new Public API Token for the current user

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **payload** | [**PublicApiTokenLoad**](PublicApiTokenLoad.md)|  | 
 **optional** | ***TokenApiPostTokenResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a TokenApiPostTokenResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**PublicApiTokenDump**](PublicApiToken-dump.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutTokenIdResource**
> PublicApiTokenDump PutTokenIdResource(ctx, token, payload, optional)
Updates the Public API Token

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **token** | **string**| Public API Token | 
  **payload** | [**PublicApiTokenLoad**](PublicApiTokenLoad.md)|  | 
 **optional** | ***TokenApiPutTokenIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a TokenApiPutTokenIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**PublicApiTokenDump**](PublicApiToken-dump.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

