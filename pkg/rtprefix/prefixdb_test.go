package rtprefix

import (
	"bufio"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var logLevel = slog.LevelError

func TestOneQuery(t *testing.T) {
	ps := NewPrefixStore(fileT)
	want := "431"
	result, err := ps.CheckPrefix("81043142406297590")
	require.NoError(t, err)
	assert.Equal(t, result, want)

	want = "79"
	result, err = ps.CheckPrefix("81079150835986")
	require.NoError(t, err)
	assert.Equal(t, result, want)
}
func TestOneAndProvider(t *testing.T) {
	ps := NewPrefixStore(fileT)
	want := "79"
	result, err := ps.CheckPrefix("81079150835986")
	require.NoError(t, err)
	assert.Equal(t, result, want)
	var pr ProviderInterface
	pr = ps.GetProvider(result)
	assert.Equal(t, "Russia-mobile", pr.GetName())
	fmt.Println("Provider", pr)
}

//go:embed rt-prefixes-2023.csv
var fileT embed.FS

//fs:embed mn102.txt
//fs:embed mn10210.txt
var filesystem embed.FS

const fn = "mn102.txt"

func setLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	// 2. Create the handler (slog.NewTextHandler or slog.NewJSONHandler)
	handler := slog.NewTextHandler(os.Stdout, opts)
	// 3. Set it as the application-wide default logger
	logger := slog.New(handler)
	return logger
}
func Test_4750_queries(t *testing.T) {

	slog.SetDefault(setLogger())

	f, err := os.Open(fn)
	if err != nil {
		t.Errorf("Error open file %s", fn)
	}
	defer f.Close()
	ps := NewPrefixStore(fileT)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		_, err := ps.CheckPrefix(line)
		require.NoError(t, err)
	}
	if err := scanner.Err(); err != nil {
		t.Error("File scan error", err)
	}

}

func TestGetCountry(t *testing.T) {
	ps := NewPrefixStore(fileT)
	cCode, err := ps.GetCountry(us)
	require.NoError(t, err)
	assert.Equal(t, "1201", cCode)

	cCode, err = ps.GetCountry(canada)
	require.NoError(t, err)
	assert.Equal(t, "1204", cCode)
	//shoukd be 216 lines as countries
	assert.Equal(t, 216, len(ps.countryMap))
}

func TestGetShortestPrefix(t *testing.T) {
	want := "1"
	ps := NewPrefixStore(fileT)
	code := ps.GetShortesPrefixByCountry(us)
	//fmt.Println("USA test", us, code)
	assert.Equal(t, want, code)
}
func TestGetShortestPrefixRussia(t *testing.T) {
	want := "7"
	ps := NewPrefixStore(fileT)
	code := ps.GetShortesPrefixByCountry("Russia")
	//fmt.Println("USA test", us, code)
	assert.Equal(t, want, code)
}
func TestGetShortestPrefixBelarus(t *testing.T) {
	want := "375"
	ps := NewPrefixStore(fileT)
	code := ps.GetShortesPrefixByCountry("Belarus")
	//fmt.Println("USA test", us, code)
	assert.Equal(t, want, code)
}
