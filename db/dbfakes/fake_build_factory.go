// Code generated by counterfeiter. DO NOT EDIT.
package dbfakes

import (
	"sync"

	"github.com/concourse/atc/db"
)

type FakeBuildFactory struct {
	BuildStub        func(int) (db.Build, bool, error)
	buildMutex       sync.RWMutex
	buildArgsForCall []struct {
		arg1 int
	}
	buildReturns struct {
		result1 db.Build
		result2 bool
		result3 error
	}
	buildReturnsOnCall map[int]struct {
		result1 db.Build
		result2 bool
		result3 error
	}
	TeamBuildsStub        func(db.Page, ...string) ([]db.Build, db.Pagination, error)
	teamBuildsMutex       sync.RWMutex
	teamBuildsArgsForCall []struct {
		arg1 db.Page
		arg2 []string
	}
	teamBuildsReturns struct {
		result1 []db.Build
		result2 db.Pagination
		result3 error
	}
	teamBuildsReturnsOnCall map[int]struct {
		result1 []db.Build
		result2 db.Pagination
		result3 error
	}
	PublicBuildsStub        func(db.Page) ([]db.Build, db.Pagination, error)
	publicBuildsMutex       sync.RWMutex
	publicBuildsArgsForCall []struct {
		arg1 db.Page
	}
	publicBuildsReturns struct {
		result1 []db.Build
		result2 db.Pagination
		result3 error
	}
	publicBuildsReturnsOnCall map[int]struct {
		result1 []db.Build
		result2 db.Pagination
		result3 error
	}
	GetAllStartedBuildsStub        func() ([]db.Build, error)
	getAllStartedBuildsMutex       sync.RWMutex
	getAllStartedBuildsArgsForCall []struct{}
	getAllStartedBuildsReturns     struct {
		result1 []db.Build
		result2 error
	}
	getAllStartedBuildsReturnsOnCall map[int]struct {
		result1 []db.Build
		result2 error
	}
	MarkNonInterceptibleBuildsStub        func() error
	markNonInterceptibleBuildsMutex       sync.RWMutex
	markNonInterceptibleBuildsArgsForCall []struct{}
	markNonInterceptibleBuildsReturns     struct {
		result1 error
	}
	markNonInterceptibleBuildsReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBuildFactory) Build(arg1 int) (db.Build, bool, error) {
	fake.buildMutex.Lock()
	ret, specificReturn := fake.buildReturnsOnCall[len(fake.buildArgsForCall)]
	fake.buildArgsForCall = append(fake.buildArgsForCall, struct {
		arg1 int
	}{arg1})
	fake.recordInvocation("Build", []interface{}{arg1})
	fake.buildMutex.Unlock()
	if fake.BuildStub != nil {
		return fake.BuildStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.buildReturns.result1, fake.buildReturns.result2, fake.buildReturns.result3
}

func (fake *FakeBuildFactory) BuildCallCount() int {
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	return len(fake.buildArgsForCall)
}

func (fake *FakeBuildFactory) BuildArgsForCall(i int) int {
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	return fake.buildArgsForCall[i].arg1
}

func (fake *FakeBuildFactory) BuildReturns(result1 db.Build, result2 bool, result3 error) {
	fake.BuildStub = nil
	fake.buildReturns = struct {
		result1 db.Build
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeBuildFactory) BuildReturnsOnCall(i int, result1 db.Build, result2 bool, result3 error) {
	fake.BuildStub = nil
	if fake.buildReturnsOnCall == nil {
		fake.buildReturnsOnCall = make(map[int]struct {
			result1 db.Build
			result2 bool
			result3 error
		})
	}
	fake.buildReturnsOnCall[i] = struct {
		result1 db.Build
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeBuildFactory) TeamBuilds(arg1 db.Page, arg2 ...string) ([]db.Build, db.Pagination, error) {
	fake.teamBuildsMutex.Lock()
	ret, specificReturn := fake.teamBuildsReturnsOnCall[len(fake.teamBuildsArgsForCall)]
	fake.teamBuildsArgsForCall = append(fake.teamBuildsArgsForCall, struct {
		arg1 db.Page
		arg2 []string
	}{arg1, arg2})
	fake.recordInvocation("TeamBuilds", []interface{}{arg1, arg2})
	fake.teamBuildsMutex.Unlock()
	if fake.TeamBuildsStub != nil {
		return fake.TeamBuildsStub(arg1, arg2...)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.teamBuildsReturns.result1, fake.teamBuildsReturns.result2, fake.teamBuildsReturns.result3
}

func (fake *FakeBuildFactory) TeamBuildsCallCount() int {
	fake.teamBuildsMutex.RLock()
	defer fake.teamBuildsMutex.RUnlock()
	return len(fake.teamBuildsArgsForCall)
}

func (fake *FakeBuildFactory) TeamBuildsArgsForCall(i int) (db.Page, []string) {
	fake.teamBuildsMutex.RLock()
	defer fake.teamBuildsMutex.RUnlock()
	return fake.teamBuildsArgsForCall[i].arg1, fake.teamBuildsArgsForCall[i].arg2
}

func (fake *FakeBuildFactory) TeamBuildsReturns(result1 []db.Build, result2 db.Pagination, result3 error) {
	fake.TeamBuildsStub = nil
	fake.teamBuildsReturns = struct {
		result1 []db.Build
		result2 db.Pagination
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeBuildFactory) TeamBuildsReturnsOnCall(i int, result1 []db.Build, result2 db.Pagination, result3 error) {
	fake.TeamBuildsStub = nil
	if fake.teamBuildsReturnsOnCall == nil {
		fake.teamBuildsReturnsOnCall = make(map[int]struct {
			result1 []db.Build
			result2 db.Pagination
			result3 error
		})
	}
	fake.teamBuildsReturnsOnCall[i] = struct {
		result1 []db.Build
		result2 db.Pagination
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeBuildFactory) PublicBuilds(arg1 db.Page) ([]db.Build, db.Pagination, error) {
	fake.publicBuildsMutex.Lock()
	ret, specificReturn := fake.publicBuildsReturnsOnCall[len(fake.publicBuildsArgsForCall)]
	fake.publicBuildsArgsForCall = append(fake.publicBuildsArgsForCall, struct {
		arg1 db.Page
	}{arg1})
	fake.recordInvocation("PublicBuilds", []interface{}{arg1})
	fake.publicBuildsMutex.Unlock()
	if fake.PublicBuildsStub != nil {
		return fake.PublicBuildsStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.publicBuildsReturns.result1, fake.publicBuildsReturns.result2, fake.publicBuildsReturns.result3
}

func (fake *FakeBuildFactory) PublicBuildsCallCount() int {
	fake.publicBuildsMutex.RLock()
	defer fake.publicBuildsMutex.RUnlock()
	return len(fake.publicBuildsArgsForCall)
}

func (fake *FakeBuildFactory) PublicBuildsArgsForCall(i int) db.Page {
	fake.publicBuildsMutex.RLock()
	defer fake.publicBuildsMutex.RUnlock()
	return fake.publicBuildsArgsForCall[i].arg1
}

func (fake *FakeBuildFactory) PublicBuildsReturns(result1 []db.Build, result2 db.Pagination, result3 error) {
	fake.PublicBuildsStub = nil
	fake.publicBuildsReturns = struct {
		result1 []db.Build
		result2 db.Pagination
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeBuildFactory) PublicBuildsReturnsOnCall(i int, result1 []db.Build, result2 db.Pagination, result3 error) {
	fake.PublicBuildsStub = nil
	if fake.publicBuildsReturnsOnCall == nil {
		fake.publicBuildsReturnsOnCall = make(map[int]struct {
			result1 []db.Build
			result2 db.Pagination
			result3 error
		})
	}
	fake.publicBuildsReturnsOnCall[i] = struct {
		result1 []db.Build
		result2 db.Pagination
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeBuildFactory) GetAllStartedBuilds() ([]db.Build, error) {
	fake.getAllStartedBuildsMutex.Lock()
	ret, specificReturn := fake.getAllStartedBuildsReturnsOnCall[len(fake.getAllStartedBuildsArgsForCall)]
	fake.getAllStartedBuildsArgsForCall = append(fake.getAllStartedBuildsArgsForCall, struct{}{})
	fake.recordInvocation("GetAllStartedBuilds", []interface{}{})
	fake.getAllStartedBuildsMutex.Unlock()
	if fake.GetAllStartedBuildsStub != nil {
		return fake.GetAllStartedBuildsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getAllStartedBuildsReturns.result1, fake.getAllStartedBuildsReturns.result2
}

func (fake *FakeBuildFactory) GetAllStartedBuildsCallCount() int {
	fake.getAllStartedBuildsMutex.RLock()
	defer fake.getAllStartedBuildsMutex.RUnlock()
	return len(fake.getAllStartedBuildsArgsForCall)
}

func (fake *FakeBuildFactory) GetAllStartedBuildsReturns(result1 []db.Build, result2 error) {
	fake.GetAllStartedBuildsStub = nil
	fake.getAllStartedBuildsReturns = struct {
		result1 []db.Build
		result2 error
	}{result1, result2}
}

func (fake *FakeBuildFactory) GetAllStartedBuildsReturnsOnCall(i int, result1 []db.Build, result2 error) {
	fake.GetAllStartedBuildsStub = nil
	if fake.getAllStartedBuildsReturnsOnCall == nil {
		fake.getAllStartedBuildsReturnsOnCall = make(map[int]struct {
			result1 []db.Build
			result2 error
		})
	}
	fake.getAllStartedBuildsReturnsOnCall[i] = struct {
		result1 []db.Build
		result2 error
	}{result1, result2}
}

func (fake *FakeBuildFactory) MarkNonInterceptibleBuilds() error {
	fake.markNonInterceptibleBuildsMutex.Lock()
	ret, specificReturn := fake.markNonInterceptibleBuildsReturnsOnCall[len(fake.markNonInterceptibleBuildsArgsForCall)]
	fake.markNonInterceptibleBuildsArgsForCall = append(fake.markNonInterceptibleBuildsArgsForCall, struct{}{})
	fake.recordInvocation("MarkNonInterceptibleBuilds", []interface{}{})
	fake.markNonInterceptibleBuildsMutex.Unlock()
	if fake.MarkNonInterceptibleBuildsStub != nil {
		return fake.MarkNonInterceptibleBuildsStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.markNonInterceptibleBuildsReturns.result1
}

func (fake *FakeBuildFactory) MarkNonInterceptibleBuildsCallCount() int {
	fake.markNonInterceptibleBuildsMutex.RLock()
	defer fake.markNonInterceptibleBuildsMutex.RUnlock()
	return len(fake.markNonInterceptibleBuildsArgsForCall)
}

func (fake *FakeBuildFactory) MarkNonInterceptibleBuildsReturns(result1 error) {
	fake.MarkNonInterceptibleBuildsStub = nil
	fake.markNonInterceptibleBuildsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildFactory) MarkNonInterceptibleBuildsReturnsOnCall(i int, result1 error) {
	fake.MarkNonInterceptibleBuildsStub = nil
	if fake.markNonInterceptibleBuildsReturnsOnCall == nil {
		fake.markNonInterceptibleBuildsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.markNonInterceptibleBuildsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	fake.teamBuildsMutex.RLock()
	defer fake.teamBuildsMutex.RUnlock()
	fake.publicBuildsMutex.RLock()
	defer fake.publicBuildsMutex.RUnlock()
	fake.getAllStartedBuildsMutex.RLock()
	defer fake.getAllStartedBuildsMutex.RUnlock()
	fake.markNonInterceptibleBuildsMutex.RLock()
	defer fake.markNonInterceptibleBuildsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBuildFactory) recordInvocation(key string, args []interface{}) {
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

var _ db.BuildFactory = new(FakeBuildFactory)
