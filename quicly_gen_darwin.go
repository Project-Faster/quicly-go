//go:generate bash deps_build.sh
//go:generate bash regen.sh
//go:generate bash -c "cd tester/ && rm -f ./tester && go build -x -v -o tester"

package quicly
