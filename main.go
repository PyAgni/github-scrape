package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "https://www.github.com"
)

var (
	username = "pyagni"
	password = ""
)

type App struct {
	Client *http.Client
}

type AuthToken struct {
	Token string
}

type TimeStamp struct {
	Stamp        string
	Stamp_Secret string
}

type Project struct {
	Name        string
	Link        string
	Description string
}

func (app *App) getTokens() (TimeStamp, AuthToken) {
	loginUrl := baseURL + "/session"
	client := app.Client

	res, err := client.Get(loginUrl)

	if err != nil {
		log.Fatalln("Error fetching response! ", err)
	}

	defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal("Error loading res body. ", err)
	}

	stampSecret, _ := document.Find("input[name='timestamp_secret']").Attr("value")
	stamp, _ := document.Find("input[name='timestamp']").Attr("value")
	token, _ := document.Find("input[name='authenticity_token']").Attr("value")

	timeStamp := TimeStamp{
		Stamp:        stamp,
		Stamp_Secret: stampSecret,
	}

	authToken := AuthToken{
		Token: token,
	}

	return timeStamp, authToken

}

func (app *App) login() {
	client := app.Client

	timeStamp, token := app.getTokens()

	loginUrl := baseURL + "/session"

	data := url.Values{
		"login":                   {username},
		"password":                {password},
		"webauthn-support":        {"supported"},
		"webauthn-iuvpaa-support": {"unsupported"},
		"required_field_d814":     {},
		"timestamp":               {timeStamp.Stamp},
		"timestamp_secret":        {timeStamp.Stamp},
		"authenticity_token":      {token.Token},
	}

	res, err := client.PostForm(loginUrl, data)

	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

}

func (app *App) getProjects() {
	projectsURL := baseURL + "/PyAgni?tab=repositories"
	client := app.Client

	res, err := client.Get(projectsURL)
	if err != nil {
		log.Fatalln("Error fetching response! ", err)
	}

	defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln("Error loading res Body! ", err)
	}

	document.Find("a[itemprop]").Each(func(i int, s *goquery.Selection) {

		link, _ := s.Attr("href")
		name := s.Text()

		fmt.Println(i, name, link)
	})
	document.Find("p[itemprop]").Each(func(i int, des *goquery.Selection) {
		description := des.Text()
		fmt.Println(i, description)
	})

}

func main() {

	jar, _ := cookiejar.New(nil)

	app := App{
		Client: &http.Client{Jar: jar},
	}

	token, auth := app.getTokens()
	app.login()
	app.getProjects()

	fmt.Println(token.Stamp, auth.Token)
}
