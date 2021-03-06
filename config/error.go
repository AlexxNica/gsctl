package config

import "github.com/juju/errgo"

// endpointNotDefinedError means the user tries to use an endpoint that is not defined
var endpointNotDefinedError = errgo.New("the given endpoint is not defined")

// IsEndpointNotDefinedError asserts endpointNotDefinedError.
func IsEndpointNotDefinedError(err error) bool {
	return errgo.Cause(err) == endpointNotDefinedError
}

// aliasMustBeUniqueError should be used if the user tries to add an alias to
// an endpoint, but the alias is already in use
var aliasMustBeUniqueError = errgo.New("alias must be unique")

// IsAliasMustBeUniqueError asserts aliasMustBeUniqueError.
func IsAliasMustBeUniqueError(err error) bool {
	return errgo.Cause(err) == aliasMustBeUniqueError
}

// credentialsRequiredError means an attempt to store incomplete credentials in the config
var credentialsRequiredError = errgo.New("email and password must not be empty")

// IsCredentialsRequiredError asserts credentialsRequiredError.
func IsCredentialsRequiredError(err error) bool {
	return errgo.Cause(err) == credentialsRequiredError
}

var garbageCollectionFailedError = errgo.New("garbage collection failed")

// IsGarbageCollectionFailedError asserts garbageCollectionFailedError.
func IsGarbageCollectionFailedError(err error) bool {
	return errgo.Cause(err) == garbageCollectionFailedError
}

var garbageCollectionPartiallyFailedError = errgo.New("garbage collection partially failed")

// IsGarbageCollectionPartiallyFailedError asserts garbageCollectionPartiallyFailedError.
func IsGarbageCollectionPartiallyFailedError(err error) bool {
	return errgo.Cause(err) == garbageCollectionPartiallyFailedError
}
