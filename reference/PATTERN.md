# First musical pattern composition

The first musical effect is the side-track 71 controller at `$50000`. It
unpacks the beginning of the side-track 72 data to `$6cb10`, then uses two
independent views of that same one-bit bitmap as display planes 1 and 2.
The bitmap is **640 by 512**, with an 80-byte row; it is not a pair of
640-by-256 images. The visible body mask supplies display plane 0.

The X motion table is 500 words at bank offset `$80e0`; the Y table is 477
words at `$84c8`. Initialization advances the pointers once. Each PAL VBL
then advances the first view by three X entries and two Y entries, and the
second view by three Y entries. Its X coordinate stays fixed. These integer
tables, including their modulo lengths, are read directly from the asset.

The copper list starts display fetch at `$20`, uses DIWSTRT `$1c71` and
DIWSTOP `$3ec7`, and exposes 342 by 290 pixels. The fine-scroll high nibble
is `(~X & 15)`. The resulting material origins are `(X + 17, Y)` and
`(32, secondY)`; the mask is sampled 32 pixels into the fetched row. The
native preview fits this visible region in its existing 352-by-290 canvas.

## Two independent clocks

The original VBL entry at `$50274` updates the material motion at 50 Hz. It
requests the level-1 software interrupt on odd VBL counters. The callback at `$502e4`
therefore advances the silhouettes at 25 Hz. Its first bank has 26 frames;
the next has 210. The triple-buffer display trails the prepared contour by
two software callbacks. The controller exits at shape count 240, after 480
VBL ticks (9.6 seconds). Completed contours are held between shape ticks;
the backdrop continues moving every VBL.

The palette is initially black. Shape count 3 installs the eight base
colors. VBL counters 20 through 36 interpolate the three body colors toward
`$d00`, `$d70`, and `$bb0`; counters 460 through 476 fade the live palette
to black. Each RGB nibble uses signed integer multiplication/division.

## Independent verification

`internal/source/testdata/original-pattern-clock.csv` contains 480 states
from executing both original callback entry points with Ghidra's 68000
interpreter. Only hardware raster-writing helpers were replaced by returns;
the contour reader, bank switch, pointer transport and palette routines ran
unchanged. Tests compare every clock, bank cursor, material origin,
fine-scroll word and live palette entry against that fixture.

The Go effect reads the recovered bitmap and contour tables, renders the
parity mask through DCK's bounded triangle batch, and combines the three
planes with a small Ebitengine shader. It allocates its targets and uniform
storage once. The reference movie is never used as runtime data.

```sh
GOWORK=off go run ./cmd/preview -bank pattern
GOWORK=off go run ./cmd/preview -bank pattern -frame 277 -capture .local/pattern
```

The standalone preview loops this one effect. Whole-production music
alignment and scene handoffs are still being reconstructed. RGB12 colors
and movement are source-derived; the GPU parity fill still needs a full
boundary comparison with the original line/fill rasterizer.

Hardware interpretation follows the original register writes and the
[Commodore Amiga Hardware Reference Manual](https://www.theflatnet.de/pub/cbm/amiga/AmigaDevDocs/hard_6.html).
