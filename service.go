package jw_scraper

type HttpService interface {
	GetLoginPage() (page string, err error)

	Login(username, password, viewState string) (jwbCookie string, err error)

	GetDefaultCourses(stuId, jwbCookie string) (page string, status int, err error)

	GetCourses(stuId, jwbCookie, schoolYear, semester, viewState, eventTarget string) (page string, status int, err error)

	GetDefaultExams(stuId, jwbCookie string) (page string, status int, err error)

	GetExams(stuId, jwbCookie, schoolYear, semester, viewState, eventTarget string) (page string, status int, err error)

	GetScores() (page string, status int, err error)

	GetMajorScores() (page string, status int, err error)

	GetTotalCredit() (page string, status int, err error)
}
