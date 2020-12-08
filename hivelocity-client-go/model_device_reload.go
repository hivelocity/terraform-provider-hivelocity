/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type DeviceReload struct {
	Body string `json:"body,omitempty"`
	Script string `json:"script,omitempty"`
	OperatingSystemId int32 `json:"operatingSystemId"`
	ControlPanelId int32 `json:"controlPanelId,omitempty"`
}