// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Amazon Software License (the "License"). You may not
// use this file except in compliance with the License. A copy of the
// License is located at
//
// http://aws.amazon.com/asl/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package executers

import (
	"github.com/aws/amazon-ssm-agent/agent/platform"
)

var instance instanceInfo = &instanceInfoImp{}

// system represents the dependency for platform
type instanceInfo interface {
	InstanceID() (string, error)
	Region() (string, error)
}

type instanceInfoImp struct{}

// InstanceID wraps platform InstanceID
func (instanceInfoImp) InstanceID() (string, error) { return platform.InstanceID() }

// Region wraps platform Region
func (instanceInfoImp) Region() (string, error) { return platform.Region() }
