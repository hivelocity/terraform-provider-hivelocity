/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Permission struct {
	Permissions []string `json:"permissions"`
	UserId int32 `json:"userId"`
}
