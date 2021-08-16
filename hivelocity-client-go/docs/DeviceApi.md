# \DeviceApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetAllDeviceTagOrderResource**](DeviceApi.md#GetAllDeviceTagOrderResource) | **Get** /device/tags-order/all | Get all device tags order
[**GetClientDeviceTagOrderResource**](DeviceApi.md#GetClientDeviceTagOrderResource) | **Get** /device/tags-order | Get device tags order for current user
[**GetClientDeviceTagResource**](DeviceApi.md#GetClientDeviceTagResource) | **Get** /device/tags | Get all device tags for current client
[**GetDeviceIdResource**](DeviceApi.md#GetDeviceIdResource) | **Get** /device/{deviceId} | Returns detailed information for a Single Device
[**GetDeviceIpminatRuleResource**](DeviceApi.md#GetDeviceIpminatRuleResource) | **Get** /device/{deviceId}/ipmi/nat | Clear NAT rules based on the device client id
[**GetDeviceResource**](DeviceApi.md#GetDeviceResource) | **Get** /device/ | Returns Active Devices and basic MetaData
[**GetDeviceTagIdResource**](DeviceApi.md#GetDeviceTagIdResource) | **Get** /device/{deviceId}/tags | Get device tags
[**GetInitialCredsIdResource**](DeviceApi.md#GetInitialCredsIdResource) | **Get** /device/{deviceId}/initial-creds | Returns initial password for the device
[**GetIpmiInfoIdResource**](DeviceApi.md#GetIpmiInfoIdResource) | **Get** /device/{deviceId}/ipmi | Returns IPMI info data
[**GetIpmiInfoLoginDataResource**](DeviceApi.md#GetIpmiInfoLoginDataResource) | **Get** /device/{deviceId}/ipmi/login-data | Returns IPMI login credentials
[**GetIpmiThresholdsIdResource**](DeviceApi.md#GetIpmiThresholdsIdResource) | **Get** /device/{deviceId}/ipmi/thresholds | Returns IPMI thresholds data
[**GetIpmiValidLoginIdResource**](DeviceApi.md#GetIpmiValidLoginIdResource) | **Get** /device/{deviceId}/ipmi/valid-login | Returns if device have valid credentials for IPMI login
[**GetPowerResource**](DeviceApi.md#GetPowerResource) | **Get** /device/{deviceId}/power | Get device&#39;s current power status
[**PostDeviceIpmiWhitelistResource**](DeviceApi.md#PostDeviceIpmiWhitelistResource) | **Post** /device/{deviceId}/ipmi/whitelist/ | Add a public IP on IPMI whitelist
[**PostPowerResource**](DeviceApi.md#PostPowerResource) | **Post** /device/{deviceId}/power | Apply action to device power
[**PutClientDeviceTagOrderResource**](DeviceApi.md#PutClientDeviceTagOrderResource) | **Put** /device/tags-order | Update device tags order for current user
[**PutDeviceIdResource**](DeviceApi.md#PutDeviceIdResource) | **Put** /device/{deviceId} | Updates Device MetaData for a Single Device
[**PutDeviceTagIdResource**](DeviceApi.md#PutDeviceTagIdResource) | **Put** /device/{deviceId}/tags | Update device tags
[**PutIpmiDevicesThresholdsIdResource**](DeviceApi.md#PutIpmiDevicesThresholdsIdResource) | **Put** /device/ipmi/thresholds | Updates IPMI thresholds for device list
[**PutIpmiThresholdsIdResource**](DeviceApi.md#PutIpmiThresholdsIdResource) | **Put** /device/{deviceId}/ipmi/thresholds | Updates IPMI thresholds data


# **GetAllDeviceTagOrderResource**
> DeviceTag GetAllDeviceTagOrderResource(ctx, optional)
Get all device tags order

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***DeviceApiGetAllDeviceTagOrderResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiGetAllDeviceTagOrderResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DeviceTag**](DeviceTag.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetClientDeviceTagOrderResource**
> DeviceTag GetClientDeviceTagOrderResource(ctx, optional)
Get device tags order for current user

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***DeviceApiGetClientDeviceTagOrderResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiGetClientDeviceTagOrderResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DeviceTag**](DeviceTag.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetClientDeviceTagResource**
> DeviceTag GetClientDeviceTagResource(ctx, optional)
Get all device tags for current client

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***DeviceApiGetClientDeviceTagResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiGetClientDeviceTagResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DeviceTag**](DeviceTag.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceIdResource**
> Device GetDeviceIdResource(ctx, deviceId, optional)
Returns detailed information for a Single Device

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to View / Update | 
 **optional** | ***DeviceApiGetDeviceIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiGetDeviceIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**Device**](Device.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceIpminatRuleResource**
> GetDeviceIpminatRuleResource(ctx, deviceId)
Clear NAT rules based on the device client id

Returns success or error

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of a client Device | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceResource**
> []Device GetDeviceResource(ctx, optional)
Returns Active Devices and basic MetaData

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***DeviceApiGetDeviceResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiGetDeviceResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]Device**](Device.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDeviceTagIdResource**
> DeviceTag GetDeviceTagIdResource(ctx, deviceId, optional)
Get device tags

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to View / Update | 
 **optional** | ***DeviceApiGetDeviceTagIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiGetDeviceTagIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DeviceTag**](DeviceTag.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetInitialCredsIdResource**
> DeviceInitialCreds GetInitialCredsIdResource(ctx, deviceId, optional)
Returns initial password for the device

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to retrieve initial authentication credentials for | 
 **optional** | ***DeviceApiGetInitialCredsIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiGetInitialCredsIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DeviceInitialCreds**](DeviceInitialCreds.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetIpmiInfoIdResource**
> DeviceIpmiInfo GetIpmiInfoIdResource(ctx, deviceId, optional)
Returns IPMI info data

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to retrieve IPMI info. | 
 **optional** | ***DeviceApiGetIpmiInfoIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiGetIpmiInfoIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DeviceIpmiInfo**](DeviceIPMIInfo.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetIpmiInfoLoginDataResource**
> IpmiLoginData GetIpmiInfoLoginDataResource(ctx, deviceId, optional)
Returns IPMI login credentials

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to retrieve IPMI Login data. | 
 **optional** | ***DeviceApiGetIpmiInfoLoginDataResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiGetIpmiInfoLoginDataResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**IpmiLoginData**](IPMILoginData.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetIpmiThresholdsIdResource**
> DeviceIpmiThresholds GetIpmiThresholdsIdResource(ctx, deviceId, optional)
Returns IPMI thresholds data

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to View | 
 **optional** | ***DeviceApiGetIpmiThresholdsIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiGetIpmiThresholdsIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DeviceIpmiThresholds**](DeviceIPMIThresholds.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetIpmiValidLoginIdResource**
> IpmiValidLogin GetIpmiValidLoginIdResource(ctx, deviceId, optional)
Returns if device have valid credentials for IPMI login

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to check IPMI credentials | 
 **optional** | ***DeviceApiGetIpmiValidLoginIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiGetIpmiValidLoginIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**IpmiValidLogin**](IPMIValidLogin.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetPowerResource**
> DevicePower GetPowerResource(ctx, deviceId, optional)
Get device's current power status

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to View / Update | 
 **optional** | ***DeviceApiGetPowerResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiGetPowerResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DevicePower**](DevicePower.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostDeviceIpmiWhitelistResource**
> PostDeviceIpmiWhitelistResource(ctx, deviceId, payload)
Add a public IP on IPMI whitelist

Returns IPMI public IP

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of the Device to put IP in Whitelist | 
  **payload** | [**DeviceIpmiWhitelistIp**](DeviceIpmiWhitelistIp.md)|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostPowerResource**
> DevicePower PostPowerResource(ctx, deviceId, action, optional)
Apply action to device power

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to View / Update | 
  **action** | **string**| Must be one of boot|reboot|shutdown | 
 **optional** | ***DeviceApiPostPowerResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiPostPowerResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DevicePower**](DevicePower.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutClientDeviceTagOrderResource**
> DeviceTag PutClientDeviceTagOrderResource(ctx, payload, optional)
Update device tags order for current user

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **payload** | [**DeviceTag**](DeviceTag.md)|  | 
 **optional** | ***DeviceApiPutClientDeviceTagOrderResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiPutClientDeviceTagOrderResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DeviceTag**](DeviceTag.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutDeviceIdResource**
> Device PutDeviceIdResource(ctx, deviceId, payload, optional)
Updates Device MetaData for a Single Device

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to View / Update | 
  **payload** | [**DeviceUpdate**](DeviceUpdate.md)|  | 
 **optional** | ***DeviceApiPutDeviceIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiPutDeviceIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**Device**](Device.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutDeviceTagIdResource**
> DeviceTag PutDeviceTagIdResource(ctx, deviceId, payload, optional)
Update device tags

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to View / Update | 
  **payload** | [**DeviceTag**](DeviceTag.md)|  | 
 **optional** | ***DeviceApiPutDeviceTagIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiPutDeviceTagIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DeviceTag**](DeviceTag.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutIpmiDevicesThresholdsIdResource**
> DevicesIpmiThresholds PutIpmiDevicesThresholdsIdResource(ctx, payload, optional)
Updates IPMI thresholds for device list

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **payload** | [**UpdateDevicesIpmiThresholds**](UpdateDevicesIpmiThresholds.md)|  | 
 **optional** | ***DeviceApiPutIpmiDevicesThresholdsIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiPutIpmiDevicesThresholdsIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DevicesIpmiThresholds**](DevicesIPMIThresholds.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutIpmiThresholdsIdResource**
> DeviceIpmiThresholds PutIpmiThresholdsIdResource(ctx, deviceId, payload, optional)
Updates IPMI thresholds data

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to View | 
  **payload** | [**DeviceIpmiThresholds**](DeviceIpmiThresholds.md)|  | 
 **optional** | ***DeviceApiPutIpmiThresholdsIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeviceApiPutIpmiThresholdsIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**DeviceIpmiThresholds**](DeviceIPMIThresholds.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

