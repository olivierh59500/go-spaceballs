# Animated tile loading passage

The side-track 109 program combines 21 packed images with a nine-row,
three-column pointer table. Each unpacked image has 3,864 bytes: three
1,288-byte planes, each **224 by 46 pixels** with a 28-byte row. Treating
those bytes as one wide monochrome sprite would lose the original lettering
and colors.

Three tile images appear in vertically adjacent bands beginning at visible
coordinates `(64, 68)`, `(64, 114)` and `(64, 160)`. The copper pointers
restart at hardware lines 96, 142 and 188 and disable the planes at line 234.
The visible window is 352 by 280 pixels; its background follows palette
entry zero. Two original eight-color palettes alternate on each tile update.

Startup executes four transport steps. The VBL handler then advances the
transport every five ticks. The first column advances one cue row each time;
the second and third also advance, with extra steps when the first wraps.
The program's working-buffer branch advances twice and then retains its
last group. The native adapter keeps that observed behavior.

## Independent verification

The original CPU executes all 21 unpacker calls. The complete 81,144-byte
output hashes to `c149c46640bd0f4021eac48e59086fa87ddfa27c67a43279fb54ddef263e939e`;
the Go unpacker reproduces it exactly.

`original-tiles-clock.csv` contains 143 executed VBL states. Tests compare
the cue positions, update phase, working/display pointers and every palette
entry. They also compare the SHA-256 of each bitmap copied into the three
displayed bands against the tile selected by the Go controller.

The native effect uses predecoded images for both palettes. DCK supplies the
effect interface and automatically selected module playback; each frame
only chooses three images and their original positions.

```sh
GOWORK=off go run ./cmd/preview -bank tiles -music
GOWORK=off go run ./cmd/preview -bank tiles -frame 70 -capture .local/tiles
```
