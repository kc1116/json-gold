// +build ignore

package main

import (
    "github.com/kc1116/json-gold/lds"
)


func main() {
	doc := map[string]interface{}{
        "@context":"https://w3id.org/identity/v1",
		"title": "Sample",
        "signature": map[string]interface{}{
			"@type": "LinkedDataSignature2015",
            "creator": "http://example.com/i/pat/keys/5",
            "created": "2011-09-23T20:21:34Z",
            "domain": "github.com/kc1116",
            "nonce": "2bbgh3dgjg2302d-d2b3gi423d42",
            "signatureValue": "OGQzNGVkMzVm4NTIyZTkZDYNmExMgoYzI43Q3ODIyOWM32NjI=",
		},
	}

    opts := lds.SignatureOptions{
        Context: "https://w3id.org/identity/v1",
        Creator: "Khalil Claybon",
        Created: "Today",
        Domain: "github.com/kc1116",
        Nonce: "fdhfakdf89fyd9afg37798",
    }

    var privk = "-----BEGIN RSA PRIVATE KEY-----"+
"MIICWwIBAAKBgGLOw75pbognskzSL6CvjDXbq94wnGk0a58tLsp+sG6kVt4T/lZq"+
"DPmSjlrrBF9sXnC7O6q6hGDyjWIDxAyE0OFC9xqrCb9klxvEHDV1joTV6uxJKXSP"+
"x/241zAsq7mZPM3HqodQzJUcNDENYbeXuP/Pi6UXoDsw8fPQkTJXkFlXAgMBAAEC"+
"gYAI+smHUIWfEhx+Jsv1Sn7vlhs0gi500TLGsJCEDqdyJrVOUXrX16N+Ovd9A8bN"+
"9UdP73Qou/Kz7Nc0hSsYCCoDcd0CDNgJI9zaKGj3xot7VTEXNk1Kr9wERSEndCuu"+
"qO08l0w4RO7CXq7HRpedhdr2nqatDkBkE9uEEptqn0N3QQJBAKo2gwYUW25GdDXS"+
"k6DrW3z+E7QyW2rG2u71wMnEdISnCHKyHvGAqPgHsqJ+MiOOJuOjijDOKlB06tTS"+
"KAUoJK8CQQCUm0mufzDncZPm3itkvpbWEfcKBrN8ZjlDBzvsT4IWpxH9FJ9W6fQo"+
"8DwqorMfjuOC+ad/mhbiOaR8iri2jQ/ZAkEAhJ5QW85EppjyNnVJXNnDwJFd3MpX"+
"e8xQDFshyJLujeRuqp6piVTLUeT9g6l7e0RofHiRVRFs2p8d0I+las8qNQJAJ6Zf"+
"PG23UKlfOwQgM9seR7O3ZDdxgEmOEbJGbMCyBvVAuXPdJ8V4XcvrYbzTaiIn1fRi"+
"mos0e9vBZXFl418z0QJANQKSXyIoHOL2deSafKmPhWlhY2zlqOnmiKJDqzuV2Enh"+
"m57XtNSt7iYq4OI+V8J4dPCLvkllPmn7npbZxun2Hw=="+
"-----END RSA PRIVATE KEY-----"
	lds.S2015(doc, opts, privk)

}
