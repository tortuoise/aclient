package aclient

import (
	"bufio"
	_ "encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"
	_ "unicode/utf8"
)

const tmpl = `           DEC | HEX | UNI | CHAR  
                {{range .}} 
                        {{.Dec}} | {{.Hex}} | {{.Uni}} | {{.Char}} 
                {{end}} `

type Odata struct {
	Dec  uint16
	Hex  string
	Uni  string
	Char string
}

// Display displays unicode table in decimal, hex, unicode, string within given range
func Display(start uint16, stop uint16, rt ...*unicode.RangeTable) {

	//chars := make([]rune,0)
	//s := ""
	out := make([]Odata, 0)

	if len(rt) == 0 {

		for j := start; j < stop; j += 1 {

			//chars = append(chars, rune(j))
			//s += string(rune(j)) + " "
			out = append(out, Odata{Dec: j, Hex: fmt.Sprintf("%x", j), Uni: fmt.Sprintf("%U", j), Char: string(rune(j))})
		}

	} else {
		for _, r16 := range rt[0].R16 {

			for j := r16.Lo; j < r16.Hi; j += r16.Stride {

				//chars = append(chars, rune(j))
				//s += string(rune(j)) + " "
				out = append(out, Odata{Dec: j, Hex: fmt.Sprintf("%x", j), Uni: fmt.Sprintf("%U", j), Char: string(rune(j))})
			}

		}
	}

	t := template.Must(template.New("out").Parse(tmpl))
	//err := t.Execute(os.Stdout, string(chars))
	err := t.Execute(os.Stdout, out)
	if err != nil {
		_, _ = io.Copy(os.Stdout, strings.NewReader(err.Error()))
	}

}

// Creates html unicode table in decimal, hex, unicode, string within given range
func Create(start uint16, stop uint16, rt ...*unicode.RangeTable) error {

	out := make([]Odata, 0)

	if len(rt) == 0 {

		for j := start; j < stop; j += 1 {

			out = append(out, Odata{Dec: j, Hex: fmt.Sprintf("%x", j), Uni: fmt.Sprintf("%U", j), Char: string(rune(j))})
		}

	} else {
		for _, r16 := range rt[0].R16 {

			for j := r16.Lo; j < r16.Hi; j += r16.Stride {

				out = append(out, Odata{Dec: j, Hex: fmt.Sprintf("%x", j), Uni: fmt.Sprintf("%U", j), Char: string(rune(j))})
			}

		}
	}

	f, err := os.Create("ascii.html")
	if err != nil {
		return nil
	}
	defer f.Close()

	by := bufio.NewWriter(f)
	t := template.Must(template.ParseFiles("base_ascii.html"))
	err = t.ExecuteTemplate(by, "base_ascii", out)
	if err != nil {
		return err
	}

	if err = by.Flush(); err != nil {
		return err
	}
	return nil

}
