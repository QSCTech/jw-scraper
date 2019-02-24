package jw_scraper

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type httpServiceImpl struct {
	baseUrl             *url.URL
	noRedirectionClient *http.Client
}

func NewHttpService(rawBaseUrl string) HttpService {
	baseUrl, err := url.Parse(rawBaseUrl)
	if err != nil {
		panic("invalid base URL: " + err.Error())
	}

	return &httpServiceImpl{
		baseUrl: baseUrl,
		noRedirectionClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (service *httpServiceImpl) NewJwCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:  "ASP.NET_SessionId",
		Value: value,
	}
}

func IsJwCookie(cookie *http.Cookie) bool {
	return cookie.Name == "ASP.NET_SessionId"
}

func (service *httpServiceImpl) GetLoginPage() (page string, err error) {
	resp, err := http.Get(service.baseUrl.String())
	if err != nil {
		return
	}
	defer func() {
		if closeBodyErr := resp.Body.Close(); closeBodyErr != nil {
			err = closeBodyErr
		}
	}()
	data, err := ioutil.ReadAll(transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())) // decoding from gbk is necessary
	if err != nil {
		return
	}
	page = string(data)
	return
}

/**
@Error: if err is kind of grpc.Status with code codes.InvalidArgument, username or password is invalid,
		otherwise, a http error occurred
*/
func (service *httpServiceImpl) Login(username, password, viewState string) (jwbCookie string, err error) {
	targetUrl := *service.baseUrl
	targetUrl.Path = path.Join(targetUrl.Path, "default2.aspx")
	req, err := http.NewRequest(
		http.MethodPost,
		targetUrl.String(),
		strings.NewReader(
			url.Values{
				"__EVENTTARGET":    {"Button1"},
				"__EVENTARGUMENT":  {""},
				"__VIEWSTATE":      {viewState},
				"TextBox1":         {username},
				"TextBox2":         {password},
				"RadioButtonList1": {"%D1%A7%C9%FA"}, // decoding in gbk: "学生"
				"Text1":            {""},
			}.Encode(),
		),
	)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := service.noRedirectionClient.Do(req)
	if err != nil {
		return
	}

	if err = resp.Body.Close(); err != nil {
		return
	}

	if resp.StatusCode == http.StatusFound {
		for _, cookie := range resp.Cookies() {
			if IsJwCookie(cookie) {
				jwbCookie = cookie.Value
				return
			}
		}
	}

	err = status.Error(codes.InvalidArgument, "authenticating fails")
	return
}

func (service *httpServiceImpl) GetDefaultCourses(stuId, jwbCookie string) (page string, status int, err error) {
	targetUrl := *service.baseUrl
	targetUrl.Path = path.Join(targetUrl.Path, "xskbcx.aspx")
	finalQuery := targetUrl.Query()
	finalQuery.Add("xh", stuId)
	targetUrl.RawQuery = finalQuery.Encode()
	req, err := http.NewRequest(http.MethodGet, targetUrl.String(), nil)
	if err != nil {
		return
	}
	req.AddCookie(service.NewJwCookie(jwbCookie))
	return service.sendRequest(req)
}

func (service *httpServiceImpl) GetCourses(stuId, jwbCookie, schoolYear, semester, viewState, eventTarget string) (page string, status int, err error) {
	xqdEncoded, _, err := transform.String(simplifiedchinese.GBK.NewEncoder(), semester)
	if err != nil {
		panic(err)
	}

	targetUrl := *service.baseUrl
	targetUrl.Path = path.Join(targetUrl.Path, "xskbcx.aspx")
	finalQuery := targetUrl.Query()
	finalQuery.Add("xh", stuId)
	targetUrl.RawQuery = finalQuery.Encode()

	req, err := http.NewRequest(
		http.MethodPost,
		targetUrl.String(),
		strings.NewReader(
			url.Values{
				"__EVENTTARGET":   {eventTarget},
				"__EVENTARGUMENT": {""},
				"__VIEWSTATE":     {viewState},
				"xxms":            {string([]byte{0xC1, 0xD0, 0xB1, 0xED})}, // decoding in gbk: "列表"
				"xnd":             {schoolYear},
				"xqd":             {xqdEncoded},
				"kcxx":            {""},
			}.Encode(),
		))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded") // It's necessary
	req.AddCookie(service.NewJwCookie(jwbCookie))
	return service.sendRequest(req)
}

func (service *httpServiceImpl) GetDefaultExams(stuId, jwbCookie string) (page string, status int, err error) {
	targetUrl := *service.baseUrl
	targetUrl.Path = path.Join(targetUrl.Path, "xskscx.aspx")
	finalQuery := targetUrl.Query()
	finalQuery.Add("xh", stuId)
	targetUrl.RawQuery = finalQuery.Encode()

	req, err := http.NewRequest(http.MethodGet, targetUrl.String(), nil)
	if err != nil {
		return
	}
	req.AddCookie(service.NewJwCookie(jwbCookie))
	return service.sendRequest(req)
}

func (service *httpServiceImpl) GetExams(stuId, jwbCookie, schoolYear, semester, viewState, eventTarget string) (page string, status int, err error) {
	targetUrl := *service.baseUrl
	targetUrl.Path = path.Join(targetUrl.Path, "xskscx.aspx")
	finalQuery := targetUrl.Query()
	finalQuery.Add("xh", stuId)
	targetUrl.RawQuery = finalQuery.Encode()

	xqdEncoded, _, err := transform.String(simplifiedchinese.GBK.NewEncoder(), semester)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		targetUrl.String(),
		strings.NewReader(
			url.Values{
				"__EVENTTARGET":   {eventTarget},
				"__EVENTARGUMENT": {""},
				"__VIEWSTATE":     {viewState},
				"xnd":             {schoolYear},
				"xqd":             {xqdEncoded},
			}.Encode(),
		))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded") // It's necessary
	req.AddCookie(service.NewJwCookie(jwbCookie))
	return service.sendRequest(req)
}

func (service *httpServiceImpl) GetScores() (page string, status int, err error) {
	panic("not implemented")
}

func (service *httpServiceImpl) GetMajorScores() (page string, status int, err error) {
	panic("not implemented")
}

func (service *httpServiceImpl) GetTotalCredit() (page string, status int, err error) {
	panic("not implemented")
}

func (service *httpServiceImpl) sendRequest(req *http.Request) (page string, status int, err error) {
	resp, err := service.noRedirectionClient.Do(req)
	if err != nil {
		return
	}

	defer func() {
		if closeBodyErr := resp.Body.Close(); err == nil && closeBodyErr != nil {
			err = closeBodyErr
			return
		}
	}()

	data, err := ioutil.ReadAll(transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return
	}
	page = string(data)
	status = resp.StatusCode
	return
}
