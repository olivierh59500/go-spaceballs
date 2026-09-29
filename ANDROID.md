# Android

The Android application runs the same Go production and DCK v1.0.0 as the
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
