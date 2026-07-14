# Java example

Requires Java 11+ (uses `java.net.http.HttpClient`, no external dependencies).

```bash
export CHECKNUMBER_API_KEY="YOUR_API_KEY"
java WhatsAppChecker.java   # reads numbers.txt (one E.164 number per line)
```

Submits to `POST /v1/tasks` (`task_type=ws`), polls `POST /v1/gettasks`, prints the result URL.

> This example keeps JSON/multipart handling dependency-free for readability. In production, use a JSON library (Jackson/Gson) and a multipart helper. Full docs: https://docs.checknumber.ai/whatsapp-bulk-checker
