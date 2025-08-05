// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package content

import (
	routing "github.com/rbenatti8/ozzo-routing/v2"
	"golang.org/x/text/language"
	"net/http"
)

// Language is the key used to store and retrieve the chosen language in routing.Context
const Language = "Language"

// LanguageNegotiator returns a content language negotiation handler.
//
// The method takes a list of languages (locale IDs) that are supported by the application.
// The negotiator will determine the best language to use by checking the Accept-Language request header.
// If no match is found, the first language will be used.
//
// In a handler, you can access the chosen language through routing.Context like the following:
//
//	func(c *routing.Context) error {
//	    language := c.Get(content.Language).(string)
//	}
//
// If you do not specify languages, the negotiator will set the language to be "en-US".
func LanguageNegotiator(languages ...string) routing.Handler {
	if len(languages) == 0 {
		languages = []string{"en-US"}
	}
	defaultLanguage := languages[0]

	return func(c *routing.Context) error {
		l := negotiateLanguage(c.Request, languages, defaultLanguage)
		c.Set(Language, l)
		return nil
	}
}

// negotiateLanguage negotiates the acceptable language according to the Accept-Language HTTP header.
func negotiateLanguage(r *http.Request, offers []string, defaultOffer string) string {
	accept := r.Header.Get("Accept-Language")

	// Parse the Accept-Language header
	tags, qs, err := language.ParseAcceptLanguage(accept)
	if err != nil || len(tags) == 0 {
		return defaultOffer
	}

	bestOffer := defaultOffer
	bestQ := float32(-1.0)

	for _, offer := range offers {
		offerTag := language.Make(offer)

		for i, tag := range tags {
			if tag == offerTag && qs[i] > bestQ {
				bestQ = qs[i]
				bestOffer = offer
			}
		}
	}

	if bestQ == 0 {
		return defaultOffer
	}

	return bestOffer
}
