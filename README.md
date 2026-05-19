

## Project overview

Go audio playback library & CLI built on [miniaudio](https://miniaud.io/) via Go CGO bindings [malgo](https://github.com/gen2brain/malgo). Supports WAV (PCM) playback, device enumeration, and low-latency notification sound replay.

## Build & run

```bash
go build ./...              # build
go vet ./...                # lint
go run main.go              # list devices (default)
go run main.go list         # list devices
go run main.go play <wav> [deviceID]  # play a WAV file
```

`malgo` wraps the C library `miniaudio` → building requires a C compiler toolchain (GCC/MinGW on Windows, `build-essential` on Linux, Xcode CLI on macOS).

## Architecture

```
audio/          ← reusable package (import by third-party Go programs)
  device.go     → DeviceInfo, ListDevices() standalone convenience
  engine.go     → AudioPlayerEngine: Init/Destroy, Preload, PlaySound, NewPlayer, ListDevices
  player.go     → Player: Play/Stop/Close/Reset/Replay/Done, WAV header parsing (stdlib encoding/binary)
  errors.go     → ErrNotInitialized
main.go          ← CLI demo
```

### Package API

**Standalone (simple, one-shot):**
```go
devices, _ := audio.ListDevices()
player, _ := audio.NewPlayer("file.wav", "")  // "" = default device
player.Play()
<-player.Done()
player.Close()
```

**Engine (shared context, sound caching, low-latency replay):**
```go
var engine audio.AudioPlayerEngine
engine.Init()
defer engine.Destroy()

engine.Preload("beep", "beep.wav")           // decode once, cache PCM
player, _ := engine.PlaySound("beep", "")    // 1st call: init device + play
<-player.Done()
player, _ = engine.PlaySound("beep", "")     // subsequent: Replay() — no device init
```

**Key Player methods:**
- `Play()` — start (non-blocking)
- `Stop()` — pause, device stays initialized → can `Play()` again
- `Reset()` — rewind to beginning
- `Replay()` — Stop + Reset + Play in one call
- `Close()` — stop + uninit device (one-shot players)
- `Done()` — channel that closes when playback reaches end

### WAV support

Parses standard PCM WAV headers via `encoding/binary` (no external dependency). Supports 8/16/24/32-bit, mono/stereo, any sample rate. Non-PCM or compressed WAV formats return an error.

### Engine lifecycle

- `Engine.Destroy()` closes all Players (both cached and standalone) before freeing the malgo context
- The `"name|deviceID"` cache key means: same sound on different devices → separate cached Players
- Context is shared across all Players from the same Engine
