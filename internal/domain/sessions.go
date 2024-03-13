package domain

import "time"

type Session struct {
	UserId       int
	RefreshToken string
	ExpiresAt    time.Time
}

func SortSessionsByTime(sessions *[]Session) {
	quickSortSessions(sessions, 0, len(*sessions))
}

func quickSortSessions(sessions *[]Session, start int, end int) {
	if end-start <= 1 {
		return
	}
	pivot := (*sessions)[start]
	var less, greater []Session
	for i := start + 1; i < end; i++ {
		if (*sessions)[i].ExpiresAt.Unix() > pivot.ExpiresAt.Unix() {
			greater = append(greater, (*sessions)[i])
		} else {
			less = append(less, (*sessions)[i])
		}
	}
	for i := 0; i < len(less); i++ {
		(*sessions)[start+i] = less[i]
	}
	(*sessions)[start+len(less)] = pivot
	for i := 0; i < len(greater); i++ {
		(*sessions)[start+len(less)+1+i] = greater[i]
	}
	quickSortSessions(sessions, start, start+len(less))
	quickSortSessions(sessions, start+len(less)+1, end)
}
