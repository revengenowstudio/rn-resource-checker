using System.Diagnostics;
using System.Security.Cryptography;
using System.Collections.Concurrent;

// 1. 设置配置
string[] targetDirs = args.Length > 0 ? args[0].Split(',', ';') : ["./"];
string[] extensions = [".mix", ".exe", ".dll", ".ext"];

var sw = Stopwatch.StartNew();
var files = new List<string>();

foreach (var dir in targetDirs)
{
    if (!Directory.Exists(dir)) continue;
    var found = Directory.EnumerateFiles(dir, "*.*", SearchOption.AllDirectories)
        .Where(f => extensions.Contains(Path.GetExtension(f).ToLower()));
    files.AddRange(found);
}

Console.WriteLine($"[INFO] Found {files.Count} files. Computing hashes...");

var results = new ConcurrentBag<(string Path, string Hash)>();

await Parallel.ForEachAsync(files, async (path, ct) =>
{
    try
    {
        // 直接使用静态 HashDataAsync，底层会自动处理 FileStream 的 Buffer 和硬件指令集
        using var stream = File.OpenRead(path);
        byte[] hashBytes = await SHA256.HashDataAsync(stream, ct);
        var hash = Convert.ToHexString(hashBytes).ToLower();
        Console.WriteLine($"[INFO] File {path} : {hash}");
        results.Add((path, hash));
    }
    catch (Exception ex)
    {
        Console.WriteLine($"[ERROR] {path}: {ex.Message}");
    }
});

// 4. 排序并保存
var sortedResults = results.OrderBy(r => r.Path).ToList();
string outName = $"result.{DateTime.Now:yyyyMMdd-HHmmss}.txt";

await File.WriteAllLinesAsync(outName,
    sortedResults.Select(r => $"{r.Path} , hash : {r.Hash}")
    .Append($"Total file numbers : {sortedResults.Count}"));

sw.Stop();
Console.WriteLine($"[INFO] Done in {sw.Elapsed.TotalSeconds:F3}s. Result: {outName}");

// 模拟等待
await Task.Delay(5000);