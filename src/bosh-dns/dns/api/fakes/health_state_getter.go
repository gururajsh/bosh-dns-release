// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"bosh-dns/dns/api"
	"sync"
)

type FakeHealthStateGetter struct {
	HealthStateStringStub        func(string) string
	healthStateStringMutex       sync.RWMutex
	healthStateStringArgsForCall []struct {
		arg1 string
	}
	healthStateStringReturns struct {
		result1 string
	}
	healthStateStringReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeHealthStateGetter) HealthStateString(arg1 string) string {
	fake.healthStateStringMutex.Lock()
	ret, specificReturn := fake.healthStateStringReturnsOnCall[len(fake.healthStateStringArgsForCall)]
	fake.healthStateStringArgsForCall = append(fake.healthStateStringArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("HealthStateString", []interface{}{arg1})
	fake.healthStateStringMutex.Unlock()
	if fake.HealthStateStringStub != nil {
		return fake.HealthStateStringStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.healthStateStringReturns
	return fakeReturns.result1
}

func (fake *FakeHealthStateGetter) HealthStateStringCallCount() int {
	fake.healthStateStringMutex.RLock()
	defer fake.healthStateStringMutex.RUnlock()
	return len(fake.healthStateStringArgsForCall)
}

func (fake *FakeHealthStateGetter) HealthStateStringCalls(stub func(string) string) {
	fake.healthStateStringMutex.Lock()
	defer fake.healthStateStringMutex.Unlock()
	fake.HealthStateStringStub = stub
}

func (fake *FakeHealthStateGetter) HealthStateStringArgsForCall(i int) string {
	fake.healthStateStringMutex.RLock()
	defer fake.healthStateStringMutex.RUnlock()
	argsForCall := fake.healthStateStringArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeHealthStateGetter) HealthStateStringReturns(result1 string) {
	fake.healthStateStringMutex.Lock()
	defer fake.healthStateStringMutex.Unlock()
	fake.HealthStateStringStub = nil
	fake.healthStateStringReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeHealthStateGetter) HealthStateStringReturnsOnCall(i int, result1 string) {
	fake.healthStateStringMutex.Lock()
	defer fake.healthStateStringMutex.Unlock()
	fake.HealthStateStringStub = nil
	if fake.healthStateStringReturnsOnCall == nil {
		fake.healthStateStringReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.healthStateStringReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeHealthStateGetter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.healthStateStringMutex.RLock()
	defer fake.healthStateStringMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeHealthStateGetter) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ api.HealthStateGetter = new(FakeHealthStateGetter)
