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
	ControlPanelId int32 `json:"controlPanelId,omitempty"`
	OperatingSystemId int32 `json:"operatingSystemId"`
	Script string `json:"script,omitempty"`
	Body string `json:"body,omitempty"`
	PublicSshKeyIds []int32 `json:"publicSshKeyIds,omitempty"`
}
