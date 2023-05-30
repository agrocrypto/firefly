// Copyright © 2022 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package definitions

import (
	"context"
	"fmt"
	"testing"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly/mocks/syncasyncmocks"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDefineFFIResolveFail(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	ffi := &fftypes.FFI{
		Name:      "ffi1",
		Version:   "1.0",
		Methods:   []*fftypes.FFIMethod{{}},
		Events:    []*fftypes.FFIEvent{{}},
		Errors:    []*fftypes.FFIError{{}},
		Published: true,
	}

	ds.mcm.On("ResolveFFI", context.Background(), ffi).Return(fmt.Errorf("pop"))

	err := ds.DefineFFI(context.Background(), ffi, false)
	assert.EqualError(t, err, "pop")
}

func TestDefineFFIFail(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	ffi := &fftypes.FFI{
		Name:      "ffi1",
		Version:   "1.0",
		Published: true,
	}

	ds.mdi.On("GetFFIByNetworkName", context.Background(), "ns1", "ffi1", "1.0").Return(nil, nil)
	ds.mcm.On("ResolveFFI", context.Background(), ffi).Return(nil)
	ds.mim.On("GetMultipartyRootOrg", context.Background()).Return(nil, fmt.Errorf("pop"))

	err := ds.DefineFFI(context.Background(), ffi, false)
	assert.EqualError(t, err, "pop")
}

func TestDefineFFIExists(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	ffi := &fftypes.FFI{
		Name:      "ffi1",
		Version:   "1.0",
		Published: true,
	}

	ds.mdi.On("GetFFIByNetworkName", context.Background(), "ns1", "ffi1", "1.0").Return(&fftypes.FFI{}, nil)
	ds.mcm.On("ResolveFFI", context.Background(), ffi).Return(nil)

	err := ds.DefineFFI(context.Background(), ffi, false)
	assert.Regexp(t, "FF10448", err)
}

func TestDefineFFIQueryFail(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	ffi := &fftypes.FFI{
		Name:      "ffi1",
		Version:   "1.0",
		Published: true,
	}

	ds.mdi.On("GetFFIByNetworkName", context.Background(), "ns1", "ffi1", "1.0").Return(nil, fmt.Errorf("pop"))
	ds.mcm.On("ResolveFFI", context.Background(), ffi).Return(nil)

	err := ds.DefineFFI(context.Background(), ffi, false)
	assert.EqualError(t, err, "pop")
}

func TestDefineFFIOk(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	ffi := &fftypes.FFI{
		Name:      "ffi1",
		Version:   "1.0",
		Published: true,
	}

	ds.mdi.On("GetFFIByNetworkName", context.Background(), "ns1", "ffi1", "1.0").Return(nil, nil)
	ds.mcm.On("ResolveFFI", context.Background(), ffi).Return(nil)
	ds.mim.On("GetMultipartyRootOrg", context.Background()).Return(&core.Identity{
		IdentityBase: core.IdentityBase{
			DID: "firefly:org1",
		},
	}, nil)
	ds.mim.On("ResolveInputSigningIdentity", context.Background(), mock.Anything).Return(nil)

	mms := &syncasyncmocks.Sender{}
	ds.mbm.On("NewBroadcast", mock.Anything).Return(mms)
	mms.On("Send", context.Background()).Return(nil)

	err := ds.DefineFFI(context.Background(), ffi, false)
	assert.NoError(t, err)

	mms.AssertExpectations(t)
}

func TestDefineFFIConfirm(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	ffi := &fftypes.FFI{
		Name:      "ffi1",
		Version:   "1.0",
		Published: true,
	}

	ds.mdi.On("GetFFIByNetworkName", context.Background(), "ns1", "ffi1", "1.0").Return(nil, nil)
	ds.mcm.On("ResolveFFI", context.Background(), ffi).Return(nil)
	ds.mim.On("GetMultipartyRootOrg", context.Background()).Return(&core.Identity{
		IdentityBase: core.IdentityBase{
			DID: "firefly:org1",
		},
	}, nil)
	ds.mim.On("ResolveInputSigningIdentity", context.Background(), mock.Anything).Return(nil)

	mms := &syncasyncmocks.Sender{}
	ds.mbm.On("NewBroadcast", mock.Anything).Return(mms)
	mms.On("SendAndWait", context.Background()).Return(nil)

	err := ds.DefineFFI(context.Background(), ffi, true)
	assert.NoError(t, err)

	mms.AssertExpectations(t)
}

