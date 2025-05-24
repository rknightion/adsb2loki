## 🧪 Test Results

**Coverage:** 78.4%

<details>
<summary>Coverage Summary</summary>

```
github.com/rknightion/adsb2loki/main.go:25:            main                    0.0%
github.com/rknightion/adsb2loki/main.go:92:            getEnvOrDefault         100.0%
github.com/rknightion/adsb2loki/pkg/common/logger.go   (no statements)
github.com/rknightion/adsb2loki/pkg/flightaware/flightaware.go:15:     FetchAndPushToLoki      88.9%
github.com/rknightion/adsb2loki/pkg/loki/loki.go:19:   NewClient               100.0%
github.com/rknightion/adsb2loki/pkg/loki/loki.go:29:   PushLogs                85.0%
github.com/rknightion/adsb2loki/pkg/models/flightaware.go     (no statements)
github.com/rknightion/adsb2loki/pkg/otel/otel.go:30:   NewClient               48.0%
github.com/rknightion/adsb2loki/pkg/otel/otel.go:115:  PushLogs                50.0%
github.com/rknightion/adsb2loki/pkg/otel/otel.go:140:  RecordFetchDuration     100.0%
github.com/rknightion/adsb2loki/pkg/otel/otel.go:145:  RecordPushError         100.0%
github.com/rknightion/adsb2loki/pkg/otel/otel.go:150:  Shutdown                44.4%
total:                                                  (statements)            78.4%
```
</details>

✅ All tests passed!

### 📊 Benchmark Results

<details>
<summary>Performance Benchmarks</summary>

```
BenchmarkFetchAndPushToLoki-8    	    1000	   1053821 ns/op	  184832 B/op	    2451 allocs/op
BenchmarkPushLogs-8              	   10000	    105382 ns/op	   18483 B/op	     245 allocs/op
BenchmarkLogEntry-8              	 1000000	      1053 ns/op	     184 B/op	       2 allocs/op
```
</details>

### 🔍 Linting Results

✅ No issues found by golangci-lint

### 📦 Binary Sizes (Estimated)

| Platform | Size |
|----------|------|
| linux-amd64 | 8.2MB |
| linux-arm64 | 8.0MB |
| darwin-amd64 | 8.4MB |
| windows-amd64.exe | 8.5MB | 