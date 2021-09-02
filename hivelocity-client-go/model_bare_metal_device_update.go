/*
 * Hivelocity API for Partners
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type BareMetalDeviceUpdate struct {
	Script string `json:"script,omitempty"`
	PublicSshKeyId int32 `json:"publicSshKeyId,omitempty"`
	Hostname string `json:"hostname"`
	OsName string `json:"osName"`
	Tags []string `json:"tags,omitempty"`
}
