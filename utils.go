package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
)

// NewMemNode returns a memory backed node which implements
// Node interface.
// Do not use in goroutines!
func NewMemNode(data TwitterData, group string) Node {
	n := &memoryNode{
		TwitterData:     data,
		id:              counterID,
		group:           group,
		internalFriends: make([]uint, 0),
	}
	counterID++
	return n
}

// GetData searchs and retrieves Twitter accounts
// passed as csv encoded data from an io.Reader.
// Data format expected:
// Twitter Name (string), Group in graph (int).
func GetData(r io.Reader) ([]Node, error) {
	nodes := make([]Node, 0)
	scanner := bufio.NewScanner(r)
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
			return nil, err
		}
		if profile.Twid != "" {
			err = GetFriends(&profile)
			if err != nil {
				return nil, err
			}
			n := NewMemNode(profile, record[1])
			nodes = append(nodes, n)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}

type dNode struct {
	Name       string `json:"name"`
	Group      string `json:"group"`
	Avatar     string `json:"avatar"`
	ScreenName string `json:"screenName"`
}
type dLink struct {
	Source uint `json:"source"`
	Target uint `json:"target"`
	Value  uint `json:"value"`
}

func WriteData(w io.Writer, nodes []Node) error {
	names := []dNode{}
	links := []dLink{}
	for _, n := range nodes {
		m := dNode{
			Name:       n.RealName(),
			Group:      n.Group(),
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
