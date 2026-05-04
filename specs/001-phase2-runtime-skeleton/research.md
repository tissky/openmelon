# Research: OpenMelon Phase 2 — Runtime Skeleton

**Date**: 2026-05-04  
**Branch**: `001-phase2-runtime-skeleton`

## Q1: Go Subcommand CLI — flag vs cobra vs urfave/cli

**Decision**: 继续使用标准库 `flag` 包，引入 `run` 子命令用 `flag.NewFlagSet` 实现

**Rationale**: Phase 2 只有 1 个子命令 (`run`)，保持零外部依赖。`flag.NewFlagSet` 支持多子命令，无需引入 cobra。Phase 3+ 如果需要 5+ 子命令再迁移 cobra。

**Alternatives considered**:
- Cobra: 功能完整，但 +50KB 二进制 + 1 依赖，Phase 2 不值得
- urfave/cli: 比 cobra 轻量，但仍是不必要的依赖

**Pattern**:
```go
// main.go
switch os.Args[1] {
case "run":
    fs := flag.NewFlagSet("run", flag.ExitOnError)
    project := fs.String("project", "", "path to project.json")
    // ... more flags
    fs.Parse(os.Args[2:])
    // delegate to engine.Run(...)
}
```

---

## Q2: Python Subprocess Management

**Decision**: `exec.CommandContext(ctx, "python3", ...)` with separated stdout/stderr capture，先验证 python3 可用

**Rationale**: context-based timeout 让调用方控制超时；stdout/stderr 分离让错误信息清晰；LookPath 验证提供好的错误提示。

**Alternatives considered**:
- `sh -lc` 包装：能继承 shell 环境但错误诊断困难，当前 main.go 用的是这种方式
- 不做超时：Phase 2 用 shell provider，LLM 可能挂住，必须加超时

**Pattern**:
```go
func (c *Compiler) Compile(ctx context.Context, req *CompileRequest) (*CompiledSkill, error) {
    if _, err := exec.LookPath("python3"); err != nil {
        return nil, fmt.Errorf("python3 not found in PATH — install Python 3.9+: %w", err)
    }
    args := []string{"-m", "skillplus_compile", req.PackagePath, "--target", req.Target, ...}
    cmd := exec.CommandContext(ctx, "python3", args...)
    cmd.Env = append(os.Environ(), "PYTHONPATH="+filepath.Clean(c.CompilerPath))
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("skillplus compile failed: %w\ndetail: %s", err, stderr.String())
    }
    // unmarshal stdout JSON → CompiledSkill
}
```

---

## Q3: Provider Interface Design

**Decision**: 2-method interface，接受 `context.Context`，返回 typed `ProviderError`

**Rationale**: context 是 Go 超时/取消的标准方式，typed error 允许调用方区分 timeout vs auth_failed vs rate_limit。Phase 2 实现 `ShellProvider`，Phase 3 添加 `HTTPProvider`。

**Alternatives considered**:
- 单纯函数 `func(Request) (string, error)`：无法添加超时，扩展性差
- 大接口 10+ 方法：违反接口隔离，过度设计

**Interface**:
```go
type Provider interface {
    Generate(ctx context.Context, req *Request) (content string, trace *Trace, err error)
}
type Trace struct {
    Provider string; Model string; Command string; DurationSec float64
}
type ProviderError struct {
    Code string // "timeout" | "non_zero_exit" | "empty_output"
    Wrapped error
}
```

---

## Q4: Provenance JSONL Append

**Decision**: `os.OpenFile(O_APPEND|O_WRONLY|O_CREATE)` + `json.Marshal` + `\n` + `f.Sync()`

**Rationale**: 单进程无并发，简单 append 模式即可。`f.Sync()` 确保崩溃不丢记录。不需要原子重命名，这是追加不是替换。

**Alternatives considered**:
- 原子 temp file rename：适合替换场景；追加场景下逻辑更复杂且无收益
- SQLite：过重，JSONL 是追加审计日志的惯用模式

**Pattern**:
```go
func AppendRecord(path string, rec *Record) error {
    f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
    if err != nil { return fmt.Errorf("open provenance file: %w", err) }
    defer f.Close()
    data, _ := json.Marshal(rec)
    f.Write(data); f.WriteString("\n"); f.Sync()
    return nil
}
```

---

## Q5: Artifact ID Stability

**Decision**: 保留现有 SHA256 截断方案（`stableID`），输入为 `project_id:workflow:stage:package_id:intent_hash`

**Rationale**: 可重现（相同输入→相同 ID），64bit 碰撞概率可忽略，与 git short-SHA 一致的惯用模式。

**改进点**: intent 参与哈希，区分同一 workflow 的不同意图产出。

**Alternatives considered**:
- UUID v4：随机，无法去重或缓存
- 单调计数器：跨运行不稳定
