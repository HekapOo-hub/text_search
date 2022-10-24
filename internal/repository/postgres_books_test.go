package repository

import (
	"context"
	"github.com/HekapOo-hub/textsearch/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSearchEnglishSingleResponse(t *testing.T) {
	ctx := context.Background()
	defer func() {
		dbPool.Exec(ctx, "TRUNCATE TABLE books")
	}()
	booksRepository := NewBooksRepository(dbPool)
	b := model.Book{
		Author: "Stone",
		Title:  "How to code better",
		Body:   "i was just copying code and became senior developer",
	}
	err := booksRepository.Create(ctx, &b)
	require.NoError(t, err)

	testCases := []struct {
		search    string
		book      model.Book
		nilResult bool
	}{
		{
			search:    "i want to code like senior",
			book:      b,
			nilResult: false,
		},
		{
			search:    "copy",
			book:      b,
			nilResult: false,
		},
		{
			search:    "alex",
			book:      model.Book{},
			nilResult: true,
		},
	}
	for _, test := range testCases {
		books, err := booksRepository.SearchText(ctx, test.search)
		require.NoError(t, err)
		if test.nilResult {
			var nilBooks []*model.Book
			require.Equal(t, nilBooks, books)
		} else {
			require.Equal(t, test.book, *books[0])
		}
	}
}

func TestSearchRussianSingleResponse(t *testing.T) {
	ctx := context.Background()
	defer func() {
		dbPool.Exec(ctx, "TRUNCATE TABLE books")
	}()
	booksRepository := NewBooksRepository(dbPool)
	b := model.Book{
		Author: "Стоун",
		Title:  "Как написать этот код",
		Body:   "я просто очень долго копировал чужой код и стал сеньором",
	}
	err := booksRepository.Create(ctx, &b)
	require.NoError(t, err)

	testCases := []struct {
		search    string
		book      model.Book
		nilResult bool
	}{
		{
			search:    "копирую код",
			book:      b,
			nilResult: false,
		},
		{
			search:    "сеньор разработчик",
			book:      b,
			nilResult: false,
		},
		{
			search:    "Автор камень",
			book:      model.Book{},
			nilResult: true,
		},
	}
	for _, test := range testCases {
		books, err := booksRepository.SearchText(ctx, test.search)
		require.NoError(t, err)
		if test.nilResult {
			var nilBooks []*model.Book
			require.Equal(t, nilBooks, books)
		} else {
			require.Equal(t, test.book, *books[0])
		}
	}
}
