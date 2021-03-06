// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

// Package twitter installs a faculty for access to the Twitter API.
package twitter

import (
	// "fmt"
	"log"
	"net/url"
	"strconv"
	"sync"

	"github.com/gocircuit/escher/faculty"
	"github.com/gocircuit/escher/kit/plumb"
	"github.com/gocircuit/escher/be"
	. "github.com/gocircuit/escher/circuit"

	"github.com/gocircuit/escher/kit/github.com/ChimeraCoder/anaconda"
)

func init() {
	faculty.Register("web.twitter.Client", Client{})
	faculty.Register("web.twitter.Answer", AnswerMaterializer{})
	faculty.Register("web.twitter.Consumer", ConsumerMaterializer{})
	faculty.Register("web.twitter.Access", AccessMaterializer{})
	faculty.Register("web.twitter.UserTimelineQuery", UserTimelineQueryMaterializer{})
}

// Client ...
type Client struct{}

func (Client) Materialize() (be.Reflex, Value) {
	var c1, a1 sync.Once
	api, consumer, access := make(chan *anaconda.TwitterApi), make(chan Circuit, 1), make(chan Circuit, 1)
	go func() { // start connecting monad
		var c, a Circuit
		for i := 0; i < 2; i++ {
			select {
			case c = <-consumer:
				consumer = nil
			case a = <-access:
				access = nil
			}
		}
		anaconda.SetConsumerKey(c.String("Key")) // dial API server
		anaconda.SetConsumerSecret(c.String("Secret"))
		y := anaconda.NewTwitterApi(a.String("Token"), a.String("Secret"))	
		for {
			api <- y // give out api server to all endpoint goroutines
		}
	}()
	// API
	query := make(chan Circuit, 5)
	reflex, eye := be.NewEyeCognizer(
		func (eye *be.Eye, valve string, value interface{}) {
			switch valve {
			case "Consumer":
				c1.Do(func () {consumer <- value.(Circuit)})
			case "Access":
				a1.Do(func() {access <- value.(Circuit)})
			case "UserTimelineQuery", "HomeTimelineQuery", "RetweetsQuery","RetweetsOfMeQuery":
				valve = valve[:len(valve)-len("Query")]
				query <- Make().Grow(valve, value)
			default:
				log.Printf("Unknown Twitter query: %s", valve)
 			}
		},
		"Consumer", "Access", // set to start connection
		"UserTimelineQuery", "UserTimelineResult", // UserTimeline
		"HomeTimelineQuery", "HomeTimelineResult", // HomeTimeline
		"RetweetsQuery", "RetweetsResult", // Retweets
		"RetweetsOfMeQuery", "RetweetsOfMeResult", // RetweetsOfMe
	)
	for i := 0; i < 3; i++ {
		go func() { // API response loop
			y := <-api
			for {
				g := <-query
				q := g.Letters()[0]
				x := g[q].(Circuit)
				uv := urlize(x)
				log.Printf("Twitter %s query %v", q, ImagineWithMaps(uv).(Circuit).PrintLine())
				var tweets []anaconda.Tweet
				var err error
				switch q {
				case "UserTimeline":
					tweets, err = y.GetUserTimeline(uv)
				case "HomeTimeline":
					tweets, err = y.GetHomeTimeline(uv)
				case "Retweets": 
					tweets, err = y.GetRetweets(int64(x.Int("Id")), uv)
				case "RetweetsOfMe":
					tweets, err = y.GetRetweetsOfMe(uv)
				}
				if err != nil {
					log.Fatalf("Problem %s query on Twitter (%v)", q, err)
				}
				eye.Show(
					q,
					Pretty(
						Circuit{
							"Name": x.Interface("Name"),
							"Sentence": Imagine(tweets),
						},
					),
				)
			}
		}()
	}
	return reflex, Client{}
}

func urlize(g Circuit) url.Values {
	uv := url.Values{}
	uv.Set("user_id", g.OptionalString("UserId"))
	uv.Set("screen_name", g.OptionalString("ScreenName"))
	if g.Has("AfterId") {
		uv.Set("since_id", strconv.Itoa(g.OptionalInt("AfterId"))) // return results indexed greater than AfterId
	}
	if g.Has("NotAfterId") {
		uv.Set("max_id", strconv.Itoa(g.OptionalInt("NotAfterId"))) // return results indexed no greater than NotAfterId
	}
	if g.Has("Count") {
		uv.Set("count", strconv.Itoa(g.OptionalInt("Count")))
	}
	if g.Has("TrimUser") {
		uv.Set("trim_user", "1")
	}	
	if g.Has("ExcludeReplies") {
		uv.Set("exclude_replies", "1")
	}	
	if g.Has("ContributorDetails") {
		uv.Set("contributor_details", "1")
	}	
	if g.Has("IncludeEntities") {
		uv.Set("include_entities", "1")
	}	
	return uv
}
