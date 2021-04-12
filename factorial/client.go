package factorial

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"golang.org/x/net/publicsuffix"
)

const BASE_URL = "https://api.factorialhr.com"

type factorialClient struct {
	http.Client
	employee_id int
	period_id   int
	calendar    []calendarDay
	shifts      []shift
	year        int
	month       int
	clock_in    string
	clock_out   string
}

type period struct {
	Id          int
	Employee_id int
	Year        int
	Month       int
}

type calendarDay struct {
	Id           string
	Day          int
	Date         string
	Is_laborable bool
	Is_leave     bool
	Leave_name   string
}

type shift struct {
	Id        int64
	Period_id int64
	Day       int
	Clock_in  string
	Clock_out string
	Minutes   int64
}

func NewFactorialClient(email, password string, year, month int, in, out string) *factorialClient {
	s := spinner.New(spinner.CharSets[14], 60*time.Millisecond)
	s.Suffix = " Logging in..."
	s.Start()
	c := new(factorialClient)
	c.year = year
	c.month = month
	c.clock_in = in
	c.clock_out = out
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&options)
	c.Client = http.Client{Jar: jar}
	c.login(email, password)
	c.setPeriodId()
	c.setCalendar()
	c.setShifts()
	s.Stop()
	return c
}

func (c *factorialClient) ClockIn(dry_run bool) {
	spinner := spinner.New(spinner.CharSets[14], 60*time.Millisecond)
	var t time.Time
	var message string
	for _, d := range c.calendar {
		spinner.Restart()
		t = time.Date(c.year, time.Month(c.month), d.Day, 0, 0, 0, 0, time.UTC)
		message = fmt.Sprintf("%s... ", t.Format("02 Jan"))
		spinner.Prefix = message + " "
		if c.clockedIn(d.Day) {
			message = fmt.Sprintf("%s ❌ Already clocked in\n", message)
		} else if d.Is_leave {
			message = fmt.Sprintf("%s ❌ %s\n", message, d.Leave_name)
		} else if !d.Is_laborable {
			message = fmt.Sprintf("%s ❌ %s\n", message, t.Format("Monday"))
		} else {
			if !dry_run {
				// clock in here!
			}
			time.Sleep(1e9)
			message = fmt.Sprintf("%s ✅ %s - %s\n", message, c.clock_in, c.clock_out)
		}
		spinner.Stop()
		fmt.Print(message)
	}
	fmt.Println("done!")
}

func (c *factorialClient) login(email, password string) {
	getCSRFToken := func(resp *http.Response) string {
		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		start := strings.Index(string(data), "<meta name=\"csrf-token\" content=\"") + 33
		end := strings.Index(string(data)[start:], "\" />")
		return string(data)[start : start+end]
	}

	getLoginError := func(resp *http.Response) string {
		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		start := strings.Index(string(data), "<div class=\"flash flash--wrong\">") + 32
		if start < 0 {
			return ""
		}
		end := strings.Index(string(data)[start:], "</div>")
		if start < 0 || end-start > 100 {
			return ""
		}
		return string(data)[start : start+end]
	}

	resp, _ := c.Get(BASE_URL + "/users/sign_in")
	csrf_token := getCSRFToken(resp)
	body := url.Values{
		"authenticity_token": {csrf_token},
		"return_host":        {"factorialhr.es"},
		"user[email]":        {email},
		"user[password]":     {password},
		"user[remember_me]":  {"0"},
		"commit":             {"Sign in"},
	}
	resp, _ = c.PostForm(BASE_URL+"/users/sign_in", body)
	if err := getLoginError(resp); err != "" {
		log.Fatal(err)
	}
}

func (c *factorialClient) setPeriodId() {
	resp, _ := c.Get(BASE_URL + "/attendance/periods")
	defer resp.Body.Close()
	var periods []period
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &periods)
	for _, p := range periods {
		if p.Year == c.year && p.Month == c.month {
			c.employee_id = p.Employee_id
			c.period_id = p.Id
			return
		}
	}
	log.Fatalf("Could not find the specified year/month in the available periods (%d/%d)\n", c.month, c.year)
}

func (c *factorialClient) setCalendar() {
	u, _ := url.Parse(BASE_URL + "/attendance/calendar")
	q := u.Query()
	q.Set("id", strconv.Itoa(c.employee_id))
	q.Set("year", strconv.Itoa(c.year))
	q.Set("month", strconv.Itoa(c.month))
	u.RawQuery = q.Encode()
	resp, _ := c.Get(u.String())
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &c.calendar)
	sort.Slice(c.calendar, func(i, j int) bool {
		return c.calendar[i].Day < c.calendar[j].Day
	})
}

func (c *factorialClient) setShifts() {
	u, _ := url.Parse(BASE_URL + "/attendance/shifts")
	q := u.Query()
	q.Set("employee_id", strconv.Itoa(c.employee_id))
	q.Set("year", strconv.Itoa(c.year))
	q.Set("month", strconv.Itoa(c.month))
	u.RawQuery = q.Encode()
	resp, _ := c.Get(u.String())
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &c.shifts)
}

func (c *factorialClient) clockedIn(day int) bool {
	for _, shift := range c.shifts {
		if shift.Day == day {
			return true
		}
	}
	return false
}
