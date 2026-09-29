# Four-plane ribbons and reversed playback

The side-track 134 director combines two retained, two-bit poses into a
four-plane image. Its source contour table was initially inventoried as
`close-up`; native captures establish its actual role as the angular ribbon
passage. The native effect keeps four working pose pairs and the original
16-color palettes.

The controller first plays frames 0–143 forward with the vertical coordinate
table reversed. It then negates the pointer step, skips five pointer entries
at the handoff and plays backward with the horizontal table reversed. A
later handoff changes direction again and selects the alternate coordinate
table. These source pointer operations, rather than independent reversed
drawings, determine the pose order and orientation.

Two working pairs supply the four visible planes; a different pair receives
the next prepared pose. The native renderer retains these pair selections
and uses a 16-entry palette shader to combine them. The palette program
contains the original black/white flash, three color sets and timed blends.

## Independent verification

`original-ribbons-clock.csv` contains 816 source VBL states. Tests compare
the software count, pose cursor, remaining count, pointer direction,
orientation mode, working/display pair addresses, coordinate-table address
and all 16 live colors. Every field matches the original CPU routines.

Native captures reproduce the layered angular bands with the source poses
and colors. The GPU parity fill and whole-production handoff/music alignment
still need final raster/sequence comparison.

```sh
GOWORK=off go run ./cmd/preview -bank ribbons -music
GOWORK=off go run ./cmd/preview -bank ribbons -frame 420 -capture .local/ribbons
```

The [continuous native director](DIRECTOR.md) now assembles this unit with
the remaining source compositions, pages, loading programs and music handoffs.
