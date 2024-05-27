package main

import (
	"errors"
	"flag"
	"os"
	"slices"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/saintfish/chardet"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 3 {
		print("Please specify the text file path and encoding source and destination")
		return
	}

	senc := args[0]
	denc := args[1]
	path := args[2]

	if _, err := os.Stat(path); err != nil {
		print("The specified file does not exist")
		return
	}

	if !ValidEncodes(senc, denc) {
		print("Please specify a valid encoding")
		return
	}

	err := convertEncodeFile(senc, denc, path)
	if err != nil {
		print(err)
	}
}

func ValidEncodes(encs ...string) bool {
	list := []string{"auto", "utf-8", "shift-jis"}

	for _, e := range encs {
		if slices.Contains(list, e) {
			return true
		}
	}

	return false
}

func convertEncodeFile(senc string, denc string, path string) error {
	bin, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if senc == "auto" {
		detector := chardet.NewTextDetector()
		rs, _ := detector.DetectAll(bin)

		for _, r := range rs {
			candidate := strings.ToLower(r.Charset)
			candidate = strings.Replace(candidate, "_", "-", -1)
			if ValidEncodes(candidate) {
				senc = candidate
				break
			}
		}

		if senc == "auto" {
			return errors.New("file character code is not supported")
		}
	}

	sbin, err := convertString(senc, bin, Decode)
	if err != nil {
		return err
	}

	dbin, err := convertString(denc, sbin, Encode)
	if err != nil {
		return err
	}

	os.WriteFile(path, dbin, os.ModePerm)
	return nil
}

type ConvertEnum int

const (
	Encode ConvertEnum = iota
	Decode
)

func convertString(enc string, bin []byte, convert ConvertEnum) ([]byte, error) {
	convs := map[string]encoding.Encoding{
		"shift-jis": japanese.ShiftJIS,
	}

	if _, exist := convs[enc]; !exist {
		return bin, nil
	}

	t := lo.Ternary[transform.Transformer](
		convert == Encode,
		convs[enc].NewEncoder(),
		convs[enc].NewDecoder())

	s, _, err := transform.Bytes(t, bin)
	return s, err
}
