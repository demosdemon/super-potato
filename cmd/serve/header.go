package serve

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"sort"
)

type Header struct {
	http.Header
}

func (h Header) Keys() []string {
	keys := make([]string, 0, len(h.Header))
	for k := range h.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (h Header) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, xml.Attr{
		Name: xml.Name{
			Local: "length",
		},
		Value: fmt.Sprintf("%d", len(h.Header)),
	})
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	keys := h.Keys()

	for _, k := range keys {
		v := h.Get(k)
		keyStart := xml.StartElement{
			Name: xml.Name{
				Local: k,
			},
		}
		err := e.EncodeToken(keyStart)
		if err != nil {
			return err
		}

		value, err := url.PathUnescape(v)
		if err != nil {
			value = v
		}

		data := xml.CharData(value)
		err = e.EncodeToken(data)
		if err != nil {
			return err
		}

		err = e.EncodeToken(keyStart.End())
		if err != nil {
			return err
		}
	}

	err = e.EncodeToken(start.End())
	if err != nil {
		return err
	}

	return nil
}
