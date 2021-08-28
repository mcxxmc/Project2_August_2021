from util import (
    capture, bgr2rgb, display
)


frame = capture(camera_seq=1)  # camera_seq may change depends on machines
frameRGB = bgr2rgb(frame)  # BGR => RGB
display(frameRGB)
