func GetStudySessionWords(db *sql.DB, sessionID, page, limit int) ([]struct {
	ID             int    `json:"id"`
	English        string `json:"english"`
	Portuguese     string `json:"portuguese"`
	CorrectCount   int    `json:"correct_count"`
	IncorrectCount int    `json:"incorrect_count"`
}, int, error) {
	offset := (page - 1) * limit

	// First get total count
	var totalItems int
	err := db.QueryRow(`
        SELECT COUNT(DISTINCT w.id)
        FROM words w
        JOIN word_review_items wri ON w.id = wri.word_id
        WHERE wri.study_session_id = ?`, sessionID).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Then get paginated results with aggregated stats
	query := `
        SELECT 
            w.id,
            w.english,
            w.portuguese,
            SUM(CASE WHEN wri.is_correct = 1 THEN 1 ELSE 0 END) as correct_count,
            SUM(CASE WHEN wri.is_correct = 0 THEN 1 ELSE 0 END) as incorrect_count
        FROM words w
        JOIN word_review_items wri ON w.id = wri.word_id
        WHERE wri.study_session_id = ?
        GROUP BY w.id, w.english, w.portuguese
        ORDER BY w.id
        LIMIT ? OFFSET ?`

	rows, err := db.Query(query, sessionID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var words []struct {
		ID             int    `json:"id"`
		English        string `json:"english"`
		Portuguese     string `json:"portuguese"`
		CorrectCount   int    `json:"correct_count"`
		IncorrectCount int    `json:"incorrect_count"`
	}

	for rows.Next() {
		var word struct {
			ID             int    `json:"id"`
			English        string `json:"english"`
			Portuguese     string `json:"portuguese"`
			CorrectCount   int    `json:"correct_count"`
			IncorrectCount int    `json:"incorrect_count"`
		}
		if err := rows.Scan(&word.ID, &word.English, &word.Portuguese, &word.CorrectCount, &word.IncorrectCount); err != nil {
			return nil, 0, err
		}
		words = append(words, word)
	}

	return words, totalItems, rows.Err()
}

func GetStudySessionWordsRaw(db *sql.DB, sessionID, page, limit int) ([]struct {
	ID         int       `json:"id"`
	WordID     int       `json:"word_id"`
	English    string    `json:"english"`
	Portuguese string    `json:"portuguese"`
	IsCorrect  bool      `json:"is_correct"`
	CreatedAt  time.Time `json:"created_at"`
}, int, error) {
	offset := (page - 1) * limit

	// First get total count
	var totalItems int
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM word_review_items wri
        WHERE wri.study_session_id = ?`, sessionID).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Then get paginated results with word details
	query := `
        SELECT 
            wri.id,
            w.id as word_id,
            w.english,
            w.portuguese,
            wri.is_correct,
            wri.created_at
        FROM word_review_items wri
        JOIN words w ON w.id = wri.word_id
        WHERE wri.study_session_id = ?
        ORDER BY wri.created_at DESC
        LIMIT ? OFFSET ?`

	rows, err := db.Query(query, sessionID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []struct {
		ID         int       `json:"id"`
		WordID     int       `json:"word_id"`
		English    string    `json:"english"`
		Portuguese string    `json:"portuguese"`
		IsCorrect  bool      `json:"is_correct"`
		CreatedAt  time.Time `json:"created_at"`
	}

	for rows.Next() {
		var review struct {
			ID         int       `json:"id"`
			WordID     int       `json:"word_id"`
			English    string    `json:"english"`
			Portuguese string    `json:"portuguese"`
			IsCorrect  bool      `json:"is_correct"`
			CreatedAt  time.Time `json:"created_at"`
		}
		if err := rows.Scan(
			&review.ID,
			&review.WordID,
			&review.English,
			&review.Portuguese,
			&review.IsCorrect,
			&review.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		reviews = append(reviews, review)
	}

	return reviews, totalItems, rows.Err()
}

func GetWordGroupStudySessions(db *sql.DB, groupID, page, limit int) ([]struct {
	ID              int       `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	ActivityName    string    `json:"activity_name"`
	TotalWords      int       `json:"total_words"`
	CorrectCount    int       `json:"correct_count"`
	IncorrectCount  int       `json:"incorrect_count"`
	SuccessRate     float64   `json:"success_rate"`
	DurationMinutes int       `json:"duration_minutes"`
}, int, error) {
	offset := (page - 1) * limit

	// First get total count
	var totalItems int
	err := db.QueryRow(`
        SELECT COUNT(DISTINCT ss.id)
        FROM study_sessions ss
        WHERE ss.group_id = ?`, groupID).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Then get paginated results with statistics
	query := `
        SELECT 
            ss.id,
            ss.created_at,
            sa.name as activity_name,
            COUNT(DISTINCT wri.word_id) as total_words,
            SUM(CASE WHEN wri.is_correct = 1 THEN 1 ELSE 0 END) as correct_count,
            SUM(CASE WHEN wri.is_correct = 0 THEN 1 ELSE 0 END) as incorrect_count,
            ROUND(AVG(CASE WHEN wri.is_correct = 1 THEN 100.0 ELSE 0.0 END), 1) as success_rate,
            ROUND((julianday(MAX(wri.created_at)) - julianday(MIN(wri.created_at))) * 24 * 60) as duration_minutes
        FROM study_sessions ss
        JOIN study_activities sa ON ss.study_activity_id = sa.id
        LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
        WHERE ss.group_id = ?
        GROUP BY ss.id, ss.created_at, sa.name
        ORDER BY ss.created_at DESC
        LIMIT ? OFFSET ?`

	rows, err := db.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []struct {
		ID              int       `json:"id"`
		CreatedAt       time.Time `json:"created_at"`
		ActivityName    string    `json:"activity_name"`
		TotalWords      int       `json:"total_words"`
		CorrectCount    int       `json:"correct_count"`
		IncorrectCount  int       `json:"incorrect_count"`
		SuccessRate     float64   `json:"success_rate"`
		DurationMinutes int       `json:"duration_minutes"`
	}

	for rows.Next() {
		var session struct {
			ID              int       `json:"id"`
			CreatedAt       time.Time `json:"created_at"`
			ActivityName    string    `json:"activity_name"`
			TotalWords      int       `json:"total_words"`
			CorrectCount    int       `json:"correct_count"`
			IncorrectCount  int       `json:"incorrect_count"`
			SuccessRate     float64   `json:"success_rate"`
			DurationMinutes int       `json:"duration_minutes"`
		}
		if err := rows.Scan(
			&session.ID,
			&session.CreatedAt,
			&session.ActivityName,
			&session.TotalWords,
			&session.CorrectCount,
			&session.IncorrectCount,
			&session.SuccessRate,
			&session.DurationMinutes,
		); err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, session)
	}

	return sessions, totalItems, rows.Err()
}

