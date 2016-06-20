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

// +build windows

// Package platform contains platform specific utilities.
package platform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/amazon-ssm-agent/agent/appconfig"
	"github.com/aws/amazon-ssm-agent/agent/log"
)

const caption = "Caption"
const version = "Version"
const sku = "OperatingSystemSKU"

func getPlatformName(log log.T) (value string, err error) {
	return getPlatformDetails(caption, log)
}

func getPlatformVersion(log log.T) (value string, err error) {
	return getPlatformDetails(version, log)
}

func getPlatformSku(log log.T) (value string, err error) {
	return getPlatformDetails(sku, log)
}

func getPlatformDetails(property string, log log.T) (value string, err error) {
	log.Debugf(gettingPlatformDetailsMessage)
	value = notAvailableMessage

	cmdName := "wmic"
	cmdArgs := []string{"OS", "get", property, "/format:list"}
	var cmdOut []byte
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		log.Debugf("There was an error running %v %v, err:%v", cmdName, cmdArgs, err)
		return
	}

	// Stringnize cmd output and trim spaces 
	value = strings.TrimSpace(string(cmdOut))

	// Match whitespaces between property and = sign and remove whitespaces
	rp := regexp.MustCompile(fmt.Sprintf("%v(\\s*)%v", property, "="))
	value = rp.ReplaceAllString(value, "")

	// Trim spaces again
	value = strings.TrimSpace(value)

	log.Debugf(commandOutputMessage, value)
	return
}

var wmicCommand = filepath.Join(appconfig.EnvWinDir, "System32", "wbem", "wmic.exe")

// fullyQualifiedDomainName returns the Fully Qualified Domain Name of the instance, otherwise the hostname
func fullyQualifiedDomainName() string {
	hostName, _ := os.Hostname()

	dnsHostName := getWMICComputerSystemValue("DNSHostName")
	domainName := getWMICComputerSystemValue("Domain")

	if dnsHostName == "" || domainName == "" {
		return hostName
	}

	return dnsHostName + "." + domainName
}

// getWMICComputerSystemValue return the value part of the wmic computersystem command for the specified attribute
func getWMICComputerSystemValue(attribute string) string {
	if contentBytes, err := exec.Command(wmicCommand, "computersystem", "get", attribute, "/value").Output(); err == nil {
		contents := string(contentBytes)
		data := strings.Split(contents, "=")
		if len(data) > 1 {
			return strings.TrimSpace(data[1])
		}
	}
	return ""
}
