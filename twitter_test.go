// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/eraclitux/cfgp"
)

func setup(t *testing.T) {
	conf = myConf{}
	err := cfgp.Parse(&conf)
	if err != nil {
		t.Fatal("unable to parse configuration", err)
	}
	SetupLoggers(os.Stderr)
}

func TestCheckRateLimits(t *testing.T) {
	setup(t)
	err := CheckRateLimits()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearchAccount(t *testing.T) {
	//setup(t)
	profile, err := SearchAccount("Andrea Masi")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("account retrieved:", profile)
}

func TestGetFriends(t *testing.T) {
	//setup(t)
	twd := &TwitterData{ScreenName: "eraclitux"}
	err := GetFriends(twd)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("friends retrieved:", len(twd.Friends))
}
