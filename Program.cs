using System.Diagnostics;
using System.Security.Cryptography;
using System.Collections.Concurrent;
using System.Text;

var waitSecond = 5;
List<string> targetDirs = ["./"];
List<string> extensions = [".mix", ".exe", ".dll", ".ext"];
bool showHelp = false;

for (int i = 0; i < args.Length; i++)
{
    switch (args[i])
    {
        case "--path" or "-p":
            if (i + 1 < args.Length) targetDirs = args[++i].Split(',', ';').Select(d => d.Trim()).ToList();
            break;
        case "--suffix" or "-s":
            if (i + 1 < args.Length) extensions = args[++i].Split(',', ';').Select(e => "." + e.Trim().TrimStart('.')).ToList();
            break;
        case "--help" or "-h":
            showHelp = true;
            break;
    }
}

if (showHelp)
{
    Console.WriteLine("""
    NAME:
       rn-resource-checker - scan paths and compute SHA256 hashes for specific suffixes

    USAGE:
       rn-resource-checker [global options]

    GLOBAL OPTIONS:
       --path string, -p string    Specify paths, separated by ',' or ';' (default: "./")
       --suffix string, -s string  Specify suffixes (default: "mix", "exe", "dll", "ext")
       --help, -h                  Show help
    """);
    return;
}

var finalExts = extensions.Select(e => e.StartsWith('.') ? e.ToLower() : "." + e.ToLower()).ToArray();

var sw = Stopwatch.StartNew();
var files = new List<string>();

foreach (var dir in targetDirs)
{
    if (!Directory.Exists(dir)) continue;
    var found = Directory.EnumerateFiles(dir, "*.*", SearchOption.AllDirectories)
        .Where(f => finalExts.Contains(Path.GetExtension(f).ToLower()));
    files.AddRange(found);
}

Console.WriteLine($"[INFO] Target paths: {string.Join(", ", targetDirs)}");
Console.WriteLine($"[INFO] Suffix whitelist: {string.Join(", ", finalExts)}");
Console.WriteLine($"[INFO] Found {files.Count} files. Computing hashes...");

var results = new ConcurrentBag<(string Path, string Hash)>();

await Parallel.ForEachAsync(files, async (path, ct) =>
{
    try
    {
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

string versionContent = "";
string versionPath = Path.Combine(targetDirs[0], "Versioncode");
if (File.Exists(versionPath))
{
    Console.WriteLine($"[INFO] Found Versioncode from: {versionPath}");
    versionContent = await File.ReadAllTextAsync(versionPath);
}
else
{
    Console.WriteLine($"[WARN] Versioncode file not found at: {versionPath}");
}

var sortedResults = results.OrderBy(r => r.Path).ToList();
string outName = $"result.{DateTime.Now:yyyyMMdd-HHmmss}.txt";

using (var swOut = new StreamWriter(outName, false, Encoding.UTF8))
{
    if (!string.IsNullOrEmpty(versionContent))
    {
        await swOut.WriteLineAsync(versionContent);
    }
    foreach (var r in sortedResults)
    {
        await swOut.WriteLineAsync($"{r.Path} , hash : {r.Hash}");
    }
    await swOut.WriteLineAsync($"Total file numbers : {sortedResults.Count}");
}

sw.Stop();
Console.WriteLine($"[INFO] Done in {sw.Elapsed.TotalSeconds:F3}s. Result: {outName}");

Console.WriteLine($"Exit after {waitSecond}s ...");
await Task.Delay(waitSecond * 1000);