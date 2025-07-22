`futurego` 

Reference Java CompletableFuture 


#### Test
```go
// fmtln prints a timestamped log message with aligned format.
func fmtln(args ...any) {
	msg := fmt.Sprint(args...)
	fmt.Printf("%-20s %s\n", time.Now().Format("15:04:05"), msg)
}

func main() {
	fmtln("=== Start Execution ===")

	f0 := futurego.Async[string](func() (string, error) {
		<-time.After(30 * time.Second)
		fmtln("➡️ f0 executed after 30s")
		return "hello world", errors.New("simulate an error")
	})

	f1 := futurego.VoidAsync(func() error {
		<-time.After(10 * time.Second)
		fmtln("➡️ f1 executed after 10s")
		return nil
	})

	f2 := futurego.Async(func() (string, error) {
		<-time.After(15 * time.Second)
		fmtln("➡️ f2 executed after 15s")
		return "hello world", nil
	})

	fmtln("Waiting for f1 and f2 to complete...")
	futurego.WaitAll(f1, f2)
	fmtln("✅ f1 and f2 are done")

	f3 := futurego.Async(func() (string, error) {
		<-time.After(2 * time.Second)
		return "➡️ f3 executed after 2s", nil
	})
	fmtln("Waiting for f3")
	f3Res, _ := f3.Get()
	fmtln("✅ f3 result:", f3Res)

	fmtln("🔔Checking f0 status. finished?", f0.IsDone())

	fmtln("Waiting for f0 error result...")
	if err := f0.Error(); err != nil {
		fmtln("⚠️ f0 contains an error:", err.Error())
	}

	fmtln("Thread enters long wait (10h)...")
	<-time.After(10 * time.Hour)
}
```

```markdown
11:52:21             === Start Execution ===
11:52:21             Waiting for f1 and f2 to complete...
11:52:31             ➡️ f1 executed after 10s
11:52:36             ➡️ f2 executed after 15s
11:52:36             ✅ f1 and f2 are done
11:52:36             Waiting for f3
11:52:38             ✅ f3 result:➡️ f3 executed after 2s
11:52:38             🔔Checking f0 status. finished?false
11:52:38             Waiting for f0 error result...
11:52:51             ➡️ f0 executed after 30s
11:52:51             ⚠️ f0 contains an error:simulate an error
11:52:51             Thread enters long wait (10h)...
```