func GetWordGroupStudySessionsRaw(db *sql.DB, groupID, page, limit int) ([]struct {
	ID              int       `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	ActivityName    string    `json:"activity_name"`
	WordID          int       `json:"word_id"`
	English         string    `json:"english"`
	Portuguese      string    `json:"portuguese"`
	IsCorrect       bool      `json:"is_correct"`
	ReviewCreatedAt time.Time `json:"review_created_at"`
}, int, error) {
	offset := (page - 1) * limit

	// First get total count
	var totalItems int
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM study_sessions ss
        JOIN word_review_items wri ON ss.id = wri.study_session_id
        WHERE ss.group_id = ?`, groupID).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Then get paginated results with all details
	query := `
        SELECT 
            ss.id,
            ss.created_at,
            sa.name as activity_name,
            w.id as word_id,
            w.english,
            w.portuguese,
            wri.is_correct,
            wri.created_at as review_created_at
        FROM study_sessions ss
        JOIN study_activities sa ON ss.study_activity_id = sa.id
        JOIN word_review_items wri ON ss.id = wri.study_session_id
        JOIN words w ON wri.word_id = w.id
        WHERE ss.group_id = ?
        ORDER BY ss.created_at DESC, wri.created_at ASC
        LIMIT ? OFFSET ?`

	rows, err := db.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []struct {
		ID              int       `json:"id"`
		CreatedAt       time.Time `json:"created_at"`
		ActivityName    string    `json:"activity_name"`
		WordID          int       `json:"word_id"`
		English         string    `json:"english"`
		Portuguese      string    `json:"portuguese"`
		IsCorrect       bool      `json:"is_correct"`
		ReviewCreatedAt time.Time `json:"review_created_at"`
	}

	for rows.Next() {
		var review struct {
			ID              int       `json:"id"`
			CreatedAt       time.Time `json:"created_at"`
			ActivityName    string    `json:"activity_name"`
			WordID          int       `json:"word_id"`
			English         string    `json:"english"`
			Portuguese      string    `json:"portuguese"`
			IsCorrect       bool      `json:"is_correct"`
			ReviewCreatedAt time.Time `json:"review_created_at"`
		}
		if err := rows.Scan(
			&review.ID,
			&review.CreatedAt,
			&review.ActivityName,
			&review.WordID,
			&review.English,
			&review.Portuguese,
			&review.IsCorrect,
			&review.ReviewCreatedAt,
		); err != nil {
			return nil, 0, err
		}
		reviews = append(reviews, review)
	}

	return reviews, totalItems, rows.Err()
}
