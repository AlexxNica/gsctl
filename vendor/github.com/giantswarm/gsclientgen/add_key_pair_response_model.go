/*
 * Giant Swarm legacy API
 *
 * Caution: This is an incomplete description of some legacy API functions.
 *
 * OpenAPI spec version: legacy
 *
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 */

package gsclientgen

type AddKeyPairResponseModel struct {
	StatusCode int32 `json:"status_code"`

	StatusText string `json:"status_text"`

	Data AddKeyPairResponseModelData `json:"data,omitempty"`
}
