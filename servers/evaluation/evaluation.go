package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	//.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	url := r.URL.Query().Get("url")
	if len(url) <= 0 {
		http.Error(w, "Could not get URL", http.StatusBadRequest)
		return
	}

	body, err := fetchHTML(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	PageSummary, err := extractSummary(url, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonErr := json.NewEncoder(w).Encode(PageSummary)
	if jsonErr != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer body.Close()

}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	response, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("%v", err)
	}

	ctype := response.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, fmt.Errorf("Content is %v and not HTML", ctype)
	}

	return response.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {

	tokenizer := html.NewTokenizer(htmlStream)
	summary := &PageSummary{}
	PreviewImages := []*PreviewImage{}
	oneImage := &PreviewImage{}

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
		}

		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if "head" == token.Data {
				break
			}
		}

		if tokenType == html.SelfClosingTagToken || tokenType == html.StartTagToken {
			token := tokenizer.Token()

			if "title" == token.Data {
				tokenType := tokenizer.Next()
				if tokenType == html.TextToken {
					token := tokenizer.Token()
					if len(summary.Title) == 0 {
						summary.Title = token.Data
					}
				}
			}

			if "meta" == token.Data {
				//meta prop tags
				siteURL, isProp := loopMetaPropAttributes(token, "og:url")
				if isProp {
					summary.URL = siteURL
				}
				siteTitle, isProp := loopMetaPropAttributes(token, "og:title")
				if isProp {
					summary.Title = siteTitle
				}
				siteType, isProp := loopMetaPropAttributes(token, "og:type")
				if isProp {
					summary.Type = siteType
				}
				siteName, isProp := loopMetaPropAttributes(token, "og:site_name")
				if isProp {
					summary.SiteName = siteName
				}
				siteDesc, isProp := loopMetaPropAttributes(token, "og:description")
				if isProp {
					summary.Description = siteDesc
				}

				// meta name tags
				desc, isName := loopMetaNamePropAttributes(token, "description")
				if isName && len(summary.Description) == 0 {
					summary.Description = desc
				}
				author, isName := loopMetaNamePropAttributes(token, "author")
				if isName {
					summary.Author = author
				}
				keywords, isName := loopMetaNamePropAttributes(token, "keywords")
				if isName {
					untrimmed := strings.Split(keywords, ",")
					for i := range untrimmed {
						untrimmed[i] = strings.TrimSpace(untrimmed[i])
					}
					summary.Keywords = untrimmed
				}

				// meta name og image tags
				imageURL, isImage := loopMetaPropAttributes(token, "og:image")

				if isImage {
					oneImage = &PreviewImage{}
					PreviewImages = append(PreviewImages, oneImage)
					absURL, err := getAbsoluteURL(imageURL, pageURL)
					if err != nil {
						return nil, err
					}
					oneImage.URL = absURL
				}

				imageSecureURL, isImageSecureURL := loopMetaPropAttributes(token, "og:image:secure_url")
				if isImageSecureURL {
					oneImage.SecureURL = imageSecureURL
				}

				imageType, isImageType := loopMetaPropAttributes(token, "og:image:type")
				if isImageType {
					oneImage.Type = imageType
				}

				imageWidth, isImageWidth := loopMetaPropAttributes(token, "og:image:width")
				if isImageWidth {
					width, err := strconv.Atoi(imageWidth)
					if err != nil {
						return nil, err
					}
					oneImage.Width = width
				}

				imageHeight, isImageHeight := loopMetaPropAttributes(token, "og:image:height")
				if isImageHeight {
					height, err := strconv.Atoi(imageHeight)
					if err != nil {
						return nil, err
					}
					oneImage.Height = height
				}

				imageAlt, isImageAlt := loopMetaPropAttributes(token, "og:image:alt")
				if isImageAlt {
					oneImage.Alt = imageAlt

				}
			}

			if "link" == token.Data {
				PreviewIcon := &PreviewImage{}
				hasRel := false
				for _, attr := range token.Attr {
					if attr.Key == "rel" && attr.Val == "icon" {
						hasRel = true
					}
				}

				if hasRel {
					for _, attr := range token.Attr {
						switch attr.Key {
						case "href":
							absURLIcon, err := getAbsoluteURL(attr.Val, pageURL)
							if err != nil {
								return nil, err
							}
							PreviewIcon.URL = absURLIcon
						case "sizes":
							sizes := strings.Split(attr.Val, "x")
							if sizes[0] != "any" {
								height, err := strconv.Atoi(sizes[0])
								if err != nil {
									return nil, err
								}
								width, err := strconv.Atoi(sizes[1])
								if err != nil {
									return nil, err
								}
								PreviewIcon.Height = height
								PreviewIcon.Width = width
							}
						case "type":
							PreviewIcon.Type = attr.Val
						case "alt":
							PreviewIcon.Alt = attr.Val
						}
					}
					summary.Icon = PreviewIcon
				}
			}
		}
	}
	if len(PreviewImages) > 0 {
		summary.Images = PreviewImages
	}
	return summary, nil
}

func loopMetaPropAttributes(token html.Token, attribute string) (content string, isProp bool) {
	for _, attr := range token.Attr {
		if attr.Key == "property" && attr.Val == attribute {
			isProp = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}
	return content, isProp
}

func loopMetaNamePropAttributes(token html.Token, attribute string) (content string, isName bool) {
	for _, attr := range token.Attr {
		if attr.Key == "name" && attr.Val == attribute {
			isName = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}
	return content, isName
}

func getAbsoluteURL(relativeURL string, pageURL string) (absoluteURL string, error error) {
	u, err := url.Parse(relativeURL)
	if err != nil {
		return "", err
	}
	base, err := url.Parse(pageURL)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(u).String(), nil
}
