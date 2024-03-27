package handlefuncs

import "sort"

type By func(p1, p2 *Comment) bool

type CommentSorter struct {
	all []Comment
	by  func(p1, p2 *Comment) bool
}

func (by By) Sort(all []Comment) {
	comsort := &CommentSorter{
		all: all,
		by:  by,
	}
	sort.Sort(comsort)
}

func (s *CommentSorter) Len() int {
	return len(s.all)
}

func (s *CommentSorter) Swap(i, j int) {
	s.all[i], s.all[j] = s.all[j], s.all[i]
}

func (s *CommentSorter) Less(i, j int) bool {
	return s.by(&s.all[i], &s.all[j])
}

func SortCommentsEarliest(all []Comment) []Comment {
	By(func(p1, p2 *Comment) bool {
		return p1.CreatedAt.Before(p2.CreatedAt)
	}).Sort(all)
	return all
}
