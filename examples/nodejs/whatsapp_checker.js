// WhatsApp Number Checker API — bulk verification example (Node.js 18+).
// Workflow: submit a file of E.164 numbers -> poll status -> download results.
// Docs: https://docs.checknumber.ai/whatsapp-bulk-checker
import fs from "node:fs";
import { setTimeout as sleep } from "node:timers/promises";

const BASE_URL = process.env.API_BASE_URL || "https://api.checknumber.ai";
const API_KEY = process.env.CHECKNUMBER_API_KEY;
if (!API_KEY) {
  throw new Error("Set the CHECKNUMBER_API_KEY environment variable");
}
const TASK_TYPE = "ws"; // ws | ws_active | ws_avatar

async function submitTask(filePath) {
  const form = new FormData();
  form.append("file", new Blob([fs.readFileSync(filePath)]), "numbers.txt");
  form.append("task_type", TASK_TYPE);
  const resp = await fetch(`${BASE_URL}/v1/tasks`, {
    method: "POST",
    headers: { "X-API-Key": API_KEY },
    body: form,
  });
  if (!resp.ok) throw new Error(`submit failed: ${resp.status}`);
  return resp.json();
}

async function getTask(taskId) {
  const form = new FormData();
  form.append("task_id", taskId);
  const resp = await fetch(`${BASE_URL}/v1/gettasks`, {
    method: "POST",
    headers: { "X-API-Key": API_KEY },
    body: form,
  });
  if (!resp.ok) throw new Error(`getTask failed: ${resp.status}`);
  return resp.json();
}

async function poll(taskId, intervalMs = 5000) {
  for (;;) {
    const task = await getTask(taskId);
    console.log(`status=${task.status} success=${task.success}/${task.total}`);
    if (task.status === "exported") return task;
    if (task.status === "failed") throw new Error("task failed");
    await sleep(intervalMs);
  }
}

async function download(resultUrl, outPath = "results.zip") {
  const resp = await fetch(resultUrl);
  if (!resp.ok) throw new Error(`download failed: ${resp.status}`);
  fs.writeFileSync(outPath, Buffer.from(await resp.arrayBuffer()));
  return outPath;
}

const task = await submitTask("numbers.txt");
console.log("task_id:", task.task_id);
const final = await poll(task.task_id);
if (final.result_url) console.log("saved to:", await download(final.result_url));
