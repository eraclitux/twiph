// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package main

import "fmt"

var counterID uint

type Node interface {
	RealName() string
	Group() string
	GroupName() string
	LocalFriends() []uint
	MapLocal([]Node)
	// Returns Twitter screen name
	TwitterName() string
	// Returns Twitter ID
	TID() string
	// Returns http addres to profile image.
	Pic() string
	// Returns local ID
	ID() uint
}

type TwitterData struct {
	Twid       string `json:"id_str"`
	ScreenName string `json:"screen_name"`
	Name       string
	Verified   bool
	Avatar     string `json:"profile_image_url"`
	// A slice of id_str
	Friends []string
}

func (t TwitterData) String() string {
	v := "not verified"
	if t.Verified {
		v = "verified"
	}
	return fmt.Sprintf("%s, @%s, %s", t.Name, t.ScreenName, v)
}

type memoryNode struct {
	TwitterData
	// Internal id used to create json for graph
	id              uint
	groupID         string
	groupName       string
	internalFriends []uint
}

func (n *memoryNode) ID() uint {
	return n.id
}

func (n *memoryNode) TID() string {
	return n.Twid
}

func (n *memoryNode) RealName() string {
	return n.Name
}
func (n *memoryNode) TwitterName() string {
	return n.ScreenName
}

func (n *memoryNode) Group() string {
	return n.groupID
}
func (n *memoryNode) GroupName() string {
	return n.groupName
}
func (n *memoryNode) Pic() string {
	return n.Avatar
}

func (n *memoryNode) LocalFriends() []uint {
	return n.internalFriends
}

func (n *memoryNode) MapLocal(nodes []Node) {
	InfoLogger.Println("finding connections for", n)
	for _, tid := range n.Friends {
		for _, m := range nodes {
			if m.TID() == tid {
				n.internalFriends = append(n.internalFriends, m.ID())
			}
		}
	}
}
