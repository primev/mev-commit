package apiserver_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/primev/mev-commit/p2p/pkg/apiserver"
)

type testHandler struct {
	called bool
}

func (h *testHandler) Handle(_ http.ResponseWriter, _ *http.Request) {
	h.called = true
}

func TestMethodHandler(t *testing.T) {
	t.Parallel()

	t.Run("method not allowed", func(t *testing.T) {
		h := &testHandler{}
		mh := apiserver.MethodHandler("GET", http.HandlerFunc(h.Handle))

		r, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		mh.ServeHTTP(w, r)

		if w.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}

		if h.called {
			t.Fatal("handler should not have been called")
		}
	})

	t.Run("method allowed", func(t *testing.T) {
		h := &testHandler{}
		mh := apiserver.MethodHandler("GET", http.HandlerFunc(h.Handle))

		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		mh.ServeHTTP(w, r)

		if !h.called {
			t.Fatal("handler should have been called")
		}
	})
}

func TestBindJSON(t *testing.T) {
	t.Parallel()

	t.Run("bad request", func(t *testing.T) {
		type v struct {
			Foo string `json:"foo"`
		}

		r, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()

		if _, err := apiserver.BindJSON[v](w, r); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("ok", func(t *testing.T) {
		type v struct {
			Foo string `json:"foo"`
		}

		b := bytes.NewBuffer([]byte(`{"foo":"bar"}`))

		r, err := http.NewRequest("POST", "/", b)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()

		vv, err := apiserver.BindJSON[v](w, r)
		if err != nil {
			t.Fatal(err)
		}

		if vv.Foo != "bar" {
			t.Fatalf("expected foo to be %q, got %q", "bar", vv.Foo)
		}
	})
}

func TestWriteResponse(t *testing.T) {
	t.Parallel()

	t.Run("string", func(t *testing.T) {
		w := httptest.NewRecorder()

		if err := apiserver.WriteResponse(w, http.StatusOK, "foo"); err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
		}

		resp := apiserver.StatusResponse{
			Code:    http.StatusOK,
			Message: "foo",
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(resp); err != nil {
			t.Fatal(err)
		}
		buf.WriteByte('\n')

		if !bytes.Equal(w.Body.Bytes(), buf.Bytes()) {
			t.Fatalf("expected body %q, got %q", buf.String(), w.Body.String())
		}
	})

	t.Run("struct", func(t *testing.T) {
		type v struct {
			Foo string `json:"foo"`
		}

		rq := v{Foo: "bar"}

		w := httptest.NewRecorder()

		if err := apiserver.WriteResponse(w, http.StatusOK, rq); err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(rq); err != nil {
			t.Fatal(err)
		}
		buf.WriteByte('\n')

		if !bytes.Equal(w.Body.Bytes(), buf.Bytes()) {
			t.Fatalf("expected body %q, got %q", buf.String(), w.Body.String())
		}
	})
}
