// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build mage

package main

import (
	// mage:import precommit
	_ "github.com/mesosphere/daggers-for-dkp/catalog/precommit"

	// mage:import test
	_ "github.com/mesosphere/daggers-for-dkp/catalog/gotest"
)
