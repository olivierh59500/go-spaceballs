# Fine contours over the returning block field

The side-track 146 director reuses the 17-band block-gradient program with
six contour banks and 71 source cues. It updates the background's four
corner colors through its own 15-step key-color transport. Each band uses
the same generated 16-step RGB lookup and first 15 horizontal entries as
the earlier block composition.

The source fill helper is an immediate return. Its pose renderer therefore
keeps thin contour lines instead of replacing them with filled silhouettes.
The copper list also addresses the same outline plane one row later in
two other display planes. The shared native block renderer accepts that
one-row sampling shift and keeps the original black line colors.

The cue program contains a counter rewind at one pointer position. Selected
banks, and the later cue region, switch to the mirrored X lookup. These
changes remain in the source controller rather than being approximated by
an independently moving outline object.

## Independent verification

`original-outline-clock.csv` contains all 1,385 source VBL states. Tests
compare the cue/shape counter, bank/frame/remaining count, working/display
addresses, coordinate-table address, four key colors and interpolation phase.
The companion gradient fixture contains 693 independently generated grids:
530,145 RGB12 words, all matching the Go gradient program.

Both original line-writing paths are omitted from the controller execution;
the CPU corner-color program and the hardware-only RGB shift/combine are
retained/evaluated. Native captures show the fine mirrored contours and
changing color blocks. The line-boundary raster comparison and full sequence
alignment remain part of final verification.

```sh
GOWORK=off go run ./cmd/preview -bank outline -music
GOWORK=off go run ./cmd/preview -bank outline -frame 420 -capture .local/outline
```

The [continuous native director](DIRECTOR.md) now assembles this unit with
the remaining source compositions, pages, loading programs and music handoffs.
