# JavaScript (browser) example

Browser-side usage with `fetch` + `FormData`, driven by a file input.

> **Security:** never expose your API key in production client-side code. Proxy these requests through your own backend. This example shows the raw request shape only.

Submits to `POST /v1/tasks` (`task_type=ws`), polls `POST /v1/gettasks`, returns the result URL. For a server-side Node.js version see [`../nodejs`](../nodejs). Full docs: https://docs.checknumber.ai/whatsapp-bulk-checker
