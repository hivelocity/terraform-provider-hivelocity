/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type CancellationCreate struct {
	DeviceId int32 `json:"deviceId"`
	Reason string `json:"reason"`
	ServiceId int32 `json:"serviceId"`
}
