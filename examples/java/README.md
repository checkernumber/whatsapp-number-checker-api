# Java example

Requires Java 11+ (uses `java.net.http.HttpClient`, no external dependencies).

From the repository root:

```bash
export CHECKNUMBER_API_KEY="YOUR_API_KEY"
cd examples
java java/WhatsAppChecker.java
```

Submits to `POST /v1/tasks` (`task_type=ws`), polls `POST /v1/gettasks`, and downloads the result to `results.zip`.

> This example keeps JSON/multipart handling dependency-free for readability. In production, use a JSON library (Jackson/Gson) and a multipart helper. Full docs: https://docs.checknumber.ai/whatsapp-bulk-checker
