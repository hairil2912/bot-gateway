package router

import (
	. "github.com/projectriri/bot-gateway/types"
	"regexp"
	"strings"
)

func route() {
	for {
		pkt := <-producerBuffer

		from := strings.ToLower(pkt.Head.From)
		to := strings.ToLower(pkt.Head.To)

		for _, cc := range consumerChannelPool {
			go func() {

				var formats []Format
				for _, ac := range cc.Accept {
					f, _ := regexp.MatchString(ac.From, from)
					t, _ := regexp.MatchString(ac.To, to)
					if f && t {
						formats = ac.Formats
						break
					}
				}

				if formats == nil {
					return
				}

				for _, format := range formats {

					if strings.ToLower(pkt.Head.Format.API) == strings.ToLower(format.API) &&
						strings.ToLower(pkt.Head.Format.Method) == strings.ToLower(format.Method) &&
						strings.ToLower(pkt.Head.Format.Protocol) == strings.ToLower(format.Protocol) {
						*cc.Buffer <- pkt
						return
					}

					for _, cvt := range converters {
						if cvt.IsConvertible(pkt.Head.Format, format) {
							ok, result := cvt.Convert(pkt, format)
							if ok && result != nil {
								for _, p := range result {
									*cc.Buffer <- p
								}
								return
							}
						}
					}
				}
			}()
		}
	}
}
