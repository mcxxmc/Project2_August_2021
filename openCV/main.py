from util import (
    capture, bgr2rgb, display, reshape_image_smaller, save_image
)


frame = capture(camera_seq=1)  # camera_seq may change depends on machines
frameRGB = bgr2rgb(frame)  # BGR => RGB
display(frameRGB)
reshaped = reshape_image_smaller(frameRGB, pooling="mid")
display(reshaped)
save_image(reshaped)
