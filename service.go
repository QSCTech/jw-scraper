package jw_scraper

type HttpService interface {
	GetLoginPage() (page string, err error)

	Login(username, password, viewState string) (jwbCookie string, err error)

	GetDefaultCourses(stuId, jwbCookie string) (page string, status int, err error)

	GetCourses(stuId, jwbCookie, schoolYear, semester, viewState, eventTarget string) (page string, status int, err error)

	GetDefaultExams(stuId, jwbCookie string) (page string, status int, err error)

	GetExams(stuId, jwbCookie, schoolYear, semester, viewState, eventTarget string) (page string, status int, err error)

	GetScoresBase(stuId, jwbCookie string) (page string, status int, err error)

	GetScores(stuId, jwbCookie, schoolYear, viewState string) (page string, status int, err error)

	GetMajorScores(stuId, jwbCookie string) (page string, status int, err error)

	GetTotalCredit(stuId, jwbCookie string) (page string, status int, err error)
}
