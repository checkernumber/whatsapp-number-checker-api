// WhatsApp Number Checker API — bulk verification example (C# / .NET 6+).
// Workflow: submit a file of E.164 numbers -> poll status -> download results.
// Docs: https://docs.checknumber.ai/whatsapp-bulk-checker
using System;
using System.IO;
using System.Net.Http;
using System.Text.Json;
using System.Threading.Tasks;

class WhatsAppChecker
{
    const string BaseUrl = "https://api.checknumber.ai";
    const string TaskType = "ws"; // ws | ws_active | ws_avatar
    static readonly string ApiKey =
        Environment.GetEnvironmentVariable("CHECKNUMBER_API_KEY") ?? "YOUR_API_KEY";
    static readonly HttpClient Http = new HttpClient();

    static async Task<JsonElement> PostAsync(string url, MultipartFormDataContent form)
    {
        var req = new HttpRequestMessage(HttpMethod.Post, url) { Content = form };
        req.Headers.Add("X-API-Key", ApiKey);
        var resp = await Http.SendAsync(req);
        resp.EnsureSuccessStatusCode();
        return JsonSerializer.Deserialize<JsonElement>(await resp.Content.ReadAsStringAsync());
    }

    static async Task<JsonElement> SubmitTask(string path)
    {
        var form = new MultipartFormDataContent
        {
            { new ByteArrayContent(File.ReadAllBytes(path)), "file", Path.GetFileName(path) },
            { new StringContent(TaskType), "task_type" }
        };
        return await PostAsync($"{BaseUrl}/v1/tasks", form);
    }

    static async Task<JsonElement> GetTask(string taskId)
    {
        var form = new MultipartFormDataContent { { new StringContent(taskId), "task_id" } };
        return await PostAsync($"{BaseUrl}/v1/gettasks", form);
    }

    static async Task Main()
    {
        var submit = await SubmitTask("numbers.txt");
        var taskId = submit.GetProperty("task_id").GetString();
        Console.WriteLine($"task_id: {taskId}");

        JsonElement task;
        while (true)
        {
            task = await GetTask(taskId);
            var status = task.GetProperty("status").GetString();
            Console.WriteLine($"status={status}");
            if (status == "exported") break;
            if (status == "failed") throw new Exception("task failed");
            await Task.Delay(5000);
        }

        if (task.TryGetProperty("result_url", out var url) && url.GetString() is string u && u.Length > 0)
        {
            File.WriteAllBytes("results.xlsx", await Http.GetByteArrayAsync(u));
            Console.WriteLine("saved to: results.xlsx");
        }
    }
}
