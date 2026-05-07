package main

type LinkStore struct {
	db *pgx.Pool
}

func (s *LinkStore) GetNextLink(ctx context.Content) (string, error) {
	var url string
	query := "SELECT url FROM links WHERE status = 'pending' LIMIT 1"
	err := s.db.QeryRow(ctx, query).Scan(&url)
	return url, err
}

func saveBatch(db *pgx.Pool, links []string) error {
	rows := [][]interface{}{}
	for _, link := range links {
		rows = append(rows, []interface{}{link, "pending"})
	}

	_, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"links"},
		[]string{"url", "status"},
		pgx.CopyFromRows(rows),
	)

	return err
}