func TestDefineFFIPublishNonMultiparty(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = false

	ffi := &fftypes.FFI{
		Name:      "ffi1",
		Version:   "1.0",
		Published: true,
	}

	err := ds.DefineFFI(context.Background(), ffi, false)
	assert.Regexp(t, "FF10414", err)
}

func TestDefineFFINonMultiparty(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)

	ffi := &fftypes.FFI{
		Name:    "ffi1",
		Version: "1.0",
	}

	ds.mcm.On("ResolveFFI", context.Background(), ffi).Return(nil)
	ds.mdi.On("InsertOrGetFFI", context.Background(), ffi).Return(nil, nil)
	ds.mdi.On("InsertEvent", context.Background(), mock.Anything).Return(nil)

	err := ds.DefineFFI(context.Background(), ffi, false)
	assert.NoError(t, err)
}

func TestDefineFFINonMultipartyFail(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)

	ffi := &fftypes.FFI{
		Name:    "ffi1",
		Version: "1.0",
	}

	ds.mcm.On("ResolveFFI", context.Background(), ffi).Return(fmt.Errorf("pop"))

	err := ds.DefineFFI(context.Background(), ffi, false)
	assert.Regexp(t, "FF10403", err)
}

func TestDefineContractAPIResolveFail(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	url := "http://firefly"
	api := &core.ContractAPI{}

	ds.mcm.On("ResolveContractAPI", context.Background(), url, api).Return(fmt.Errorf("pop"))

	err := ds.DefineContractAPI(context.Background(), url, api, false)
	assert.EqualError(t, err, "pop")
}

func TestDefineContractAPIFail(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	url := "http://firefly"
	api := &core.ContractAPI{}

	ds.mcm.On("ResolveContractAPI", context.Background(), url, api).Return(nil)
	ds.mim.On("GetMultipartyRootOrg", context.Background()).Return(nil, fmt.Errorf("pop"))

	err := ds.DefineContractAPI(context.Background(), url, api, false)
	assert.EqualError(t, err, "pop")
}

func TestDefineContractAPIOk(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	url := "http://firefly"
	api := &core.ContractAPI{}

	ds.mcm.On("ResolveContractAPI", context.Background(), url, api).Return(nil)
	ds.mim.On("GetMultipartyRootOrg", context.Background()).Return(&core.Identity{
		IdentityBase: core.IdentityBase{
			DID: "firefly:org1",
		},
	}, nil)
	ds.mim.On("ResolveInputSigningIdentity", context.Background(), mock.Anything).Return(nil)

	mms := &syncasyncmocks.Sender{}
	ds.mbm.On("NewBroadcast", mock.Anything).Return(mms)
	mms.On("Send", context.Background()).Return(nil)

	err := ds.DefineContractAPI(context.Background(), url, api, false)
	assert.NoError(t, err)

	mms.AssertExpectations(t)
}

func TestDefineContractAPINonMultiparty(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)

	url := "http://firefly"
	api := &core.ContractAPI{}

	err := ds.DefineContractAPI(context.Background(), url, api, false)
	assert.Regexp(t, "FF10403", err)
}

