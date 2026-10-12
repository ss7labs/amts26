package r11

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/ss7labs/amts26/pkg/rtprefix"

	"os"
	"slices"

	"github.com/jedib0t/go-pretty/table"
	"github.com/ss7labs/amts26/internal/app"
	"github.com/ss7labs/amts26/internal/model"
	"gorm.io/gorm"
)

const (
	//trunks names
	uz    = "BWTASH,OTASH,BTASH2"
	kz    = "BWAATA,OAATA,OASTAN,ASTANA,AATANG"
	ank   = "BWANK"
	ist   = "BWIST"
	uae   = "BWEMIR"
	mn2   = "BWMN2"
	mn102 = "MN10.2"
	//operators names
	ANK   = "ANK"
	IST   = "IST"
	UAE   = "UAE"
	MN2   = "MN2"
	MN102 = "MN102"
)

type report struct {
	appCtx          app.AppContext
	ps              *rtprefix.PrefixStore
	operators       map[string]operator
	opsNames        []string
	currentOperator operator
	period          string
}
type operator struct {
	trunks  string
	exclude string
	round   bool
}

func Report(date string, app app.AppContext) {
	rp := &report{appCtx: app}
	rp.period = date
	rp.ps = app.GetPrefixStore()
	rp.initAllOps()
	for _, opName := range rp.opsNames {
		rp.currentOperator = rp.operators[opName]
		if rp.opIsDead() {
			continue
		}
		fmt.Printf("Трафик исходящих(транзитных) переговоров на направление %s за %s.\n", opName, date)
		rp.walkOnProvider(date, rp.currentOperator)
	}
}
func (rp *report) opIsDead() (dead bool) {
	dead = false
	var count int64
	ops := createOps(rp.currentOperator.trunks)
	db := rp.appCtx.DB()
	clause := "b_number NOT LIKE '810" + rp.currentOperator.exclude + "%'"
	db.Model(&model.Event{}).
		Scopes(filterByOperator(ops)).
		Where(clause).
		Where("to_char(start_time, 'YYYYMM') = ?", rp.period).
		Count(&count)

	if count <= 0 {
		dead = true
	}

	return
}

func (rp *report) initAllOps() {
	rp.opsNames = []string{ANK, IST, UAE, MN2, MN102}
	rp.operators = make(map[string]operator)
	rp.operators[rp.opsNames[0]] = operator{ank, "90", false}
	rp.operators[rp.opsNames[1]] = operator{ist, "90", false}
	rp.operators[rp.opsNames[2]] = operator{uae, "971", false}
	rp.operators[rp.opsNames[3]] = operator{mn2, "7", true}
	rp.operators[rp.opsNames[4]] = operator{mn102, "7", true}
}

func filterByTime(date string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("to_char(start_time, 'YYYYMM') = ?", date)
	}
}
func filterByOperator(op []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("out_trunk IN (?)", op)
	}
}
func (rp *report) walkOnProvider(date string, provider operator) {
	ops := createOps(provider.trunks)
	var events []model.Event
	var clause string
	if provider.exclude != "7" {
		clause = "b_number NOT LIKE '810" + provider.exclude + "%'"
	}
	db := rp.appCtx.DB()
	if clause == "" {
		db.Scopes(filterByTime(date), filterByOperator(ops)).Find(&events)
	} else {
		db.Scopes(filterByTime(date), filterByOperator(ops)).Where(clause).Find(&events)
	}
	rp.walkOnEvents(&events)
}

type total struct {
	prefix         string
	count          int
	duration       int
	mobileCount    int
	mobileDuration int
	fixedCount     int
	fixedDuration  int
}

// TODO
// Define provider rounding schema
// Apply rounding schema depends on provider
var patternR = `(8107\d+[a-z]|8107[a-z])`
var reR = regexp.MustCompile(patternR)

