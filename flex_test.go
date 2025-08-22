package jsonflex_test

import (
	"errors"
	"testing"

	"github.com/krelinga/go-jsonflex"
)

type Movie jsonflex.Object

func (m Movie) Adult() (bool, error) {
	return jsonflex.GetField(m, "adult", jsonflex.AsBool())
}

func (m Movie) Title() (string, error) {
	return jsonflex.GetField(m, "title", jsonflex.AsString())
}

func (m Movie) ID() (int32, error) {
	return jsonflex.GetField(m, "id", jsonflex.AsInt32())
}

func (m Movie) GenreIDs() ([]int32, error) {
	return jsonflex.GetField(m, "genre_ids", jsonflex.AsArray(jsonflex.AsInt32()))
}

func (m Movie) Genres() ([]Genre, error) {
	return jsonflex.GetField(m, "genres", jsonflex.AsArray(jsonflex.AsObject[Genre]()))
}

type Genre jsonflex.Object

func (g Genre) ID() (int32, error) {
	return jsonflex.GetField(g, "id", jsonflex.AsInt32())
}

func (g Genre) Name() (string, error) {
	return jsonflex.GetField(g, "name", jsonflex.AsString())
}

func assertNoError[T any](in T, err error) func(*testing.T) T {
	return func(t *testing.T) T {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return in
	}
}

func TestFoo(t *testing.T) {
	movie := Movie{
		"adult":     false,
		"title":     "Inception",
		"id":        jsonflex.Number(12345),
		"genre_ids": jsonflex.Array{jsonflex.Number(28), jsonflex.Number(12), jsonflex.Number(878)},
		"genres": jsonflex.Array{
			jsonflex.Object{"id": jsonflex.Number(28), "name": "Action"},
			jsonflex.Object{"id": jsonflex.Number(12), "name": "Adventure"},
			jsonflex.Object{"id": jsonflex.Number(878), "name": "Science Fiction"},
		},
	}

	adult, err := movie.Adult()
	if err != nil || adult {
		t.Errorf("expected adult to be false, got %v with error %v", adult, err)
	}

	title, err := movie.Title()
	if err != nil || title != "Inception" {
		t.Errorf("expected title 'Inception', got '%s' with error %v", title, err)
	}

	id, err := movie.ID()
	if err != nil || id != 12345 {
		t.Errorf("expected id 12345, got %d with error %v", id, err)
	}

	genreIDs, err := movie.GenreIDs()
	if err != nil || len(genreIDs) != 3 || genreIDs[0] != 28 || genreIDs[1] != 12 || genreIDs[2] != 878 {
		t.Errorf("expected genre_ids [28, 12, 878], got %v with error %v", genreIDs, err)
	}

	genres, err := movie.Genres()
	if err != nil || len(genres) != 3 || assertNoError(genres[0].ID())(t) != 28 || assertNoError(genres[1].ID())(t) != 12 || assertNoError(genres[2].ID())(t) != 878 {
		t.Errorf("expected genres with ids [28, 12, 878], got %v with error %v", genres, err)
	}

	genre := Genre(movie)
	id, err = genre.ID()
	if err != nil || id != 12345 {
		t.Errorf("expected genre id 12345, got %d with error %v", id, err)
	}
	_, err = genre.Name()
	if err == nil {
		t.Error("expected error when accessing name on Genre object without 'name' field")
	}
}

func TestBar(t *testing.T) {
	keywords, err := jsonflex.FromArray(jsonflex.Array{
		jsonflex.Object{"id": jsonflex.Number(1), "name": "Action"},
		jsonflex.Object{"id": jsonflex.Number(2), "name": "Adventure"},
		jsonflex.Object{"id": jsonflex.Number(3), "name": "Science Fiction"},
	}, jsonflex.AsObject[Genre]())
	if err != nil {
		t.Fatal("expected keywords to be successfully converted to Genre objects")
	}
	if len(keywords) != 3 {
		t.Fatalf("expected 3 keywords, got %d", len(keywords))
	}
	if assertNoError(keywords[0].ID())(t) != 1 || assertNoError(keywords[1].ID())(t) != 2 || assertNoError(keywords[2].ID())(t) != 3 {
		t.Fatalf("expected keyword IDs to be [1, 2, 3], got [%d, %d, %d]", assertNoError(keywords[0].ID())(t), assertNoError(keywords[1].ID())(t), assertNoError(keywords[2].ID())(t))
	}
	if assertNoError(keywords[0].Name())(t) != "Action" || assertNoError(keywords[1].Name())(t) != "Adventure" || assertNoError(keywords[2].Name())(t) != "Science Fiction" {
		t.Fatalf("expected keyword names to be ['Action', 'Adventure', 'Science Fiction'], got ['%s', '%s', '%s']", assertNoError(keywords[0].Name())(t), assertNoError(keywords[1].Name())(t), assertNoError(keywords[2].Name())(t))
	}
}

func TestErrors(t *testing.T) {
	// Test nil object access
	_, err := jsonflex.GetField(nil, "title", jsonflex.AsString())
	if err == nil || errors.Is(err, jsonflex.ErrFieldNotFound) || errors.Is(err, jsonflex.ErrCannotConvert) {
		t.Errorf("expected non-specific error accessing field on nil object, got %v", err)
	}

	// Test field not found
	movie := Movie{"title": "Test Movie"}
	_, err = jsonflex.GetField(movie, "non_existent_field", jsonflex.AsInt32())
	if err == nil || !errors.Is(err, jsonflex.ErrFieldNotFound) {
		t.Errorf("expected field not found error, got %v", err)
	}

	// Test conversion error
	_, err = jsonflex.GetField(movie, "title", jsonflex.AsBool())
	if err == nil || !errors.Is(err, jsonflex.ErrCannotConvert) {
		t.Errorf("expected conversion error, got %v", err)
	}
}
