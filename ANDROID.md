# Android

The Android application runs the same Go production and DCK v1.0.8 as the
desktop command, at 50 ticks per second. Its fixed canvas is fitted in landscape
without stretching. The activity keeps the display awake while the demo runs
and suspends rendering/audio when the application leaves the foreground.

## Build and install

```sh
./scripts/run-android.sh
./scripts/run-android.sh --build-only
```

The script uses Ebitengine and ebitenmobile 2.9.11, Android SDK 36, NDK r28,
Java 17, Gradle 8.11.1 and Android Gradle Plugin 8.10.1. Set `ANDROID_HOME`
and `JAVA_HOME` if the SDK/JDK use different paths. `ANDROID_SERIAL` can select
an authorized USB device when several are connected.

The APK is `android/app/build/outputs/apk/debug/app-debug.apk` and its package
is `com.olivierh.spaceballs`. This development build targets `arm64-v8a` with
minimum Android API 23. Generated libraries and APKs are excluded from Git.

The mobile host creates the production at its first update, after Android's
view/context initialization. Reopening the activity starts a fresh production.
The native Android Back action leaves the application.

## Device checkpoints

```sh
adb shell am start -S -W -n com.olivierh.spaceballs/.MainActivity \
    --ei spaceballs_verify_tick 2200
```

This optional launch extra selects a native production tick and logs scene,
TPS/FPS and Go heap measurements. A positive checkpoint uses a muted seek;
tick zero exercises normal music playback. Checkpoint launches can display
over a locked screen while preserving the device's keyguard security.

## Native video export

```sh
GOWORK=off go run ./cmd/video
```

DCK records the production canvas and its own synchronized PCM output. The
default export is a 250-second H.264/AAC MP4 at 704 by 580 pixels and 50 FPS,
including the closing hold. It also writes a PNG poster and chapter report.
FFmpeg is required. The original recording is not used by this exporter.

## Verified delivery

The ARM64 development APK was installed and exercised on a Google Pixel 10a
on 29 September 2026. The complete sequence reached its stationary closing
page without a crash. After startup, five-second samples measured 49.2–51.0
ticks/s around the 50 Hz target and approximately 60 displayed frames/s.
The measured Go heap peaked at 65.4 MiB; this excludes native/GPU memory.
Backgrounding and resuming the activity retained the same running process.
Both ZIP packaging and all native ELF load segments passed 16 KiB alignment.

The 30 September 2026 package was rebuilt after the shared contour and offset
material migrations. The embedded native library records DCK v1.0.5; its four
ELF load segments and APK library placement pass 16 KiB alignment. Desktop parity
covers 690 sampled complete frames through all fourteen effect units with the
published module. The package was installed on the Pixel 10a on 30 September and
completed the whole sequence through its closing page. Five-second measurements
recorded 49.8–50.9 ticks/s, 57.0–60.1 displayed frames/s and a peak Go heap of
67.0 MiB, excluding native/GPU memory. No crash was observed. This runtime check
does not replace the separate visual-reference comparisons.

The subsequent DCK v1.0.8 ARM64 package includes the shared spatial palette
material. It retains all 690 published-module frame samples and reduces the
Blocks/Outline palette texture to 2,040 bytes. Its four native ELF load segments
and APK library placement pass 16 KiB alignment.

This v1.0.8 package was installed on the Pixel 10a and completed its director
through the closing page. Five-second samples measured 49.2–50.9 ticks/s,
58.9–60.1 displayed frames/s and a peak Go heap of 66.5 MiB without an observed
crash. Go heap excludes native/GPU memory; this is a runtime check rather than
a new hardware visual-reference comparison.

The complete export contains 12,500 frames over 250 seconds at 704 by 580
pixels, with stereo 48 kHz music. The portfolio WebM uses VP9/Opus; the retained
MP4 uses H.264/AAC. Both loading and main music are present. A separate
140-second MP4 excerpt starts at 43.78 seconds for short video sharing.
Browser validation covered French/English copy, poster loading, video/audio
decoding, seeks across the production, search filtering and narrow layouts.
