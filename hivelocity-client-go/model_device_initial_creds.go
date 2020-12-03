/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type DeviceInitialCreds struct {
	Password string `json:"password,omitempty"`
	User string `json:"user,omitempty"`
	LockerUrl string `json:"lockerUrl,omitempty"`
	Port int32 `json:"port,omitempty"`
}
