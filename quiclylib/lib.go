package quiclylib

import (
	"github.com/Project-Faster/quicly-go/internal/bindings"
	"github.com/Project-Faster/quicly-go/quiclylib/errors"
	"github.com/Project-Faster/quicly-go/quiclylib/types"
	"os"
	"path/filepath"
	"strings"
)

var _ types.ServerSession = &QServerSession{}
var _ types.ClientSession = &QClientSession{}
var _ types.Stream = &QStream{}

const (
	READ_SIZE         = 512 * 1024
	SMALL_BUFFER_SIZE = 4 * 1024
)

func QuiclyInitializeEngine(options types.Options) int {
	bindings.ResetRegistry()

	proto, cc_req, certfile, certkey, idle_timeout := options.Get()

	cwd, _ := os.Getwd()
	if !filepath.IsAbs(certfile) {
		certfile = filepath.Join(cwd, certfile)
	}
	if !filepath.IsAbs(certkey) {
		certkey = filepath.Join(cwd, certkey)
	}

	cc_algo := errors.QUICLY_CC_RENO
	switch strings.ToLower(cc_req) {
	case "cubic":
		cc_algo = errors.QUICLY_CC_CUBIC
		break
	case "pico":
		cc_algo = errors.QUICLY_CC_PICO
		break
	case "search":
		cc_algo = errors.QUICLY_CC_SEARCH
		break

	case "reno":
		fallthrough
	default:
		cc_algo = errors.QUICLY_CC_RENO
		break
	}

	result := bindings.QuiclyInitializeEngine(proto, certfile, certkey, idle_timeout, uint64(cc_algo))
	return int(result)
}

func QuiclyCloseEngine() int {
	result := bindings.QuiclyCloseEngine()
	return int(result)
}