func TestPublishFFI(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	mms := &syncasyncmocks.Sender{}

	ffi := &fftypes.FFI{
		Name:      "ffi1",
		Version:   "1.0",
		Namespace: "ns1",
		Published: false,
	}

	ds.mdi.On("GetFFIByNetworkName", context.Background(), "ns1", "ffi1-shared", "1.0").Return(nil, nil)
	ds.mcm.On("GetFFI", context.Background(), "ffi1", "1.0").Return(ffi, nil)
	ds.mcm.On("ResolveFFI", context.Background(), ffi).Return(nil)
	ds.mim.On("GetMultipartyRootOrg", context.Background()).Return(&core.Identity{
		IdentityBase: core.IdentityBase{
			DID: "firefly:org1",
		},
	}, nil)
	ds.mim.On("ResolveInputSigningIdentity", mock.Anything, mock.Anything).Return(nil)
	ds.mbm.On("NewBroadcast", mock.Anything).Return(mms)
	mms.On("Prepare", context.Background()).Return(nil)
	mms.On("Send", context.Background()).Return(nil)
	mockRunAsGroupPassthrough(ds.mdi)

	result, err := ds.PublishFFI(context.Background(), "ffi1", "1.0", "ffi1-shared", false)
	assert.NoError(t, err)
	assert.Equal(t, ffi, result)
	assert.True(t, ffi.Published)

	mms.AssertExpectations(t)
}

func TestPublishFFIAlreadyPublished(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	ffi := &fftypes.FFI{
		Name:      "ffi1",
		Version:   "1.0",
		Namespace: "ns1",
		Published: true,
	}

	ds.mcm.On("GetFFI", context.Background(), "ffi1", "1.0").Return(ffi, nil)
	mockRunAsGroupPassthrough(ds.mdi)

	_, err := ds.PublishFFI(context.Background(), "ffi1", "1.0", "ffi1-shared", false)
	assert.Regexp(t, "FF10450", err)
}

func TestPublishFFIQueryFail(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	ds.mcm.On("GetFFI", context.Background(), "ffi1", "1.0").Return(nil, fmt.Errorf("pop"))
	mockRunAsGroupPassthrough(ds.mdi)

	_, err := ds.PublishFFI(context.Background(), "ffi1", "1.0", "ffi1-shared", false)
	assert.EqualError(t, err, "pop")
}

func TestPublishFFIResolveFail(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	ffi := &fftypes.FFI{
		Name:      "ffi1",
		Version:   "1.0",
		Namespace: "ns1",
		Published: false,
	}

	ds.mcm.On("GetFFI", context.Background(), "ffi1", "1.0").Return(ffi, nil)
	ds.mcm.On("ResolveFFI", context.Background(), ffi).Return(fmt.Errorf("pop"))
	mockRunAsGroupPassthrough(ds.mdi)

	_, err := ds.PublishFFI(context.Background(), "ffi1", "1.0", "ffi1-shared", false)
	assert.EqualError(t, err, "pop")
}

func TestPublishFFIPrepareFail(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = true

	mms := &syncasyncmocks.Sender{}

	ffi := &fftypes.FFI{
		Name:      "ffi1",
		Version:   "1.0",
		Namespace: "ns1",
		Published: false,
	}

	ds.mdi.On("GetFFIByNetworkName", context.Background(), "ns1", "ffi1-shared", "1.0").Return(nil, nil)
	ds.mcm.On("GetFFI", context.Background(), "ffi1", "1.0").Return(ffi, nil)
	ds.mcm.On("ResolveFFI", context.Background(), ffi).Return(nil)
	ds.mim.On("GetMultipartyRootOrg", context.Background()).Return(&core.Identity{
		IdentityBase: core.IdentityBase{
			DID: "firefly:org1",
		},
	}, nil)
	ds.mim.On("ResolveInputSigningIdentity", mock.Anything, mock.Anything).Return(nil)
	ds.mbm.On("NewBroadcast", mock.Anything).Return(mms)
	mms.On("Prepare", context.Background()).Return(fmt.Errorf("pop"))
	mockRunAsGroupPassthrough(ds.mdi)

	_, err := ds.PublishFFI(context.Background(), "ffi1", "1.0", "ffi1-shared", false)
	assert.EqualError(t, err, "pop")

	mms.AssertExpectations(t)
}

func TestPublishFFINonMultiparty(t *testing.T) {
	ds := newTestDefinitionSender(t)
	defer ds.cleanup(t)
	ds.multiparty = false

	_, err := ds.PublishFFI(context.Background(), "ffi1", "1.0", "ffi1-shared", false)
	assert.Regexp(t, "FF10414", err)
}
