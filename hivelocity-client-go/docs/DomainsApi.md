# \DomainsApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteARecordIdResource**](DomainsApi.md#DeleteARecordIdResource) | **Delete** /domains/{domainID}/a-record/{recordId} | Delete a single A Record value for given zone
[**DeleteAaaaRecordIdResource**](DomainsApi.md#DeleteAaaaRecordIdResource) | **Delete** /domains/{domainID}/aaaa-record/{recordId} | Delete an AAAA Record for given zone
[**DeleteDomainIdResource**](DomainsApi.md#DeleteDomainIdResource) | **Delete** /domains/{domainId} | Delete domain for given id
[**DeleteMxRecordIdResource**](DomainsApi.md#DeleteMxRecordIdResource) | **Delete** /domains/{domainID}/mx-record/{recordId} | Delete an MX Record for given zone
[**GetARecordIdResource**](DomainsApi.md#GetARecordIdResource) | **Get** /domains/{domainID}/a-record/{recordId} | Return Single A Records found for given id
[**GetARecordResource**](DomainsApi.md#GetARecordResource) | **Get** /domains/{domainID}/a-record | Return all A Records found for given zone
[**GetAaaaRecordIdResource**](DomainsApi.md#GetAaaaRecordIdResource) | **Get** /domains/{domainID}/aaaa-record/{recordId} | Return Single AAAA Records found for given id
[**GetAaaaRecordResource**](DomainsApi.md#GetAaaaRecordResource) | **Get** /domains/{domainID}/aaaa-record | Return all AAAA Records found for given zone
[**GetDomainIdResource**](DomainsApi.md#GetDomainIdResource) | **Get** /domains/{domainId} | Return a domain for given id
[**GetDomainResource**](DomainsApi.md#GetDomainResource) | **Get** /domains/ | Return the domains of a client
[**GetMxRecordIdResource**](DomainsApi.md#GetMxRecordIdResource) | **Get** /domains/{domainID}/mx-record/{recordId} | Return Single MX Records found for given id
[**GetMxRecordResource**](DomainsApi.md#GetMxRecordResource) | **Get** /domains/{domainID}/mx-record | Return all MX Records found for given zone
[**PostARecordResource**](DomainsApi.md#PostARecordResource) | **Post** /domains/{domainID}/a-record | Create a new A Record for given zone
[**PostAaaaRecordResource**](DomainsApi.md#PostAaaaRecordResource) | **Post** /domains/{domainID}/aaaa-record | Create a new AAAA Record for given zone
[**PostDomainResource**](DomainsApi.md#PostDomainResource) | **Post** /domains/ | Create a new domain
[**PostMxRecordResource**](DomainsApi.md#PostMxRecordResource) | **Post** /domains/{domainID}/mx-record | Create a new MX Record for given zone
[**PutARecordIdResource**](DomainsApi.md#PutARecordIdResource) | **Put** /domains/{domainID}/a-record/{recordId} | Update an A Record for given zone
[**PutAaaaRecordIdResource**](DomainsApi.md#PutAaaaRecordIdResource) | **Put** /domains/{domainID}/aaaa-record/{recordId} | Update an AAAA Record for given id
[**PutMxRecordIdResource**](DomainsApi.md#PutMxRecordIdResource) | **Put** /domains/{domainID}/mx-record/{recordId} | Update an MX Record for given zone


# **DeleteARecordIdResource**
> DeleteARecordIdResource(ctx, domainID, recordId)
Delete a single A Record value for given zone

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
  **recordId** | **int32**|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteAaaaRecordIdResource**
> DeleteAaaaRecordIdResource(ctx, domainID, recordId)
Delete an AAAA Record for given zone

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
  **recordId** | **int32**|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteDomainIdResource**
> DeleteDomainIdResource(ctx, domainId)
Delete domain for given id

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainId** | **int32**|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteMxRecordIdResource**
> DeleteMxRecordIdResource(ctx, domainID, recordId)
Delete an MX Record for given zone

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
  **recordId** | **int32**|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetARecordIdResource**
> SingleARecordReturn GetARecordIdResource(ctx, domainID, recordId, optional)
Return Single A Records found for given id

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
  **recordId** | **int32**|  | 
 **optional** | ***DomainsApiGetARecordIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiGetARecordIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**SingleARecordReturn**](SingleARecordReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetARecordResource**
> []ARecordReturn GetARecordResource(ctx, domainID, optional)
Return all A Records found for given zone

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
 **optional** | ***DomainsApiGetARecordResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiGetARecordResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]ARecordReturn**](ARecordReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAaaaRecordIdResource**
> SingleARecordReturn GetAaaaRecordIdResource(ctx, domainID, recordId, optional)
Return Single AAAA Records found for given id

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
  **recordId** | **int32**|  | 
 **optional** | ***DomainsApiGetAaaaRecordIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiGetAaaaRecordIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**SingleARecordReturn**](SingleARecordReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAaaaRecordResource**
> []AaaaRecordReturn GetAaaaRecordResource(ctx, domainID, optional)
Return all AAAA Records found for given zone

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
 **optional** | ***DomainsApiGetAaaaRecordResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiGetAaaaRecordResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]AaaaRecordReturn**](AAAARecordReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDomainIdResource**
> DomainReturn GetDomainIdResource(ctx, domainId, optional)
Return a domain for given id

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainId** | **int32**|  | 
 **optional** | ***DomainsApiGetDomainIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiGetDomainIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DomainReturn**](DomainReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDomainResource**
> []DomainReturn GetDomainResource(ctx, optional)
Return the domains of a client

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***DomainsApiGetDomainResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiGetDomainResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]DomainReturn**](DomainReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetMxRecordIdResource**
> SingleMxRecordReturn GetMxRecordIdResource(ctx, domainID, recordId, optional)
Return Single MX Records found for given id

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
  **recordId** | **int32**|  | 
 **optional** | ***DomainsApiGetMxRecordIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiGetMxRecordIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**SingleMxRecordReturn**](SingleMXRecordReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetMxRecordResource**
> []MxRecordReturn GetMxRecordResource(ctx, domainID, optional)
Return all MX Records found for given zone

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
 **optional** | ***DomainsApiGetMxRecordResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiGetMxRecordResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]MxRecordReturn**](MXRecordReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostARecordResource**
> SingleARecordReturn PostARecordResource(ctx, domainID, payload, optional)
Create a new A Record for given zone

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
  **payload** | [**ARecordCreate**](ARecordCreate.md)|  | 
 **optional** | ***DomainsApiPostARecordResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiPostARecordResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**SingleARecordReturn**](SingleARecordReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostAaaaRecordResource**
> SingleARecordReturn PostAaaaRecordResource(ctx, domainID, payload, optional)
Create a new AAAA Record for given zone

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
  **payload** | [**AaaaRecordCreate**](AaaaRecordCreate.md)|  | 
 **optional** | ***DomainsApiPostAaaaRecordResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiPostAaaaRecordResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**SingleARecordReturn**](SingleARecordReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostDomainResource**
> DomainReturn PostDomainResource(ctx, payload, optional)
Create a new domain

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **payload** | [**DomainCreate**](DomainCreate.md)|  | 
 **optional** | ***DomainsApiPostDomainResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiPostDomainResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DomainReturn**](DomainReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostMxRecordResource**
> SingleMxRecordReturn PostMxRecordResource(ctx, domainID, payload, optional)
Create a new MX Record for given zone

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
  **payload** | [**SingleMxRecordCreate**](SingleMxRecordCreate.md)|  | 
 **optional** | ***DomainsApiPostMxRecordResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiPostMxRecordResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**SingleMxRecordReturn**](SingleMXRecordReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutARecordIdResource**
> SingleARecordReturn PutARecordIdResource(ctx, domainID, recordId, payload, optional)
Update an A Record for given zone

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
  **recordId** | **int32**|  | 
  **payload** | [**ARecordUpdate**](ARecordUpdate.md)|  | 
 **optional** | ***DomainsApiPutARecordIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiPutARecordIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**SingleARecordReturn**](SingleARecordReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutAaaaRecordIdResource**
> SingleARecordReturn PutAaaaRecordIdResource(ctx, domainID, recordId, payload, optional)
Update an AAAA Record for given id

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
  **recordId** | **int32**|  | 
  **payload** | [**AaaaRecordUpdate**](AaaaRecordUpdate.md)|  | 
 **optional** | ***DomainsApiPutAaaaRecordIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiPutAaaaRecordIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**SingleARecordReturn**](SingleARecordReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutMxRecordIdResource**
> SingleMxRecordReturn PutMxRecordIdResource(ctx, domainID, recordId, payload, optional)
Update an MX Record for given zone

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **domainID** | **int32**|  | 
  **recordId** | **int32**|  | 
  **payload** | [**SingleMxRecordUpdate**](SingleMxRecordUpdate.md)|  | 
 **optional** | ***DomainsApiPutMxRecordIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DomainsApiPutMxRecordIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**SingleMxRecordReturn**](SingleMXRecordReturn.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

