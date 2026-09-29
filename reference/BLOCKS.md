# Copper block gradients and the dancer shadow

The side-track 101 controller at `$50000` reads the 211-frame table loaded
at `$a0000`. It rotates six one-bit masks and displays the current mask
together with the one three positions older. Palette entries 1 and 3 are
black; entry 2 shades the older pose. This produces the black dancer and
the colored shadow without an independently animated shadow object.

The block colors come from two loops of four corner colors at controller
offsets `$a72` and `$f18`. The source builds a 4,096-word interpolation
lookup in memory, then combines its red, green and blue fields with blitter
shifts. Interpolation floors the complete weighted component, unlike the
signed-difference interpolation used by the neighboring fade programs.
The key-color transport uses 15 intermediate steps. Vertical color lookup
uses 16 steps; horizontal copper bands use its first 15 entries.

The copper list has 17 bands and 15 groups of palette writes per band.
Its short band at the vertical-counter wrap repeats the preceding color,
so the native raster joins those adjacent bands. The first band begins
before the visible window and lasts until hardware line 64. Native palette
uploads contain 510 RGBA cells (2,040 bytes), and occur only when a new
software callback changes the color program.

## Verification and remaining timing work

`original-blocks-clock.csv` contains 422 states from the original VBL and
software callbacks. `original-blocks-grids.bin` independently records 212
generated color grids, including initialization: 162,180 RGB12 words. The
original CPU generated the interpolation lookup and executed the gradient
program; its hardware-only shift/combine operation was evaluated from the
source registers. Every recorded color and controller field matches Go.
All contours in this bank are inline byte-coordinate polygons.

Native captures reproduce the black silhouette, its old-pose shadow and
the changing colored blocks. Comparing the same source pose against the
recording also shows that earlier disk-loading handoffs cannot be replaced
by the nominal minimum counter waits alone. Full-production alignment is
still pending; the standalone effect's local clock is verified separately.
The precise copper color-write pixel boundaries also need final comparison.

```sh
GOWORK=off go run ./cmd/preview -bank blocks-effect -music
GOWORK=off go run ./cmd/preview -bank blocks-effect -frame 120 -capture .local/blocks
```

The [continuous native director](DIRECTOR.md) now assembles this unit with
the remaining source compositions, pages, loading programs and music handoffs.
