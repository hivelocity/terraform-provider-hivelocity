# \SshKeyApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteSshKeyIdResource**](SshKeyApi.md#DeleteSshKeyIdResource) | **Delete** /ssh_key/{sshKeyId} | Removes public ssh key
[**GetSshKeyIdResource**](SshKeyApi.md#GetSshKeyIdResource) | **Get** /ssh_key/{sshKeyId} | Get public ssh key
[**GetSshKeyResource**](SshKeyApi.md#GetSshKeyResource) | **Get** /ssh_key/ | Gets all public ssh key
[**PostSshKeyResource**](SshKeyApi.md#PostSshKeyResource) | **Post** /ssh_key/ | Adds public ssh key
[**PutSshKeyIdResource**](SshKeyApi.md#PutSshKeyIdResource) | **Put** /ssh_key/{sshKeyId} | Updates public ssh key


# **DeleteSshKeyIdResource**
> DeleteSshKeyIdResource(ctx, sshKeyId)
Removes public ssh key

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **sshKeyId** | **int32**|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetSshKeyIdResource**
> SshKeyResponse GetSshKeyIdResource(ctx, sshKeyId, optional)
Get public ssh key

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **sshKeyId** | **int32**|  | 
 **optional** | ***SshKeyApiGetSshKeyIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a SshKeyApiGetSshKeyIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**SshKeyResponse**](SshKeyResponse.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetSshKeyResource**
> []SshKeyResponse GetSshKeyResource(ctx, optional)
Gets all public ssh key

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***SshKeyApiGetSshKeyResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a SshKeyApiGetSshKeyResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]SshKeyResponse**](SshKeyResponse.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostSshKeyResource**
> SshKeyResponse PostSshKeyResource(ctx, payload, optional)
Adds public ssh key

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **payload** | [**SshKey**](SshKey.md)|  | 
 **optional** | ***SshKeyApiPostSshKeyResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a SshKeyApiPostSshKeyResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**SshKeyResponse**](SshKeyResponse.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutSshKeyIdResource**
> SshKeyResponse PutSshKeyIdResource(ctx, sshKeyId, payload, optional)
Updates public ssh key

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **sshKeyId** | **int32**|  | 
  **payload** | [**SshKeyUpdate**](SshKeyUpdate.md)|  | 
 **optional** | ***SshKeyApiPutSshKeyIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a SshKeyApiPutSshKeyIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**SshKeyResponse**](SshKeyResponse.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

