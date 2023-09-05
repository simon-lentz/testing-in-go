package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func getText(node *html.Node) string {
	var nodes []string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		switch child.Type {
		case html.TextNode:
			nodes = append(nodes, strings.TrimSpace(child.Data))
		case html.ElementNode:
			getText(child)
		}
	}
	return strings.Join(nodes, " ")
}

func findNodes(body, desiredTag, desiredClass string) ([]string, error) {
	// Adapted from html package examples.
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	var nodes []string
	var find func(node *html.Node)
	find = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == desiredTag {
			for _, attr := range node.Attr {
				if attr.Key == "class" {
					classes := strings.Fields(attr.Val)
					for _, foundClass := range classes {
						if foundClass == desiredClass {
							nodes = append(nodes, getText(node))
						}
					}
					break
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			find(child)
		}
	}
	find(doc)
	return nodes, nil
}

func TestApp(t *testing.T) {
	app := &App{}

	type checkFun func(r *http.Response, body string) error
	hasNoAlerts := func() checkFun {
		return func(r *http.Response, body string) error {
			nodes, err := findNodes(body, "div", "alert")
			if err != nil {
				return fmt.Errorf("findNodes() err = %s", err)
			}
			if len(nodes) != 0 {
				return fmt.Errorf("len(alerts)=%d; want len(alerts)=0", len(nodes))
			}
			return nil
		}
	}
	hasAlert := func(msg string) checkFun {
		return func(r *http.Response, body string) error {
			nodes, err := findNodes(body, "div", "alert")
			if err != nil {
				return fmt.Errorf("findNodes() err = %s", err)
			}
			for _, node := range nodes {
				if node == msg {
					return nil
				}
			}
			return fmt.Errorf("missing alert: %q", msg)
		}
	}

	tests := []struct {
		method string
		path   string
		body   io.Reader
		checks []checkFun
	}{
		{http.MethodGet, "/", nil, []checkFun{hasNoAlerts()}},
		{http.MethodGet, "/alert", nil, []checkFun{hasAlert("Something is Wrong!")}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s %s", tt.method, tt.path),
			func(t *testing.T) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(tt.method, tt.path, nil)
				if err != nil {
					t.Fatalf("http.NewRequest() err = %s", err)
				}
				app.ServeHTTP(w, r)
				res := w.Result()
				var sb strings.Builder
				if _, err = io.Copy(&sb, res.Body); err != nil {
					t.Fatalf("io.Copy() err = %s", err)
				}
				for _, check := range tt.checks {
					if err := check(res, sb.String()); err != nil {
						t.Error(err)
					}
				}
			})
	}
}
