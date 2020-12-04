package dataSource

import (
	"ConfigurationTools/model"
	"encoding/xml"
	"io"
	"net/http"
)

func ListCanConfig(url string) (*model.CanXml, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cache := []byte{}
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)

		if n > 0 {
			cache = append(cache, buf[:n]...)
		}

		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
	}

	canCfg := new(model.CanXml)
	err = xml.Unmarshal(cache, canCfg)
	if err != nil {
		return nil, err
	}

	return canCfg, nil
}
