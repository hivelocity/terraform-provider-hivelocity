/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ContactCreate struct {
	Email string `json:"email"`
	Phone string `json:"phone,omitempty"`
	Active int32 `json:"active"`
	Password string `json:"password,omitempty"`
	Description string `json:"description,omitempty"`
	FullName string `json:"fullName"`
}
