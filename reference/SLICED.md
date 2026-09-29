# Sliced dancer, camera zoom and authored cue rewinds

The side-track 127 director loads two contour banks and 71 cue entries.
It retains the preceding textured material and the shared six-mask ring.
Its first three software callbacks do not draw; later callbacks advance the
shape program, choose cues and select normal or close-up coordinate tables.
The original minimum exit wait is 1,618 VBL ticks.

The cue program includes explicit counter rewinds at selected pointer
positions. Four one-shot flags prevent some rewinds from repeating. One
other cue always subtracts three shape ticks. These corrections affect both
the choreography and the 44-position camera cycle, so simply reading every
bank from beginning to end would not preserve the original sequence.

The close-up lookup scales X by 680/256 with a -130 shift, and Y by 550/204
with a -135 shift, clipping both coordinates to the display domain. The
second contour bank also exchanges each edge's endpoint Y coordinates
before line setup. This produces the sliced/glitch-like deformation. The
native renderer therefore fills parity spans from those modified edges,
using DCK's bounded geometry batch, instead of assuming an ordinary polygon.

The palette inverts entries 1–31 at the source's trail-reset phase, while
entry zero alternates between `$104` and `$500`. Its half-bright entries
follow the base RGB12 colors. Near the final counter wait the source also
subtracts `$100` from entry zero on selected VBL ticks; this operation is
kept as word arithmetic rather than replaced with an arbitrary fade.

## Independent verification

`original-sliced-clock.csv` contains all 1,618 executed source states. Tests
compare software/shape counts, bank/frame/remaining count, working-mask
address, material pointer, zoom-table pointer, five display pointers and
all 32 base colors. Every field matches the native controller, including
the cue rewinds and final color arithmetic.

The bitmap material and shader reuse the earlier independently checked
six-plane composition. Native captures now show the close-up glow and
swapped-edge deformation. The parity-span renderer preserves the authored
edge geometry; its precise raster boundaries still need comparison with
the original one-dot Amiga line/fill operation.

```sh
GOWORK=off go run ./cmd/preview -bank sliced -music
GOWORK=off go run ./cmd/preview -bank sliced -frame 930 -capture .local/sliced
```
