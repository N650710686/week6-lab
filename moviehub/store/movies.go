package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Movie struct {
	ID         int64    `json:"id"`
	Title      string   `json:"title"`
	TitleTh    *string  `json:"titleTh"`
	Poster     *string  `json:"poster"`
	Year       *int     `json:"year"`
	Genre      *string  `json:"genre"`
	Rating     *float64 `json:"rating"`
	Overview   *string  `json:"overview"`
	OverviewTh *string  `json:"overviewTh"`
}

func (s *Store) SearchMovies(ctx context.Context, q string, genreID int) ([]Movie, error) {
	sql := `SELECT m.id, m.title, m.title_th,
                   CASE WHEN m.poster_path LIKE 'http%' THEN m.poster_path
                        ELSE 'https://image.tmdb.org/t/p/w500' || m.poster_path END,
                   EXTRACT(YEAR FROM m.release_date)::int,
                   (SELECT g.name_en FROM movie_genres mg JOIN genres g ON g.id = mg.genre_id WHERE mg.movie_id = m.id AND mg.position = 0),
                   round(m.vote_average, 1)::float8
            FROM movies m
            WHERE 1=1`

	args := []any{}
	if q != "" {
		args = append(args, "%"+q+"%")
		sql += fmt.Sprintf(" AND (m.title ILIKE $%d OR m.title_th ILIKE $%d)", len(args), len(args))
	}
	if genreID > 0 {
		args = append(args, genreID)
		sql += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM movie_genres WHERE movie_id = m.id AND genre_id = $%d)", len(args))
	}
	sql += " ORDER BY m.popularity DESC LIMIT 100"

	rows, err := s.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Movie{}
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.TitleTh, &m.Poster, &m.Year, &m.Genre, &m.Rating); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *Store) GetMovie(ctx context.Context, id int64) (Movie, error) {
	var m Movie
	err := s.Pool.QueryRow(ctx,
		`SELECT m.id, m.title, m.title_th,
				CASE WHEN m.poster_path LIKE 'http%' THEN m.poster_path
					 ELSE 'https://image.tmdb.org/t/p/w500' || m.poster_path END,
				EXTRACT(YEAR FROM m.release_date)::int,
				(SELECT g.name_en FROM movie_genres mg JOIN genres g ON g.id = mg.genre_id WHERE mg.movie_id = m.id AND mg.position = 0),
				round(m.vote_average, 1)::float8,
				m.overview, m.overview_th
		 FROM movies m WHERE m.id = $1`, id).
		Scan(&m.ID, &m.Title, &m.TitleTh, &m.Poster, &m.Year, &m.Genre, &m.Rating, &m.Overview, &m.OverviewTh)
	if errors.Is(err, pgx.ErrNoRows) {
		return m, ErrNotFound
	}
	return m, err
}
