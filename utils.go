// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package main

import (
	"encoding/csv"
	"encoding/json"
	"html/template"
	"io"
	"regexp"
	"strconv"
)

func createIndexWithAvatar(w io.Writer) error {
	t, err := template.ParseFiles("templates/index_avatar.html")
	if err != nil {
		return err
	}
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}
	return nil
}

func createIndexWithGroups(w io.Writer) error {
	t, err := template.ParseFiles("templates/index_groups.html")
	if err != nil {
		return err
	}
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}
	return nil
}

// NewMemNode returns a memory backed node which implements
// Node interface.
// Do not use in goroutines!
func NewMemNode(data TwitterData, groupID, groupName string) Node {
	n := &memoryNode{
		TwitterData:     data,
		id:              counterID,
		groupID:         groupID,
		groupName:       groupName,
		internalFriends: make([]uint, 0),
	}
	counterID++
	return n
}

// GetData searchs and retrieves Twitter accounts
// passed as csv encoded data from an io.Reader.
// Data format expected:
// Twitter Name (string), Group id in graph (int), Group name in graph
func GetData(r io.Reader) ([]Node, error) {
	nodes := make([]Node, 0)
	c := csv.NewReader(r)
	for {
		record, err := c.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		profile, err := GetProfile(record[0])
		if err != nil {
			ErrorLogger.Println("calling remote APIs:", err)
			continue
			//return nil, err
		}
		if profile.Twid != "" {
			err = GetFriends(&profile)
			if err != nil {
				ErrorLogger.Println("calling remote APIs:", err)
				continue
				//return nil, err
			}
			n := NewMemNode(profile, record[1], record[2])
			nodes = append(nodes, n)
		}
	}
	return nodes, nil
}

type dNode struct {
	Name       string `json:"name"`
	Group      int    `json:"group"`
	GroupName  string `json:"group_name"`
	Avatar     string `json:"avatar"`
	ScreenName string `json:"screenName"`
}
type dLink struct {
	Source uint `json:"source"`
	Target uint `json:"target"`
	Value  uint `json:"value"`
}

func parseGroupId(g string) int {
	rgxp := regexp.MustCompile(`[[:space:]]`)
	g = rgxp.ReplaceAllString(g, "")
	i, err := strconv.ParseInt(g, 10, 32)
	if err != nil {
		ErrorLogger.Println("converting to int, will be 0:", g)
		return 1
	}
	return int(i)
}
func WriteData(w io.Writer, nodes []Node) error {
	InfoLogger.Println("writing data...")
	names := []dNode{}
	links := []dLink{}
	for _, n := range nodes {
		gId := parseGroupId(n.Group())
		m := dNode{
			Name:       n.RealName(),
			Group:      gId,
			GroupName:  n.GroupName(),
			Avatar:     n.Pic(),
			ScreenName: n.TwitterName(),
		}
		names = append(names, m)
		n.MapLocal(nodes)
		for _, f := range n.LocalFriends() {
			l := dLink{
				Source: n.ID(),
				Target: f,
				Value:  1,
			}
			links = append(links, l)
		}
	}
	values := map[string]interface{}{
		"nodes": names,
		"links": links,
	}
	err := json.NewEncoder(w).Encode(values)
	if err != nil {
		return err
	}
	return nil
}
