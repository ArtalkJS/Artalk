package entity

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model

	Content string

	PageKey  string `gorm:"index;size:255"`
	SiteName string `gorm:"index;size:255"`

	// SHA-256 of the decoded browser secret. Empty means no guest proof; never
	// infer ownership for legacy/imported comments from their name or email.
	// Internal persistence only: not part of the public comment/transfer format.
	GuestTokenHash string `gorm:"index;size:64;not null;default:''" json:"-"`

	UserID     uint `gorm:"index"`
	IsVerified bool `gorm:"default:false"`
	UA         string
	IP         string

	Rid uint `gorm:"index"` // Parent Node ID

	IsCollapsed bool `gorm:"default:false"`
	IsPending   bool `gorm:"default:false"`
	IsPinned    bool `gorm:"default:false"`

	VoteUp   int
	VoteDown int

	RootID uint `gorm:"index"` // Root Node ID (can be derived from `Rid`)

	// Annotation anchor (W3C TextQuoteSelector). Empty means a normal page comment.
	AnchorExact    string `gorm:"type:text;not null;default:''"`
	AnchorPrefix   string `gorm:"type:text;not null;default:''"`
	AnchorSuffix   string `gorm:"type:text;not null;default:''"`
	AnchorStart    int    `gorm:"not null;default:-1"`
	AnchorEnd      int    `gorm:"not null;default:-1"`
	AnchorSection  string `gorm:"size:255;not null;default:''"`
	AnchorResolved bool   `gorm:"not null;default:false"`

	// IsAnnotationThread is true for annotation roots and all their replies.
	// It lets the normal comment list exclude annotation threads with one
	// indexed condition, keeping counts and pagination consistent.
	IsAnnotationThread bool `gorm:"index;not null;default:false"`

	// Associated Page
	//
	// Use Composite Foreign Keys for multiple-site support.
	Page *Page `gorm:"foreignKey:page_key,site_name;references:key,site_name"`

	// Associated User
	User *User `gorm:"foreignKey:user_id;references:id"`
}

func (c Comment) IsEmpty() bool {
	return c.ID == 0
}

func (c Comment) IsAllowReply() bool {
	return !c.IsCollapsed && !c.IsPending
}

func (c Comment) HasAnchor() bool {
	return c.AnchorExact != ""
}

// AnchorDTO builds the API anchor representation from the stored columns.
func (c Comment) AnchorDTO() *AnchorDTO {
	if c.AnchorExact == "" {
		return nil
	}

	a := &AnchorDTO{
		Quote: QuoteSelectorDTO{
			Type:   "TextQuoteSelector",
			Exact:  c.AnchorExact,
			Prefix: c.AnchorPrefix,
			Suffix: c.AnchorSuffix,
		},
		Section: c.AnchorSection,
	}
	if c.AnchorStart >= 0 && c.AnchorEnd >= c.AnchorStart {
		a.Position = &TextPositionDTO{Start: c.AnchorStart, End: c.AnchorEnd}
	}
	return a
}
