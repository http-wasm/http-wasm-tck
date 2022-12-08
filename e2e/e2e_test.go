package e2e

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	wasm "github.com/http-wasm/http-wasm-host-go/handler/nethttp"
)

func TestE2E(t *testing.T) {
	dir := t.TempDir()
	tckPath := filepath.Join(dir, "http-wasm-tck")
	guestPath := filepath.Join(dir, "guest.wasm")

	buildCmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-o", tckPath, "-c", "..")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("error building tck binary: %v", err)
	}

	extractCmd := exec.Command(tckPath, "extract-guest", guestPath)
	if err := extractCmd.Run(); err != nil {
		t.Fatalf("error extracting guest: %v", err)
	}

	guestWasm, err := os.ReadFile(guestPath)
	if err != nil {
		t.Fatalf("error reading guest wasm file: %v", err)
	}

	backendCmd := exec.Command(tckPath, "backend", "-addr", "0.0.0.0:0")
	backendOut := &strings.Builder{}
	backendCmd.Stdout = backendOut
	if err := backendCmd.Start(); err != nil {
		t.Fatalf("error starting backend: %v", err)
	}
	defer backendCmd.Process.Signal(syscall.SIGTERM)
	backendAddr := ""
	for i := 0; i < 20; i++ {
		s := backendOut.String()
		if strings.Contains(s, "Started backend server on") {
			_, addr, _ := strings.Cut(s, "Started backend server on ")
			backendAddr = strings.TrimSpace(addr)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if backendAddr == "" {
		t.Fatal("backend did not start after 20 attempts")
	}

	_, backendPort, _ := net.SplitHostPort(backendAddr)
	backendHost := fmt.Sprintf("localhost:%s", backendPort)

	for _, succeed := range []bool{true, false} {
		tn := "succeed"
		if !succeed {
			tn = "fail"
		}
		t.Run(tn, func(t *testing.T) {
			mw, err := wasm.NewMiddleware(context.Background(), guestWasm)
			if err != nil {
				t.Fatal(err)
			}
			// Proxying requests is not trivial
			h := mw.NewHandler(context.Background(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.URL.Scheme = "http"
				r.URL.Host = backendHost
				r.RequestURI = ""
				resp, err := http.DefaultClient.Do(r)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(err.Error()))
					return
				}
				w.WriteHeader(resp.StatusCode)

				// Simplest way to make the TCK to fail is to not forward headers.
				if succeed {
					for k, vs := range resp.Header {
						for _, v := range vs {
							w.Header().Add(k, v)
						}
					}
				}
				if resp.Body != nil {
					io.Copy(w, resp.Body)
				}
			}))

			server := httptest.NewServer(h)
			defer server.Close()

			runCmd := exec.Command(tckPath, "run", "-url", server.URL)
			err = runCmd.Run()
			if err != nil && succeed {
				t.Fatalf("error running tck, expected success: %v", err)
			} else if err == nil && !succeed {
				t.Fatal("no error running tck, expected failure")
			}
		})
	}
}
