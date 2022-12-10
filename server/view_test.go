package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_viewApp(t *testing.T) {
	app := viewApp()
	html, clean := createIndexHTML(t)
	defer clean()

	t.Run("GET /", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res, err := app.Test(req)
		if err != nil {
			t.Error(err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}

		if string(body) != string(html) {
			t.Fatalf("want: %v , got: %v", html, body)
		}
	})

	t.Run("GET /unknown-path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/unknown-path", nil)
		res, err := app.Test(req)
		if err != nil {
			t.Error(err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}

		if string(body) != string(html) {
			t.Fatalf("want: %v , got: %v", html, body)
		}
	})
}

func createIndexHTML(t testing.TB) ([]byte, func()) {
	t.Helper()
	html := []byte(`<html>
<head><title>index.html</title></head>
<body></body>
</html>
`)
	err := os.MkdirAll("./views/dist/", 0o755)
	if err != nil {
		t.Error(err)
	}

	err = os.WriteFile("./views/dist/index.html", html, 0o755)
	if err != nil {
		t.Error(err)
	}

	rm := func() {
		err := os.RemoveAll("./views")
		if err != nil {
			t.Error(err)
		}
	}

	return html, rm
}
