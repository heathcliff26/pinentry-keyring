package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPinentry(t *testing.T) {
	// Prepare mock pinentry
	tmpDir := t.TempDir()
	pinentry := filepath.Join(tmpDir, "mock")
	out, err := exec.Command("go", "build", "-o", pinentry, "testdata/mock.go").CombinedOutput()
	require.NoError(t, err, "Should compile mock pinentry")
	t.Log(string(out))

	t.Setenv("PINENTRY", pinentry)

	t.Run("WithoutCaching", func(t *testing.T) {
		require := require.New(t)

		output, stdout, stdin := newMainTest(t)

		finished := make(chan struct{})
		go func() {
			main()
			close(finished)
		}()

		require.True(stdout.Scan())
		require.Equal("OK Pleased to meet you", stdout.Text(), "Should send greeting")

		_, err := stdin.Write([]byte("SETKEYINFO --clear\n"))
		require.NoError(err, "Should send SETKEYINFO")
		require.True(stdout.Scan())
		require.Equal("OK", stdout.Text(), "Should receive confirmation")

		_, err = stdin.Write([]byte("GETPIN\n"))
		require.NoError(err, "Should send GETPIN")
		require.True(stdout.Scan())
		require.Equal("S PASSWORD_FROM_CACHE", stdout.Text(), "Should receive password prompt")
		require.True(stdout.Scan())
		require.Equal("D 1234", stdout.Text(), "Should receive password")

		_, err = stdin.Write([]byte("BYE\n"))
		require.NoError(err, "Should send BYE")
		require.True(stdout.Scan())
		require.Equal("OK closing connection", stdout.Text(), "Should receive confirmation")

		select {
		case <-finished:
		case <-time.After(100 * time.Millisecond):
			t.Fatal("main() should have finished")
		}

		buf, err := os.ReadFile(output)
		require.NoError(err, "Should read output file")
		require.Equal(expectedOutputWithoutCache, string(buf), "Should match expected output")
	})
	t.Run("Cache", func(t *testing.T) {
		require := require.New(t)

		output, stdout, stdin := newMainTest(t)

		finished := make(chan struct{})
		go func() {
			main()
			close(finished)
		}()

		require.True(stdout.Scan())
		require.Equal("OK Pleased to meet you", stdout.Text(), "Should send greeting")

		_, err := stdin.Write([]byte("OPTION allow-external-password-cache\n"))
		require.NoError(err, "Should send OPTION")
		require.True(stdout.Scan())
		require.Equal("OK", stdout.Text(), "Should receive confirmation")

		_, err = stdin.Write([]byte("SETKEYINFO test-key\n"))
		require.NoError(err, "Should send SETKEYINFO")
		require.True(stdout.Scan())
		require.Equal("OK", stdout.Text(), "Should receive confirmation")

		_, err = stdin.Write([]byte("GETPIN\n"))
		require.NoError(err, "Should send GETPIN")
		require.True(stdout.Scan())
		require.Equal("S PASSWORD_FROM_CACHE", stdout.Text(), "Should receive password prompt")
		require.True(stdout.Scan())
		require.Equal("D 1234", stdout.Text(), "Should receive password")

		_, err = stdin.Write([]byte("BYE\n"))
		require.NoError(err, "Should send BYE")
		require.True(stdout.Scan())
		require.Equal("OK closing connection", stdout.Text(), "Should receive confirmation")

		select {
		case <-finished:
		case <-time.After(100 * time.Millisecond):
			t.Fatal("main() should have finished")
		}

		buf, err := os.ReadFile(output)
		require.NoError(err, "Should read output file")
		require.Equal(expectedOutputCache, string(buf), "Should match expected output")
	})
}

func TestCreatePinentryCMD(t *testing.T) {
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})
	testArgs := []string{"pinentry-keyring", "foo", "bar"}
	os.Args = testArgs

	t.Run("Default", func(t *testing.T) {
		assert := assert.New(t)

		t.Setenv("PINENTRY", "")
		cmd := createPinentryCMD()
		assert.Equal("/usr/bin/pinentry", cmd.Path, "Should have correct path")
		assert.Equal(testArgs[1:], cmd.Args[1:], "Should preserve args")
	})
	t.Run("WithEnv", func(t *testing.T) {
		assert := assert.New(t)

		t.Setenv("PINENTRY", "/foo/bar/testcommand")
		cmd := createPinentryCMD()
		assert.Equal("/foo/bar/testcommand", cmd.Path, "Should have correct path")
		assert.Equal(testArgs[1:], cmd.Args[1:], "Should preserve args")
	})

}

func newMainTest(t *testing.T) (string, *bufio.Scanner, io.Writer) {
	t.Helper()

	require := require.New(t)

	tmpDir := t.TempDir()
	output := filepath.Join(tmpDir, "output.txt")
	t.Setenv("PINENTRY_MOCK_OUTPUT", output)

	// Save os.Stdout/in and restore after test
	osStdout := os.Stdout
	osStdin := os.Stdin
	t.Cleanup(func() {
		os.Stdout = osStdout
		os.Stdin = osStdin
	})

	outR, outW, err := os.Pipe()
	require.NoError(err, "Should create pipe for Stdout")
	t.Cleanup(func() {
		outW.Close()
		outR.Close()
	})
	os.Stdout = outW
	out := bufio.NewScanner(outR)

	inR, inW, err := os.Pipe()
	require.NoError(err, "Should create pipe for Stdin")
	t.Cleanup(func() {
		inW.Close()
		inR.Close()
	})
	os.Stdin = inR

	return output, out, inW
}

var (
	expectedOutputWithoutCache = `OK Pleased to meet you
OPTION allow-external-password-cache
OK
SETKEYINFO pinenty-keyring-default-key
OK
GETPIN
S PASSWORD_FROM_CACHE
D 1234
BYE
OK closing connection
`
	expectedOutputCache = `OK Pleased to meet you
OPTION allow-external-password-cache
OK
OPTION allow-external-password-cache
OK
SETKEYINFO test-key
OK
GETPIN
S PASSWORD_FROM_CACHE
D 1234
BYE
OK closing connection
`
)
