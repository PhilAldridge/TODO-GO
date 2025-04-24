package lib

// Parse input into command and arguments
func ParseInput(input string) []string {
    var parts []string
    var current string
    inQuotes := false

    for _, char := range input {
        switch char {
        case ' ':
            if inQuotes {
                current += string(char)
            } else if current != "" {
                parts = append(parts, current)
                current = ""
            }
        case '"':
            inQuotes = !inQuotes
        default:
            current += string(char)
        }
    }
    if current != "" {
        parts = append(parts, current)
    }
    return parts
}