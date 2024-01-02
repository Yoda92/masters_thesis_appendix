/*
Wasp API

REST API for the Wasp node

API version: 0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"encoding/json"
)

// checks if the BlockReceiptError type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BlockReceiptError{}

// BlockReceiptError struct for BlockReceiptError
type BlockReceiptError struct {
	ErrorMessage string `json:"errorMessage"`
}

// NewBlockReceiptError instantiates a new BlockReceiptError object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBlockReceiptError(errorMessage string) *BlockReceiptError {
	this := BlockReceiptError{}
	this.ErrorMessage = errorMessage
	return &this
}

// NewBlockReceiptErrorWithDefaults instantiates a new BlockReceiptError object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBlockReceiptErrorWithDefaults() *BlockReceiptError {
	this := BlockReceiptError{}
	return &this
}

// GetErrorMessage returns the ErrorMessage field value
func (o *BlockReceiptError) GetErrorMessage() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ErrorMessage
}

// GetErrorMessageOk returns a tuple with the ErrorMessage field value
// and a boolean to check if the value has been set.
func (o *BlockReceiptError) GetErrorMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ErrorMessage, true
}

// SetErrorMessage sets field value
func (o *BlockReceiptError) SetErrorMessage(v string) {
	o.ErrorMessage = v
}

func (o BlockReceiptError) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BlockReceiptError) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["errorMessage"] = o.ErrorMessage
	return toSerialize, nil
}

type NullableBlockReceiptError struct {
	value *BlockReceiptError
	isSet bool
}

func (v NullableBlockReceiptError) Get() *BlockReceiptError {
	return v.value
}

func (v *NullableBlockReceiptError) Set(val *BlockReceiptError) {
	v.value = val
	v.isSet = true
}

func (v NullableBlockReceiptError) IsSet() bool {
	return v.isSet
}

func (v *NullableBlockReceiptError) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBlockReceiptError(val *BlockReceiptError) *NullableBlockReceiptError {
	return &NullableBlockReceiptError{value: val, isSet: true}
}

func (v NullableBlockReceiptError) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBlockReceiptError) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

