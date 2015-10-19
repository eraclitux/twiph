package main

import (
	"encoding/csv"
	"encoding/json"
	"html/template"
	"io"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
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
	c := csv.NewReader(r)
	// FIXME this will not work, we must intercept this in spleep function
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			ErrorLogger.Println("interrupt signal intercepted")
			return nodes, nil
		default:
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
				n := NewMemNode(profile, record[1])
				nodes = append(nodes, n)
			}
		}
	}
	return nodes, nil
}

type dNode struct {
	Name       string `json:"name"`
	Group      int    `json:"group"`
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
	names := []dNode{}
	links := []dLink{}
	for _, n := range nodes {
		gId := parseGroupId(n.Group())
		m := dNode{
			Name:       n.RealName(),
			Group:      gId,
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
