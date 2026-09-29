# Verified early animation banks

The opening bank contains 142 frames and the first dancer bank 335 frames.
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
Additional animation banks are loaded by the later effect directors and remain
to be decoded before the full production can be reconstructed.
