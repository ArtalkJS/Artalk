package comments_get

type SortRule string

const (
	SortByDateDesc SortRule = "date_desc"
	SortByDateAsc  SortRule = "date_asc"
	SortByVote     SortRule = "vote"
)

// Get sort rule
func GetSortSQL(scope Scope, sortBy SortRule) string {
	switch sortBy {
	case SortByDateDesc:
		return "created_at DESC"
	case SortByDateAsc:
		return "created_at ASC"
	case SortByVote:
		return "vote_up DESC, created_at DESC"
	}

	if scope == ScopePage {
		return "is_pinned DESC, created_at DESC"
	}

	return "created_at DESC"
}
