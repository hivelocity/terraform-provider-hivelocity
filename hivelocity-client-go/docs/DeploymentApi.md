# \DeploymentApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteDeploymentIdResource**](DeploymentApi.md#DeleteDeploymentIdResource) | **Delete** /deploy/{deploymentId} | Delete the specified deployment
[**GetDeploymentIdResource**](DeploymentApi.md#GetDeploymentIdResource) | **Get** /deploy/{deploymentId} | Return a dictionary with deployment information
[**GetDeploymentResource**](DeploymentApi.md#GetDeploymentResource) | **Get** /deploy/ | Return a list with all client deployments
[**PostDeploymentIdResource**](DeploymentApi.md#PostDeploymentIdResource) | **Post** /deploy/{deploymentId} | Input a billing info id and script to process and finish a deployment
[**PostDeploymentResource**](DeploymentApi.md#PostDeploymentResource) | **Post** /deploy/ | Start a new deployment
[**PutDeploymentIdResource**](DeploymentApi.md#PutDeploymentIdResource) | **Put** /deploy/{deploymentId} | Receive product, quantity and options to be added on the deployment


# **DeleteDeploymentIdResource**
> DeleteDeploymentIdResource(ctx, deploymentId)
Delete the specified deployment

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deploymentId** | **int32**| Id of the deployment to interact with | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeploymentIdResource**
> Deployment GetDeploymentIdResource(ctx, deploymentId, optional)
Return a dictionary with deployment information

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deploymentId** | **int32**| Id of the deployment to interact with | 
 **optional** | ***DeploymentApiGetDeploymentIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeploymentApiGetDeploymentIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**Deployment**](Deployment.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeploymentResource**
> []Deployment GetDeploymentResource(ctx, optional)
Return a list with all client deployments

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***DeploymentApiGetDeploymentResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeploymentApiGetDeploymentResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]Deployment**](Deployment.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostDeploymentIdResource**
> PostDeploymentIdResource(ctx, deploymentId, payload)
Input a billing info id and script to process and finish a deployment

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deploymentId** | **int32**| Id of the deployment to interact with | 
  **payload** | [**DeploymentStart**](DeploymentStart.md)|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostDeploymentResource**
> Deployment PostDeploymentResource(ctx, optional)
Start a new deployment

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***DeploymentApiPostDeploymentResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeploymentApiPostDeploymentResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deploymentName** | **optional.String**|  | 
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**Deployment**](Deployment.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutDeploymentIdResource**
> Deployment PutDeploymentIdResource(ctx, deploymentId, payload, optional)
Receive product, quantity and options to be added on the deployment

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deploymentId** | **int32**| Id of the deployment to interact with | 
  **payload** | [**DeploymentCustomization**](DeploymentCustomization.md)|  | 
 **optional** | ***DeploymentApiPutDeploymentIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeploymentApiPutDeploymentIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**Deployment**](Deployment.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

