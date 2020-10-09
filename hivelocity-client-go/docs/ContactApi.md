# \ContactApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteContactIdResource**](ContactApi.md#DeleteContactIdResource) | **Delete** /contact/{contactId} | Delete a Contact
[**GetContactIdResource**](ContactApi.md#GetContactIdResource) | **Get** /contact/{contactId} | Return details of a specific Contact
[**GetContactResource**](ContactApi.md#GetContactResource) | **Get** /contact/ | Return a list with all Contacts
[**PostContactResource**](ContactApi.md#PostContactResource) | **Post** /contact/ | Create a new Contact
[**PostPasswordReset**](ContactApi.md#PostPasswordReset) | **Post** /contact/password | Recieve a request to reset password, and send a link by email with token
[**PutContactIdResource**](ContactApi.md#PutContactIdResource) | **Put** /contact/{contactId} | Edit a Contact
[**PutPasswordReset**](ContactApi.md#PutPasswordReset) | **Put** /contact/password | Recieve a token and password, verify the user and reset your password
[**PutPasswordResource**](ContactApi.md#PutPasswordResource) | **Put** /contact/{contactId}/password | Change password for a Contact


# **DeleteContactIdResource**
> DeleteContactIdResource(ctx, contactId)
Delete a Contact

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **contactId** | **int32**| ID of Contact to View / Update | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetContactIdResource**
> ContactDump GetContactIdResource(ctx, contactId, optional)
Return details of a specific Contact

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **contactId** | **int32**| ID of Contact to View / Update | 
 **optional** | ***ContactApiGetContactIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ContactApiGetContactIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**ContactDump**](Contact-dump.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetContactResource**
> []ContactDump GetContactResource(ctx, optional)
Return a list with all Contacts

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ContactApiGetContactResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ContactApiGetContactResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]ContactDump**](Contact-dump.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostContactResource**
> ContactDump PostContactResource(ctx, payload, optional)
Create a new Contact

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **payload** | [**ContactCreate**](ContactCreate.md)|  | 
 **optional** | ***ContactApiPostContactResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ContactApiPostContactResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**ContactDump**](Contact-dump.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostPasswordReset**
> PostPasswordReset(ctx, payload)
Recieve a request to reset password, and send a link by email with token

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **payload** | [**PasswordRequest**](PasswordRequest.md)|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutContactIdResource**
> ContactDump PutContactIdResource(ctx, contactId, payload, optional)
Edit a Contact

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **contactId** | **int32**| ID of Contact to View / Update | 
  **payload** | [**ContactUpdate**](ContactUpdate.md)|  | 
 **optional** | ***ContactApiPutContactIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ContactApiPutContactIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**ContactDump**](Contact-dump.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutPasswordReset**
> PutPasswordReset(ctx, payload)
Recieve a token and password, verify the user and reset your password

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **payload** | [**PasswordResponse**](PasswordResponse.md)|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutPasswordResource**
> PutPasswordResource(ctx, contactId, payload)
Change password for a Contact

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **contactId** | **int32**| \&quot;ID of Contact to change password\&quot; | 
  **payload** | [**Password**](Password.md)|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

