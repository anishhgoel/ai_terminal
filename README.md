
### Shall I use fmt.scan or bufio.Reader to read the user input?
- Will use bufio as it allows:
1) Multi-word commands. The terminal might require to process commands ike ls -la or cd Downloads. bufio allows to caputure entire lines, including spaces
2) bufio.Reader offers advanced input handling.
3) bufio.Reader gives more control over error handling and input processing as compared to fmt.Scan.

Since users will likely enter multi-word commands or instructions, bufio.Reader is the better choice