func (rp *report) walkOnEvents(events *[]model.Event) {
	aggregate := make(map[string]*total)
	for _, event := range *events {
		code, err := rp.ps.CheckPrefix(event.BNumber)
		if err != nil {
			continue
		}
		m := reR.MatchString(event.BNumber)
		if m {
			continue
		}

		if rp.currentOperator.exclude == "7" && prfxSearch("8107", event.BNumber) {

			if excludeByOuterSymbolOnly("8107", event.BNumber) {
				continue
			}
		}

		pr := rp.ps.GetProvider(code)
		country := pr.GetCountry()

		durationRaw := int(event.Duration)
		durationRounded := durationRaw
		if rp.currentOperator.round {
			durationRounded = rp.roundDuration(country, event.BNumber, durationRaw)
		}

		if _, ok := aggregate[country]; !ok {
			aggregate[country] = &total{count: 1, duration: durationRounded, prefix: code}
			continue
		}
		aggregate[country].count++
		aggregate[country].duration += durationRounded

		isMobile := func(i int) bool {
			if i == 0 {
				return true
			}
			return false
		}

		if isMobile(pr.GetType()) {
			aggregate[country].mobileCount++
			aggregate[country].mobileDuration += durationRounded
			continue
		}
		aggregate[country].fixedCount++
		aggregate[country].fixedDuration += durationRounded
	}
	rp.drawTable(aggregate)
}
func excludeByOuterSymbolOnly(patt, num string) bool {
	//8107
	exclude := true
	//81077
	lastChar := patt[len(patt)-1:]
	if prfxSearch(patt+lastChar, num) {
		exclude = false
	}
	return exclude
}

const (
	AF = "Afghanistan"
	CN = "China"
	AZ = "Azerbaijan"
	KR = "Korea (Rep.)"
	TH = "Thailand"
	MY = "Malaysia"
	VN = "Viet Nam"
	IN = "India"
)

var round_0_60_60 map[string]any = map[string]any{
	IN: nil,
}
var roundCountries map[string]any = map[string]any{
	CN: nil,
	AZ: nil,
	KR: nil,
	TH: nil,
	MY: nil,
	VN: nil,
	IN: nil,
}

func (rp *report) roundDuration(country, bnum string, durationRaw int) (seconds int) {
	var roundByPrefix = func(patt string) {
		if prfxSearch(patt, bnum) && durationRaw < 60 {
			seconds = 60
			return
		}
	}
	roundByPrefix("8109377")
	roundByPrefix("81025")
	roundByPrefix("81022")

	_, ok := roundCountries[country]
	if ok && durationRaw < 60 {
		seconds = 60
		return
	}
	// India rounding as 0-60-60
	// if duration < 60 then above rounding
	_, ok = round_0_60_60[country]
	if ok && durationRaw > 60 {
		seconds = int(math.Ceil(float64(durationRaw)/60) * 60)
		return
	}

	seconds = durationRaw
	return
}

func prfxSearch(prfx, bnum string) bool {
	pattern := `^` + prfx + `\d+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(bnum)
}

func createOps(provider string) []string {
	var ops []string
	tmp := strings.Split(provider, ",")
	for _, op := range tmp {
		ops = append(ops, op)
	}
	return ops
}

func (rp *report) drawTable(mfs map[string]*total) {
	tLen := len(mfs)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"КОД", "СТРАНА", "СОТ_КОЛ", "СОТ_МИН", "ФИК_КОЛ", "ФИК_МИН", "ИТОГО_КОЛ", "ИТОГО_МИН"})

	var cntT, minT int
	orderedList := mapToList(mfs)
	for _, country := range orderedList {
		mf := mfs[country]
		if mf.duration < 60 {
			continue
		}
		cntT += mf.count
		minutes := int(mf.duration / 60)
		minT += minutes
		prfx := rp.ps.GetShortesPrefixByCountry(country)
		if prfx == "" {
			prfx = mf.prefix
		}
		t.AppendRow(table.Row{prfx, country, 0, 0, 0, 0, mf.count, minutes})
	}

	if tLen > 1 {
		t.AppendFooter(table.Row{"Итого", "", 0, 0, 0, 0, cntT, minT})
	}
	t.Render()
}

func mapToList(mfs map[string]*total) []string {
	var list []string
	for key, _ := range mfs {
		list = append(list, key)
	}
	slices.Sort(list)
	return list
}
