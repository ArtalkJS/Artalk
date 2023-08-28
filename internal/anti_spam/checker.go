package anti_spam

type CheckerName string

type CheckerParams struct {
	BlogURL string

	Content   string
	CommentID uint

	UserName  string
	UserEmail string
	UserID    uint
	UserIP    string
	UserAgent string
}

type Checker interface {
	Name() string
	Check(p *CheckerParams) (bool, error)
}
