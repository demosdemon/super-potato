package scrape

import (
	"context"
	"io"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

const query = ".JCAZList-list .az-list > ul > li > a"

var re = regexp.MustCompile(`^([^(]+?)(?: \(([^)]+)\))?$`)

type Character struct {
	Name string
	Path string
	Note string
}

type CharacterPage struct {
	outCh chan Character
	errCh chan error
}

func NewCharacterPage(r io.Reader) *CharacterPage {
	outCh := make(chan Character)
	errCh := make(chan error)

	go func() {
		defer close(outCh)
		defer close(errCh)

		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			errCh <- err
			return
		}

		doc.Find(query).Each(func(idx int, s *goquery.Selection) {
			path, _ := s.Attr("href")
			name, note := extractNameNote(s.Text())
			outCh <- Character{
				Name: name,
				Path: path,
				Note: note,
			}
		})
	}()

	return &CharacterPage{outCh, errCh}
}

func (p *CharacterPage) Next(ctx context.Context) (c Character, err error) {
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return c, err
	case c, ok := <-p.outCh:
		if !ok {
			err = io.EOF
		}
		return c, err
	case err, ok := <-p.errCh:
		if !ok {
			err = io.EOF
		}
		return c, err
	}
}

func (p *CharacterPage) Collect(ctx context.Context) ([]Character, error) {
	var rv []Character

	for {
		switch c, err := p.Next(ctx); err {
		case nil:
			rv = append(rv, c)
		case io.EOF:
			return rv, nil
		default:
			return nil, err
		}
	}
}

func extractNameNote(s string) (name, note string) {
	sub := re.FindStringSubmatch(s)
	if len(sub) > 1 {
		name = sub[1]
		note = sub[2]
	}

	return name, note
}
