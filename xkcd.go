// Package xkcd builds and searches an offline index of xkcd comics.
package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Comic struct {
	Num        int
	Title      string
	Transcript string
	Year       string
}

func genJSONURL(comicNum int) string {
	return fmt.Sprintf("https://xkcd.com/%d/info.0.json", comicNum)
}

func fetchLatestNum() (int, error) {
	resp, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("http error: %s", resp.Status)
	}

	var c Comic
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return 0, err
	}
	return c.Num, nil
}

func fetchComic(URL string) (*Comic, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %s", resp.Status)
	}

	var comic Comic
	if err := json.Unmarshal(data, &comic); err != nil {
		return nil, fmt.Errorf("unmarshaling %s: %s", URL, err)
	}
	return &comic, nil
}

type Index struct {
	filename string
}

// NewIndex creates an index that will be stored in the filename.
func NewIndex(filename string) *Index {
	return &Index{filename: filename}
}

// Build fetches all comics from 1 to the latest and saves them as JSON to the
// index filename. If the filename already exists fs.ErrExist is returned and no
// building happens. The maxConcurrency limits the number of concurrent HTTP
// requests.
func (idx *Index) Build(maxConcurrency int) error {
	if _, err := os.Stat(idx.filename); err == nil {
		return fs.ErrExist
	}

	latestComicNum, err := fetchLatestNum()
	if err != nil {
		return fmt.Errorf("getting latest comic number: %v", err)
	}

	comicsChan := make(chan Comic, latestComicNum)

	var tokens = make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	for i := 1; i <= latestComicNum; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tokens <- struct{}{}        // acquire a token
			defer func() { <-tokens }() // release a token
			u := genJSONURL(i)
			c, err := fetchComic(u)
			if err != nil {
				log.Printf("fetching comic %s: %v", u, err)
				return
			}
			comicsChan <- *c
		}(i)
	}

	go func() {
		wg.Wait()
		close(comicsChan)
	}()

	var comics []Comic
	for c := range comicsChan {
		comics = append(comics, c)
	}

	f, err := os.Create(idx.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(comics)
}

// Search searches for comics containing the term in their title or transcript.
// The search is case-insensitive.
func (idx *Index) Search(term string) ([]Comic, error) {
	f, err := os.Open(idx.filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var allComics []Comic
	if err := json.NewDecoder(f).Decode(&allComics); err != nil {
		return nil, err
	}

	var comics []Comic
	term = strings.ToLower(term)
	for _, c := range allComics {
		transcript := strings.ToLower(c.Transcript)
		title := strings.ToLower(c.Title)
		if strings.Contains(transcript, term) || strings.Contains(title, term) {
			comics = append(comics, c)
		}
	}
	return comics, nil
}
