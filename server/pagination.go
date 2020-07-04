package server

import (
	"calories-counter/common"
	"github.com/gin-gonic/gin"
	"math"
	"net/url"
	"strconv"
	"strings"
)

type Link struct {
	Rel  string  `json:"rel"`
	Href *string `json:"href"`
}

func PageParams(c *gin.Context) (page int, perPage int, filter string) {
	page, _ = strconv.Atoi(c.Query("page"))
	perPage, _ = strconv.Atoi(c.Query("per_page"))
	if perPage > 100 {
		perPage = 100
	} else if perPage <= 0 {
		perPage = 10
	}

	return page, perPage, c.Query("filter")
}

func CreateLinks(c *gin.Context, total, page, perPage int) []Link {
	params := url.Values{}
	if c.Query("page") == "" {
		params.Add("page", strconv.Itoa(page))
	}
	if c.Query("per_page") == "" {
		params.Add("per_page", strconv.Itoa(perPage))
	}
	param := params.Encode()

	query := c.Request.URL.Query()
	appUrl := c.Request.Host + c.Request.URL.RequestURI()
	if len(query) == 0 {
		appUrl = appUrl + "?" + param
	} else if param != "" {
		appUrl = appUrl + "&" + param
	}

	var links []Link
	lastPage := 0
	if total > 0 {
		lastPage = int(math.Ceil(float64(total)/float64(perPage))) - 1
	}

	first := strings.Replace(appUrl, "page="+strconv.Itoa(page), "page=0", 1)

	var prev *string
	if page > 0 {
		prev = common.String(strings.Replace(appUrl, "page="+strconv.Itoa(page), "page="+strconv.Itoa(page-1), 1))
	}

	var next *string
	if page < lastPage {
		next = common.String(strings.Replace(appUrl, "page="+strconv.Itoa(page), "page="+strconv.Itoa(page+1), 1))
	}

	last := strings.Replace(appUrl, "page="+strconv.Itoa(page), "page="+strconv.Itoa(lastPage), 1)

	links = []Link{
		{Rel: "self", Href: &appUrl},
		{Rel: "first", Href: &first},
		{Rel: "prev", Href: prev},
		{Rel: "next", Href: next},
		{Rel: "last", Href: &last},
	}

	return links
}
