package dataSource

import (
	"ConfigurationTools/model"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

func ListCanConfig() *model.CanXml {
	resp, err := http.Get(model.CONFIG_URL)
	if err != nil {
		log.Print(fmt.Sprintf("%+v", err))
		return nil
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
				log.Print(fmt.Println("read error:", err))
			}
			break
		}
	}

	canCfg := new(model.CanXml)
	err = xml.Unmarshal(cache, canCfg)
	if err != nil {
		log.Print(fmt.Println("xml反序列化", err))
		return nil
	}

	return canCfg
}
