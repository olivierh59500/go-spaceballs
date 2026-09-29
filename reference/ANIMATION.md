# Verified early animation banks

The separately loaded opening bank contains 142 frames and the first mixed
hands/dancer bank 335 frames. A 66-frame hands table is also embedded in the
resident program; it supplies the initial foreground plane.
Their word headers are DBF counts; the following longwords point to frame
payloads relative to the bank. Each frame describes one or more byte contours
or six-byte references that interpolate between authored contours.

Inline commands hold plane masks and explicit Y/X byte pairs. Referenced
commands retain backward/forward offsets, a blend numerator and denominator.
The original routine visits the forward contour first. When the two vertex
counts differ, it resamples the shorter contour using Q3 edge accumulators,
DBF group lengths and signed division, then computes the weighted result.

The Go parser's morph coordinates match an independent execution of the
original `0x3f62a` reader for both complete banks. The fixture contains 318
prepared morph contours, including unequal-count pairs and skipped zero-count
shapes. Rendering still needs to preserve the source coordinate lookup,
bitplane parity/fill convention, palette and placement clocks.

These verified banks are reusable source poses, not video-derived silhouettes.
Later directors load additional tables. The inventory below retains those
payloads; their display composition and choreography still need reconstruction.

## Recovered bank inventory

Following the complete loader chain, including the resident table, recovers
33 distinct byte-contour
tables. `internal/source/banks.go` records their file offsets and frame counts.
Every listed table parses its full pointer set and morph payloads with bounded
coordinates. The names describe the inventory; their final scene roles,
source playback offsets and ordering still require choreography verification.

Several small tables are embedded in effect-code banks. Later bodies are loaded
to the same memory addresses as earlier ones, so they cannot be interpreted
through one static memory image. The recovered transfer order retains those
overlays, including the final reload of the effect code from side tracks 88–89.
The source code finishes on a stationary picture after the final animation.

The [resident opening](OPENING.md) and [first musical composition](PATTERN.md)
now retain their verified bank cursors, display clocks and palette programs.
