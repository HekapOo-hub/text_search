// Package repository. Full text search 2 approaches: generate additional column to store to_tsvector value and create index on it.
// or create index on to_tsvector and evaluate to verify index matches
// It is possible to
// setweight(to_tsvector(coalesce(title,'')), 'A')    ||
// setweight(to_tsvector(coalesce(body,'')), 'D'); and order by rank then
package repository

import (
	"context"
	"fmt"
	"github.com/HekapOo-hub/textsearch/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"regexp"
	"strings"
)

const (
	russian = "russian"
	english = "english"
)

var (
	rusRegexp = regexp.MustCompile("[а-яА-Я]")
	engRegexp = regexp.MustCompile("[a-zA-Z]")
)

type BooksRepository interface {
	Create(ctx context.Context, b *model.Book) error
	SearchText(ctx context.Context, search string) ([]*model.Book, error)
}

type booksPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewBooksRepository(pool *pgxpool.Pool) *booksPostgresRepository {
	return &booksPostgresRepository{
		pool: pool,
	}
}

func (br *booksPostgresRepository) Create(ctx context.Context, book *model.Book) error {
	configName := determineLanguage(book.Title)
	sql := `INSERT INTO books (author, title, body, config_name) VALUES ($1,$2,$3,$4)`
	_, err := br.pool.Exec(ctx, sql, book.Author, book.Title, book.Body, configName)
	if err != nil {
		return fmt.Errorf("books repository create: %w", err)
	}
	return nil
}

func determineLanguage(s string) string {
	switch {
	case rusRegexp.Match([]byte(s)):
		return russian
	case engRegexp.Match([]byte(s)):
		return english
	default:
		return ""
	}
}

func (br *booksPostgresRepository) SearchText2(ctx context.Context, search string) ([]*model.Book, error) {
	tsQueryForm := getTsQueryForm(search)
	tsVector := "setweight(to_tsvector(config_name, coalesce(author,'')), 'A') || setweight(to_tsvector(config_name, coalesce(title,'')), 'B') || setweight(to_tsvector(config_name, coalesce(body,'')), 'C')"
	sql := fmt.Sprintf(`SELECT author, title, body FROM books, to_tsquery($1) query WHERE %s @@ query ORDER BY ts_rank_cd(%s, query) DESC`, tsVector, tsVector)
	rows, err := br.pool.Query(ctx, sql, tsQueryForm)
	if err != nil {
		return nil, fmt.Errorf("books repository: search text: %w", err)
	}
	var books []*model.Book
	for rows.Next() {
		var b model.Book
		err = rows.Scan(&b.Author, &b.Title, &b.Body)
		if err != nil {
			return nil, fmt.Errorf("books repository: search text: rows scan: %w", err)
		}
		books = append(books, &b)
	}
	return books, nil
}

func getTsQueryForm(s string) string {
	tokens := strings.Split(s, " ")
	return strings.Join(tokens, " | ")
}

func (br *booksPostgresRepository) SearchText(ctx context.Context, search string) ([]*model.Book, error) {
	tsQueryForm := getTsQueryForm(search)
	sql := `SELECT author, title, body FROM books, to_tsquery($1) query WHERE textsearchable_index_col @@ query ORDER BY ts_rank_cd(textsearchable_index_col, query) DESC`
	rows, err := br.pool.Query(ctx, sql, tsQueryForm)
	if err != nil {
		return nil, fmt.Errorf("books repository: search text: %w", err)
	}
	var books []*model.Book
	for rows.Next() {
		var b model.Book
		err = rows.Scan(&b.Author, &b.Title, &b.Body)
		if err != nil {
			return nil, fmt.Errorf("books repository: search text: rows scan: %w", err)
		}
		books = append(books, &b)
	}
	return books, nil
}
