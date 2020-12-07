package gomemory

import (
	"fmt"
	"testing"
)

func TestGenericTranslate(t *testing.T) {
	text := "Hello World!"
	src := "en"
	dest := "ru"
	mime := "text/plain"

	resp, err := Translate(Parameters{
		Text:     text,
		Src:      src,
		Dest:     dest,
		MimeType: mime,
		Email:    "gen",
	})
	if err != nil {
		t.Error(err)
	}

	if resp.Data.Text != "Привет мир!" {
		t.Error("Incorrect translation:", resp.Data.Text)
	}
}

func TestBigTranslate(t *testing.T) {
	text := "Browse thousands of apartments for sale and condos for sale on our innovative listing platform powered by real-time updates and see only what’s available. Broker One lets you search for apartments for sale through responsive filtering data that is refreshed minute-by-minute so you never miss an opportunity to find exactly what you’re looking for based on price, location, size, condo amenities, facilities and other unique features. Discover a gorgeous high-rise condo in the heart of the city with breathtaking views or waterfront condos for sale that place you just steps away from the water and local marine activities. Whether you’re seeking a primary residence or to add an investment property to your portfolio, Broker One gives you access to all available apartments and condos for sale with up-to-date realtor contacts and scheduling."
	text += "Browsing apartments for rent often involves meeting a strict budget and schedule while accommodating specific locations, floor spaces and local amenities. Broker One is designed to simplify your search for apartments for rent and a house for rent by giving you access to real-time rental availability, detailed data filters, leasing terms and pricing that can see you moved in a matter of days. Whether you’re seeking to relocate to a new area or are in a hurry to start working at your new job, Broker One can put you in touch with local agents ready to sign on an apartment or home rental right away. Our broad listings platform is being updated minute-by-minute so you’re always just a click away from finding the perfect rental for your needs and lifestyle."
	text += "Browse thousands of listings on commercial real estate for sale and narrow down your search to the ideal property for company relocation, expansion or the development of a brand new business. Broker One is your accessory of choice in finding commercial property for sale with access to real-time purchase data aligned to your unique business terms, move-in timeline, renovation and construction requirements and company specifications with instant agent contact and appointment scheduling. Whether you are investing in a new business, relocating or expanding from out-of-state, or are a local company seeking a more accommodating location for your business, Broker One’s powerful directory platform is updated minute-by-minute and lets you browse and select commercial real estate for sale with the seamless efficiency to match your company goals."
	src := "en"
	dest := "ru"
	mime := "text/plain"

	resp, err := Translate(Parameters{
		Text:     text,
		Src:      src,
		Dest:     dest,
		MimeType: mime,
		Email:    "gen",
	})
	if err != nil {
		t.Error(err)
	}

	fmt.Println(resp.Data.Text)
}
