# Continuous native production director

The root command plays the recovered production continuously: opening,
four title pages, compact credits, illustration, both animated loading
programs, every musical composition and the stationary closing artwork.
`ProgramSegments` records its 50 Hz intervals. DCK's timeline selects them;
an effect is constructed on entry and closed on exit, keeping GPU buffers
bounded to the active composition.

The title/drive/decompression holds are anchored to the supplied recording.
Musical effects retain the independently verified source controllers. The
material sequence begins at tick 2189 (43.78 seconds); later units use the
source minimum counter waits, including the second spirograph's 755 ticks.
The closing picture fades in over 17 steps and remains stationary from
247.26 seconds. The reference continues to 249.754 seconds on that page.

DCK opens both recovered modules and owns decoder selection. The loader
tune starts at tick 292; the main module starts at tick 2206. Sample-envelope
comparison aligns them near 5.85 and 44.12 seconds in the supplied recording.
No reference-movie frames or audio are sampled during native playback.

Indexed page tests retain the recovered artwork and color-register order
before source RGB12 transitions. The hires illustration is fitted using
its nonblank artwork bounds and the recorded viewport aspect. Title and
closing pages fit their source display window in the common viewport.

## Verification

The complete director was traversed through every update/draw to the final
stationary page. Deterministic DCK capture produced 62 frames covering the
middle/end of every segment and the final hold. Contact-sheet inspection
checks that all compositions and transitions are present. Original CPU
fixtures verify each musical controller and both loading programs. The
material shader's independent CPU planar oracle compares every output byte,
including atlas-backed source textures. Device runs exercise both the music
handoff and the final fade/closing-page path.

This is a close native reconstruction. Individual GPU line/fill boundaries
can differ from Amiga one-dot rasterization, and the recorded viewport is
reproduced with page-specific image fitting. Source data, choreography,
palettes and material transforms drive the effects; the movie is solely an
excluded analysis reference.

```sh
GOWORK=off go run .
GOWORK=off go run . -silent
GOWORK=off go run . -start 2189 -ticks 250
GOWORK=off go run . -start 12346 -capture .local/closing-check
```

Escape exits. Reference files, tools, captures, analysis and the hosting
exchange directory remain excluded from Git.
