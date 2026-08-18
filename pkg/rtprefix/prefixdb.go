package rtprefix

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

var files embed.FS

const prefixCsv = "rt-prefixes-2023.csv"

type LongestPrefixLength struct {
	length int
	prefix string
}
type carrier int

const (
	us     = "United States of America"
	canada = "Canada"
)
const (
	mobile carrier = iota
	fixed
)

type rounding int

const (
	sec rounding = iota
	min
)

type ProviderInterface interface {
	GetName() string
	GetCountry() string
	GetType() int
}

type Provider struct {
	name      string
	country   string
	cType     carrier
	roundType rounding
}

type PrefixStore struct {
	prefixDict    map[string]Provider
	longestPrefix LongestPrefixLength
	countryMap    map[string]string
}

func (pr Provider) GetName() string {
	return pr.name
}
func (pr Provider) GetCountry() (cntry string) {
	cntry = pr.country
	if cntry == "" && pr.cType == 0 {
		cntryTmp := strings.Split(pr.GetName(), "-")
		cntry = cntryTmp[0]
		return
	}
	if cntry == "" && pr.name == us {
		cntry = us
		return
	}
	if cntry == "" && pr.name == canada {
		cntry = canada
		return
	}
	if cntry == "" && pr.cType == 1 {
		cntryTmp := strings.Split(pr.GetName(), "-")
		cntry = cntryTmp[0]
	}
	return cntry
}
func (pr Provider) GetType() int {
	return int(pr.cType)
}
func (ps *PrefixStore) GetDict() map[string]Provider {
	return ps.prefixDict
}
func (ps *PrefixStore) GetProvider(prfx string) Provider {
	return ps.prefixDict[prfx]
}
func (ps *PrefixStore) GetShortesPrefixByCountry(country string) (prfx string) {
	if country == us {
		prfx = "1"
		return
	}
	if country == canada {
		prfx = "1"
		return
	}
	prfx = ps.countryMap[country]
	return
}
func NewPrefixStore(rootFiles embed.FS) *PrefixStore {
	files = rootFiles
	ps := &PrefixStore{
		prefixDict:    make(map[string]Provider),
		countryMap:    make(map[string]string),
		longestPrefix: LongestPrefixLength{length: 0, prefix: ""},
	}
	ps.initAll()
	return ps
}

func (ps *PrefixStore) initAll() {
	f, err := os.Open(prefixCsv)
	if err != nil {
		slog.Error("Open embedding FS issue", "file", prefixCsv)
		return
	}
	defer f.Close()
	ps.fileRead(f)
}

func (ps *PrefixStore) CheckPrefix(query string) (string, error) {
	q := strings.TrimPrefix(query, "810")
	var longest string
	var prefix string
	queryLength := ps.longestPrefix.length
	qLen := len(q)
	if qLen <= queryLength {
		queryLength = qLen - 1
	}
	slog.Debug("query Length", "length", queryLength, "query length", qLen, "query string", q)
	for i := 0; i <= queryLength; i++ {
		prefix = prefix + string(q[i])
		_, ok := ps.prefixDict[prefix]
		if ok {
			longest = prefix
		}
	}
	if longest == "" {
		return "", errors.New("Prefix not found")
	}
	res := strings.TrimSpace(longest)

	return res, nil
}

func (ps *PrefixStore) fileRead(file *os.File) {

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		tmp := strings.Split(scanner.Text(), "|")
		if len(tmp) != 2 {
			continue
		}
		ps.parsePrefix(tmp)
	}

	if err := scanner.Err(); err != nil {
		slog.Error("File scan error", "Error:", err)
	}
}

func (ps *PrefixStore) parsePrefix(stra []string) {
	var prfx []string
	providerName := stra[0]
	tmp := strings.Split(stra[1], ";")
	for _, v := range tmp {
		p := strings.TrimSpace(v)
		rng := strings.Split(p, "-")
		if len(rng) != 2 {
			prfx = append(prfx, p)
			continue
		}
		startDigit, err := strconv.Atoi(rng[0])
		if err != nil {
			continue
		}
		endDigit, err := strconv.Atoi(rng[1])
		if err != nil {
			continue
		}
		for i := startDigit; i <= endDigit; i++ {
			prfx = append(prfx, strconv.Itoa(i))
		}
	}
	ps.addPrefixToMap(providerName, prfx)
}

func (ps *PrefixStore) addPrefixToMap(providerName string, prfx []string) {
	for _, v := range prfx {
		var provider Provider
		provider.name = providerName
		provider.roundType = sec
		provider.cType = getCarrierType(providerName)
		countryName := getCountry(providerName)
		err := ps.addCountryName(countryName, v)
		if err == nil {
			provider.country = countryName
		}
		ps.prefixDict[v] = provider
		ps.setLongestPrefixLength(v)
	}
}
func (ps *PrefixStore) addCountryName(cname, code string) error {
	_, ok := ps.countryMap[cname]
	if ok {
		return errors.New(fmt.Sprintf("Country %s entry already exists", cname))
	}
	ps.countryMap[cname] = code
	return nil
}
func getCountry(pname string) string {
	tmp := strings.Split(pname, "-")
	return tmp[0]
}
func getCarrierType(providerName string) carrier {
	if strings.Contains(providerName, "mobile") {
		return mobile
	}
	return fixed
}
func (ps *PrefixStore) setLongestPrefixLength(str string) {
	l := len(str)
	if l > ps.longestPrefix.length {
		ps.longestPrefix.length = l
		ps.longestPrefix.prefix = str
	}
}

func (ps *PrefixStore) GetCountry(cname string) (string, error) {
	err := errors.New("GetCountry error")
	cCode, ok := ps.countryMap[cname]
	if !ok {
		return "", err
	}
	return cCode, nil
}
