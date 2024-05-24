func ConvertToUTF16LEBase64String(command string) string {
    utf16Command := utf16.Encode([]rune(command))
    buffer := new(bytes.Buffer)
    for _, code := range utf16Command {
        buffer.WriteByte(byte(code))
        buffer.WriteByte(byte(code >> 8))
    }
    return base64.StdEncoding.EncodeToString(buffer.Bytes())
}
