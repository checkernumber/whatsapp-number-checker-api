// WhatsApp Number Checker API — browser/JavaScript example.
// Workflow: submit a file of E.164 numbers -> poll status -> download results.
// Docs: https://docs.checknumber.ai/whatsapp-bulk-checker
//
// Browser code must never hold an API key. This example calls a same-origin
// backend proxy (e.g. /api/async-check-proxy) that your server implements;
// the proxy is responsible for attaching the real API key server-side and
// forwarding to https://api.checknumber.ai.

const TASK_TYPE = "ws"; // ws | ws_active | ws_avatar

async function submitTask(file) {
  const form = new FormData();
  form.append("file", file); // a File/Blob, e.g. from <input type="file">
  form.append("task_type", TASK_TYPE);
  const resp = await fetch("/api/async-check-proxy", { method: "POST", body: form });
  if (!resp.ok) throw new Error(`submit failed: ${resp.status}`);
  return resp.json();
}

async function getTask(taskId) {
  const resp = await fetch(`/api/async-check-proxy?task_id=${encodeURIComponent(taskId)}`);
  if (!resp.ok) throw new Error(`getTask failed: ${resp.status}`);
  return resp.json();
}

async function checkWhatsApp(file) {
  const { task_id } = await submitTask(file);
  console.log("task_id:", task_id);
  for (;;) {
    const task = await getTask(task_id);
    console.log(`status=${task.status} success=${task.success}/${task.total}`);
    if (task.status === "exported") return task.result_url;
    if (task.status === "failed") throw new Error("task failed");
    await new Promise((r) => setTimeout(r, 5000));
  }
}

// Example wiring:
// document.querySelector("#file").addEventListener("change", async (e) => {
//   const url = await checkWhatsApp(e.target.files[0]);
//   console.log("results:", url);
// });
