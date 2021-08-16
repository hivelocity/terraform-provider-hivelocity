/*
 * Hivelocity API for Partners
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type DeviceInitialCreds struct {
	Port int32 `json:"port,omitempty"`
	Password string `json:"password,omitempty"`
	LockerUrl string `json:"lockerUrl,omitempty"`
	User string `json:"user,omitempty"`
}
