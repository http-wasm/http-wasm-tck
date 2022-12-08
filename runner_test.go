package http_wasm_tck

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/http-wasm/http-wasm-host-go/tck"
)

// TestMain is the entrypoint for the test binary. Commands like backend and extract-guest
// don't actually run tests, but we can still have a single binary for all commands by
// implementing them as part of TestMain.
func TestMain(m *testing.M) {
	var help bool
	flag.BoolVar(&help, "h", false, "print usage")

	flag.Parse()

	if help || flag.NArg() == 0 {
		printUsage(os.Stderr)
		os.Exit(0)
	}

	subCmd := flag.Arg(0)
	switch subCmd {
	case "backend":
		backend(flag.Args()[1:])
	case "extract-guest":
		extractGuest(flag.Args()[1:])
	case "run":
		// Make sure test runs are outputted.
		_ = flag.Lookup("test.v").Value.Set("true")

		flags := flag.NewFlagSet("run", flag.ExitOnError)
		flags.SetOutput(os.Stderr)

		var help bool
		flags.BoolVar(&help, "h", false, "print usage")

		var addr string
		flags.StringVar(&addr, "url", "http://localhost:8080", "URL to send test requests too")

		_ = flags.Parse(flag.Args()[1:])

		if help {
			printRunUsage(os.Stderr, flags)
			os.Exit(0)
		}

		os.Exit(m.Run())
	}
}

func backend(args []string) {
	flags := flag.NewFlagSet("backend", flag.ExitOnError)

	var help bool
	flags.BoolVar(&help, "h", false, "print usage")

	var addr string
	flags.StringVar(&addr, "addr", "0.0.0.0:9080", "address to listen on")

	_ = flags.Parse(args)

	if help {
		printBackendUsage(os.Stderr, flags)
		os.Exit(0)
	}

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)

	s := tck.StartBackend(addr)
	// Don't change this string since we allow it to be parsed to find the address.
	fmt.Printf("Started backend server on %s\n", s.Listener.Addr())

	for {
		select {
		case <-exitCh:
			return
		}
	}
}

func extractGuest(args []string) {
	flags := flag.NewFlagSet("extract-guest", flag.ExitOnError)
	flags.SetOutput(os.Stderr)

	var help bool
	flags.BoolVar(&help, "h", false, "print usage")

	_ = flags.Parse(args)

	if help {
		printExtractGuestUsage(os.Stderr)
		os.Exit(0)
	}

	if flags.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing path to extract wasm file to")
		printExtractGuestUsage(os.Stderr)
		os.Exit(1)
	}

	wasmPath := flags.Arg(0)

	if err := os.WriteFile(wasmPath, tck.GuestWASM, 0o644); err != nil {
		fmt.Printf("error opening path: %v\n", err)
		os.Exit(1)
	}
}

func printUsage(stdErr io.Writer) {
	fmt.Fprintln(stdErr, "http-wasm TCK")
	fmt.Fprintln(stdErr)
	fmt.Fprintln(stdErr, "Usage:\n  http-wasm-tck <command>")
	fmt.Fprintln(stdErr)
	fmt.Fprintln(stdErr, "Commands:")
	fmt.Fprintln(stdErr, "  extract\t\tExtracts the guest wasm file")
	fmt.Fprintln(stdErr, "  backend\t\tStarts the backend server")
	fmt.Fprintln(stdErr, "  run\t\t\tRuns the tests")
}

func printBackendUsage(stdErr io.Writer, flags *flag.FlagSet) {
	fmt.Fprintln(stdErr, "http-wasm TCK - backend server")
	fmt.Fprintln(stdErr)
	fmt.Fprintln(stdErr, "Usage:\n  http-wasm-tck backend <options>")
	fmt.Fprintln(stdErr)
	fmt.Fprintln(stdErr, "Options:")
	flags.PrintDefaults()
}

func printExtractGuestUsage(stdErr io.Writer) {
	fmt.Fprintln(stdErr, "http-wasm TCK - extract guest")
	fmt.Fprintln(stdErr)
	fmt.Fprintln(stdErr, "Usage:\n  http-wasm-tck extract-guest <path to wasm file>")
	fmt.Fprintln(stdErr)
}

func printRunUsage(stdErr io.Writer, flags *flag.FlagSet) {
	fmt.Fprintln(stdErr, "http-wasm TCK - run the tests")
	fmt.Fprintln(stdErr)
	fmt.Fprintln(stdErr, "Usage:\n  http-wasm-tck run <options>")
	fmt.Fprintln(stdErr)
	fmt.Fprintln(stdErr, "Options:")
	flags.PrintDefaults()
}

// TestTCK is a standard unit test function. It is not invoked by us directly, instead being
// run by Go's test framework reflectively when testing.M.Run is invoked in TestMain.
func TestTCK(t *testing.T) {
	flags := flag.NewFlagSet("run", flag.ExitOnError)
	flags.SetOutput(os.Stderr)

	var addr string
	flags.StringVar(&addr, "url", "http://localhost:8080", "URL to send test requests too")

	_ = flags.Parse(flag.Args()[1:])

	tck.Run(t, addr)
}
