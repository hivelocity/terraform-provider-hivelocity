/*
 * Hivelocity API for Partners
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ContactCreate struct {
	Phone string `json:"phone,omitempty"`
	Password string `json:"password,omitempty"`
	Email string `json:"email"`
	Description string `json:"description,omitempty"`
	Active int32 `json:"active"`
	FullName string `json:"fullName"`
}
