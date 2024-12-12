package conversions

import (
	"encoding/base64"
	"encoding/binary"
	"unicode/utf16"
)

func ConvertToUTF16LEString(input string) (output string) {
	enc := ConvertToUTF16Slice(input)
	return string(ConvertUTF16ToLEBytes(enc))
}

func ConvertToUTF16LEBase64String(input string) (output string) {
	enc := ConvertToUTF16Slice(input)
	b := ConvertUTF16ToLEBytes(enc)
	return base64.StdEncoding.EncodeToString(b)
}

func ConvertToUTF16Slice(input string) (output []uint16) {
	return utf16.Encode([]rune(input))
}

func ConvertUTF16ToLEBytes(input []uint16) (output []byte) {
	b := make([]byte, 2*len(input)) // 2 bytes per rune
	for idx, r := range input {
		binary.LittleEndian.PutUint16(b[idx*2:], r)
	}
	return b
}