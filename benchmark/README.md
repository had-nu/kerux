# Kerux Packet Benchmark

## Structure

```
benchmark/
├── spec_projeto.md          ← Control spec (fixed, identical for both arms)
├── transition-log.csv       ← Token measurements per arm per run
├── sigcheck-arm-j/          ← ARM-J deliverable (JSON packets)
├── sigcheck-arm-m/          ← ARM-M deliverable (marker packets)
└── report.md                ← Final benchmark report (post-experiment)
```

## Execution Order
1. ARM-J Run 1 → ARM-J Run 2 (fresh sessions)
2. ARM-M Run 1 → ARM-M Run 2 (fresh sessions)

## Token Counting
Estimated via tokenizer (no API access).

## Reference
See `KERUX_PACKET_BENCHMARK_SPEC.md` for full experiment design.
