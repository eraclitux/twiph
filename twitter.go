package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/mrjones/oauth"
)

var twitterEndPoint string = "https://api.twitter.com/1.1"

func CallTwitter(ep string, query map[string]string) (io.ReadCloser, error) {
	consumer := oauth.NewConsumer(conf.ConsumerKey, conf.ConsumerSecret, oauth.ServiceProvider{})
	//consumer.Debug(true)
	accessToken := &oauth.AccessToken{
		Token:  conf.AccessToken,
		Secret: conf.AccessTokenSecret,
	}
	//	client, err := consumer.MakeHttpClient(accessToken)
	//	if err != nil {
	//		return t, err
	//	}
	//	query := url.QueryEscape("q=" + s + "&page=1&count=20")
	//	ep := fmt.Sprintf("%s/users/search.json?%s", twitterEndPoint, query)
	//	response, err := client.Get(ep)
	// FIXME this is deprecated but MakeHttpClient makes a bad
	// auth because do not encode "&" in "&amp"
	response, err := consumer.Get(twitterEndPoint+ep, query, accessToken)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		b, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("%s, %s", http.StatusText(response.StatusCode), string(b))
	}
	return response.Body, nil
}

type Resource struct {
	Limit     int   `json:"limit"`
	Remaining int   `json:"remaining"`
	Reset     int64 `json:"reset"`
}
type limits struct {
	Resources map[string]map[string]Resource `json:"resources"`
}

func sleepMin(l limits) {
	var min int
	var reset int64
	resource := ""
	min = math.MaxInt32
	for _, v := range l.Resources {
		for k, t := range v {
			if t.Remaining <= min {
				resource = k
				min = t.Remaining
				reset = t.Reset
			}
		}
	}
	if min <= 1 {
		// FIXME
		//signalChan := make(chan os.Signal, 1)
		//signal.Notify(signalChan, os.Interrupt)
		//signal.Notify(signalChan, syscall.SIGTERM)
		//select {
		//case <-signalChan:
		//	ErrorLogger.Println("interrupt signal intercepted")
		//	return nodes, nil
		//default:
		secs := reset - time.Now().Unix()
		InfoLogger.Printf("reached limit for \"%s\", sleeping %d seconds", resource, secs)
		// FIXME print a count down
		time.Sleep(time.Duration(secs) * time.Second)
	}
}

// CheckRateLimits checks is some rate limit is
// aproaching and sleeps until it is resetted.
func CheckRateLimits() error {
	query := map[string]string{"resources": "users,friends,application"}
	r, err := CallTwitter("/application/rate_limit_status.json", query)
	if err != nil {
		return err
	}
	defer r.Close()
	data := limits{}
	err = json.NewDecoder(r).Decode(&data)
	if err != nil {
		return err
	}
	sleepMin(data)
	return nil
}

// SearchAccount needs user auth to call REST end point.
func SearchAccount(s string) (TwitterData, error) {
	InfoLogger.Println("searching for", s)
	t := TwitterData{}
	err := CheckRateLimits()
	if err != nil {
		return t, err
	}
	query := map[string]string{"q": s}
	r, err := CallTwitter("/users/search.json", query)
	if err != nil {
		return t, err
	}
	defer r.Close()
	data := []TwitterData{}
	err = json.NewDecoder(r).Decode(&data)
	if err != nil {
		return t, err
	}
	// FIXME return verified only if any.
	if len(data) == 0 {
		InfoLogger.Println("not found")
		return t, nil
	}
	InfoLogger.Println("found:", data[0])
	return data[0], nil
}

// GetFriends retrieves friends from Twitter.
func GetFriends(d *TwitterData) error {
	// FIXME iterate over cursors.
	err := CheckRateLimits()
	if err != nil {
		return err
	}
	query := map[string]string{"screen_name": d.ScreenName, "cursor": "-1", "stringify_ids": "true", "count": "5000"}
	r, err := CallTwitter("/friends/ids.json", query)
	if err != nil {
		return err
	}
	defer r.Close()
	data := struct {
		Ids        []string `json:"ids"`
		NextCursor string   `json:"next_cursor_str"`
	}{}
	err = json.NewDecoder(r).Decode(&data)
	if err != nil {
		return err
	}
	d.Friends = data.Ids
	return nil
}

func GetProfile(s string) (TwitterData, error) {
	t, err := SearchAccount(s)
	if err != nil {
		return t, err
	}
	return t, nil
}
