## rdrop - Command-Line File & Message Sharing Tool

A super-simple command-line utility for quickly sharing a **file**, **text message**, or **plain text file** over your local network.

Supports any combination of:

- 📁 a single file
- 💬 a short message
- 📄 a content file (displayed directly as text)

### How to Use It

**To share a file:**

```bash
rdrop -i ./file.txt
````

**To send a short message:**

```bash
rdrop -m "Here's the meeting link."
```

**To display a content file as plain text:**

```bash
rdrop -I ./note.txt
```

**To combine multiple types:**

```bash
rdrop -i ./report.pdf -m "See the attached report."
```

**To specify a custom port (optional):**

```bash
rdrop -i ./file.txt -p 1130
```

Once running, `rdrop` will display the access URLs in your terminal. Anyone on the same local network can open them in a browser to view or download the shared content.

> ⚠️ `rdrop` is designed for **trusted local networks** only and does **not** include authentication or encryption.
