package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"git.sr.ht/~ft/cue"
	"github.com/asdfsx/gochardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
)

func decoderToUTF8For(b []byte) (dec *encoding.Decoder, err error) {
	var best *chardet.Result
	var enc encoding.Encoding
	if best, err = chardet.NewTextDetector().DetectBest(b); err != nil {
		return
	} else if best.Confidence >= 25 {
		if enc, err = ianaindex.IANA.Encoding(best.Charset); err != nil {
			err = nil
		}
	}
	if enc == nil {
		enc = encoding.Nop
	}
	dec = enc.NewDecoder()
	return
}

func openFileUTF8(path string) (r io.ReadCloser, err error) {
	if r, err = os.Open(path); err != nil {
		return
	}
	defer r.Close()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(r); err != nil {
		return
	}
	var dec *encoding.Decoder
	if dec, err = decoderToUTF8For(buf.Bytes()); err != nil {
		return
	}
	r = ioutil.NopCloser(dec.Reader(buf))
	return
}

func cueSheetFromBytes(raw []byte) (sheet *cue.Sheet, err error) {
	if sheet, err = cue.Parse(bytes.NewBuffer(raw), 0); err != nil {
		return
	}

	buf := new(bytes.Buffer)
	buf.WriteString(sheet.Performer)
	buf.WriteString(sheet.Songwriter)
	buf.WriteString(sheet.Title)
	for _, f := range sheet.Files {
		if f.Type == cue.FileTypeWave {
			for _, t := range f.Tracks {
				buf.WriteString(t.Title)
			}
			buf.WriteString(f.Name)
		}
	}
	var dec *encoding.Decoder
	if dec, err = decoderToUTF8For(buf.Bytes()); err == nil {
		sheet, err = cue.Parse(dec.Reader(bytes.NewBuffer(raw)), 0)
	}
	return
}

func pathReplaceChars(s string) string {
	return strings.ReplaceAll(s, "/", "∕")
}
