# JavaScript (browser) example

Browser-side usage with `fetch` + `FormData`, driven by a file input. Calls a
same-origin backend proxy (`/api/async-check-proxy`) and never holds an API
key client-side; your backend attaches the real key and forwards to
`https://api.checknumber.ai`.

Submits via the proxy (`task_type=ws`), polls the same proxy endpoint, returns
the result URL. For a server-side Node.js version that talks to the real API
directly, see [`../nodejs`](../nodejs). Full docs:
https://docs.checknumber.ai/whatsapp-bulk-checker
