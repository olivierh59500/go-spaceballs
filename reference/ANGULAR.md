# Authored material trails and angular poses

The side-track 114 controller sequences five contour banks through 32 cues.
Its sixth display plane reuses the original textured materials loaded by
side tracks 118–126. The nine-track overlay replaces 50,688 bytes and leaves
352 bytes of the preceding material resident; the native adapter preserves
that tail rather than silently zeroing the fourth plane.

The texture's eight-position counter advances every software callback. The
first five callbacks do not draw a pose. Subsequent callbacks advance a
shape counter, switch cues at multiples of eleven and prepare one of six
retained masks. Five display pointers progressively address older masks;
the trail age resets when the selected bank changes. Thus the source's
short bank changes also alter the visible trail structure.

The original RGB12 palette supplies the dark field, warm/cool pose colors
and half-bright material versions. The native effect shares `maskRing`,
four-bit packing and the six-plane material shader with the earlier
textured-dancer composition. Its angular and rotating shapes come from the
authored contour banks, not a generic replacement animation.

## Independent verification

`original-angular-clock.csv` contains 715 executed original states. Every
software/shape count, bank/frame/remaining count, trail age, working mask,
texture pointer, five display pointers and base palette color matches Go.
The shared GPU composition also retains the six-plane pixel-oracle test.

```sh
GOWORK=off go run ./cmd/preview -bank angular -music
GOWORK=off go run ./cmd/preview -bank angular -frame 250 -capture .local/angular
```

The minimum source wait for this unit is 715 VBL ticks. Its final position
in the uninterrupted production remains part of director/music alignment.
The next source director introduces different coordinate and edge programs;
those must be recovered separately rather than treating all later banks as
identical filled silhouettes.
