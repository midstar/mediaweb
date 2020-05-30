# MediaWEB performance test

## Test setup

Following MediaWEB configuration parameters where enabled:

    genthumbsonstartup = on
    enablepreview = on
    previewmaxside = 1280
    genpreviewonstartup = on

The mediapath included 50 images of which:

* 9 required thumbnail generation (no EXIF and larger than thumb size)
* 43 required preview generation (larger than 1280 px)
* Average size of 3 MB per image

The cache was emptied before each test and the time to complete the
thumbnail and preview generation was measured.

## MediaWEB executables 

| Name             | OS     | Architecture             | GOOS  | GOARCH | GOARM |
|------------------|--------|--------------------------|-------|--------|-------|
| mediaweb_x86-64  | Linux  | x86 64-bit               | linux | amd64  | N/A   |
| mediaweb_arm_v5  | Linux  | ARMv5 - softfloat 32-bit | linux | arm    | 5     |
| mediaweb_arm_v6  | Linux  | ARMv6 32-bit             | linux | arm    | 6     |
| mediaweb_arm_v7  | Linux  | ARMv7 32-bit             | linux | arm    | 7     |
| mediaweb_arm64   | Linux  | ARMv8 64-bit             | linux | arm64  | N/A   |

## Architectures

| Brand / model         | CPU Model              | Speed   | Cores | RAM  | OS                           |
|-----------------------|------------------------|---------|-------|------|------------------------------|
| Acer Chromebook CB714 | Intel Core i3 8130U    | 2.2 GHz | 2     | 4 GB | Chromium OS 10.0 64-bit      |
| ROCK64                | ARM Cortex A53 (ARMv8) | 1.5 GHz | 4     | 4 GB | Armbian Linux 4.4.184 64-bit |
| Banana Pi BPI-M1      | ARM Cortex A7 (ARMv7)  | 1 GHz   | 2     | 1 GB | Debian 8 32-bit              |

## Summary

The below table illustrates the result using the "best" 
MediaWEB executable for each target

| Measure                      | Acer (x86 64-bit) | ROCK65 (ARMv8) | BananaPi (ARMv7) |
|------------------------------|-------------------|----------------|------------------|'
| Total time                   | 27 sec            | 2 min 31 sec   | 5 min 26 sec     |
| Preview generation average   | 390 ms            | 2 110 ms       | 5 259 ms         |
| Thumbnail generation average | 226 ms            | 844 ms         | 1 854 ms         |   

## Logs from Acer Chromebook CB714

### mediaweb_x86-64

    joelmidstjarna@penguin:~/performance_test$ ./mediaweb_x86-64 
    2020/05/30 06:38:16 settings.go:56: INFO - Loading configuration: mediaweb.conf
    2020/05/30 06:38:16 main_common.go:15: INFO - Version: 
    2020/05/30 06:38:16 main_common.go:16: INFO - Build time: Sat May 30 06:30:44 CEST 2020
    2020/05/30 06:38:16 main_common.go:17: INFO - Git hash: 9a22d457e7cb90f7e0504a1a78d2a7212269a102
    2020/05/30 06:38:16 media.go:49: INFO - Media path: testmedia
    2020/05/30 06:38:16 media.go:59: INFO - Cache path: tmpcache
    2020/05/30 06:38:16 media.go:64: INFO - JPEG auto rotate: true
    2020/05/30 06:38:16 media.go:65: INFO - Image preview: true  (max width/height 1280 px)
    2020/05/30 06:38:16 media.go:74: INFO - Video thumbnails supported (ffmpeg installed): true
    2020/05/30 06:38:16 webapi.go:47: INFO - Starting Web API on port :9999
    2020/05/30 06:38:16 media.go:864: INFO - Pre-generating cache (thumbnails: true, preview: true)
    2020/05/30 06:38:16 watcher.go:50: INFO - Starting media watcher
    2020/05/30 06:38:16 updater.go:48: INFO - Starting updater
    2020/05/30 06:38:16 media.go:709: INFO - Creating new preview file for DSCN5369.JPG
    2020/05/30 06:38:16 media.go:718: INFO - Preview done for DSCN5369.JPG (conversion time: 436 ms)
    2020/05/30 06:38:17 media.go:709: INFO - Creating new preview file for DSCN5370.JPG
    2020/05/30 06:38:17 media.go:718: INFO - Preview done for DSCN5370.JPG (conversion time: 418 ms)
    2020/05/30 06:38:17 media.go:709: INFO - Creating new preview file for DSCN5371.JPG
    2020/05/30 06:38:18 media.go:718: INFO - Preview done for DSCN5371.JPG (conversion time: 364 ms)
    2020/05/30 06:38:18 media.go:709: INFO - Creating new preview file for DSCN5372.JPG
    2020/05/30 06:38:18 media.go:718: INFO - Preview done for DSCN5372.JPG (conversion time: 373 ms)
    2020/05/30 06:38:19 media.go:709: INFO - Creating new preview file for DSCN5373.JPG
    2020/05/30 06:38:19 media.go:718: INFO - Preview done for DSCN5373.JPG (conversion time: 383 ms)
    2020/05/30 06:38:19 media.go:709: INFO - Creating new preview file for DSCN5374.JPG
    2020/05/30 06:38:19 media.go:718: INFO - Preview done for DSCN5374.JPG (conversion time: 380 ms)
    2020/05/30 06:38:20 media.go:709: INFO - Creating new preview file for DSCN5375.JPG
    2020/05/30 06:38:20 media.go:718: INFO - Preview done for DSCN5375.JPG (conversion time: 385 ms)
    2020/05/30 06:38:20 media.go:709: INFO - Creating new preview file for DSCN5376.JPG
    2020/05/30 06:38:21 media.go:718: INFO - Preview done for DSCN5376.JPG (conversion time: 386 ms)
    2020/05/30 06:38:21 media.go:709: INFO - Creating new preview file for DSCN5377.JPG
    2020/05/30 06:38:21 media.go:718: INFO - Preview done for DSCN5377.JPG (conversion time: 459 ms)
    2020/05/30 06:38:22 media.go:709: INFO - Creating new preview file for DSCN5378.JPG
    2020/05/30 06:38:22 media.go:718: INFO - Preview done for DSCN5378.JPG (conversion time: 369 ms)
    2020/05/30 06:38:22 media.go:709: INFO - Creating new preview file for DSCN5379.JPG
    2020/05/30 06:38:23 media.go:718: INFO - Preview done for DSCN5379.JPG (conversion time: 363 ms)
    2020/05/30 06:38:23 media.go:709: INFO - Creating new preview file for DSCN5380.JPG
    2020/05/30 06:38:23 media.go:718: INFO - Preview done for DSCN5380.JPG (conversion time: 391 ms)
    2020/05/30 06:38:24 media.go:709: INFO - Creating new preview file for DSCN5381.JPG
    2020/05/30 06:38:24 media.go:718: INFO - Preview done for DSCN5381.JPG (conversion time: 404 ms)
    2020/05/30 06:38:24 media.go:709: INFO - Creating new preview file for DSCN5382.JPG
    2020/05/30 06:38:25 media.go:718: INFO - Preview done for DSCN5382.JPG (conversion time: 386 ms)
    2020/05/30 06:38:25 media.go:709: INFO - Creating new preview file for DSCN5383.JPG
    2020/05/30 06:38:25 media.go:718: INFO - Preview done for DSCN5383.JPG (conversion time: 371 ms)
    2020/05/30 06:38:25 media.go:709: INFO - Creating new preview file for DSCN5384.JPG
    2020/05/30 06:38:26 media.go:718: INFO - Preview done for DSCN5384.JPG (conversion time: 454 ms)
    2020/05/30 06:38:26 media.go:709: INFO - Creating new preview file for DSCN5385.JPG
    2020/05/30 06:38:26 media.go:718: INFO - Preview done for DSCN5385.JPG (conversion time: 410 ms)
    2020/05/30 06:38:27 media.go:709: INFO - Creating new preview file for DSCN5386.JPG
    2020/05/30 06:38:27 media.go:718: INFO - Preview done for DSCN5386.JPG (conversion time: 415 ms)
    2020/05/30 06:38:27 media.go:709: INFO - Creating new preview file for DSCN5387.JPG
    2020/05/30 06:38:28 media.go:718: INFO - Preview done for DSCN5387.JPG (conversion time: 400 ms)
    2020/05/30 06:38:28 media.go:709: INFO - Creating new preview file for DSCN5388.JPG
    2020/05/30 06:38:28 media.go:718: INFO - Preview done for DSCN5388.JPG (conversion time: 380 ms)
    2020/05/30 06:38:29 media.go:709: INFO - Creating new preview file for DSCN5389.JPG
    2020/05/30 06:38:29 media.go:718: INFO - Preview done for DSCN5389.JPG (conversion time: 364 ms)
    2020/05/30 06:38:29 media.go:709: INFO - Creating new preview file for DSCN5390.JPG
    2020/05/30 06:38:30 media.go:718: INFO - Preview done for DSCN5390.JPG (conversion time: 418 ms)
    2020/05/30 06:38:30 media.go:709: INFO - Creating new preview file for DSCN5391.JPG
    2020/05/30 06:38:30 media.go:718: INFO - Preview done for DSCN5391.JPG (conversion time: 392 ms)
    2020/05/30 06:38:31 media.go:709: INFO - Creating new preview file for DSCN5392.JPG
    2020/05/30 06:38:31 media.go:718: INFO - Preview done for DSCN5392.JPG (conversion time: 388 ms)
    2020/05/30 06:38:31 media.go:709: INFO - Creating new preview file for DSCN5393.JPG
    2020/05/30 06:38:32 media.go:718: INFO - Preview done for DSCN5393.JPG (conversion time: 371 ms)
    2020/05/30 06:38:32 media.go:709: INFO - Creating new preview file for DSCN5394.JPG
    2020/05/30 06:38:32 media.go:718: INFO - Preview done for DSCN5394.JPG (conversion time: 385 ms)
    2020/05/30 06:38:32 media.go:709: INFO - Creating new preview file for DSCN5395.JPG
    2020/05/30 06:38:33 media.go:718: INFO - Preview done for DSCN5395.JPG (conversion time: 362 ms)
    2020/05/30 06:38:33 media.go:709: INFO - Creating new preview file for DSCN5396.JPG
    2020/05/30 06:38:33 media.go:718: INFO - Preview done for DSCN5396.JPG (conversion time: 371 ms)
    2020/05/30 06:38:34 media.go:709: INFO - Creating new preview file for DSCN5397.JPG
    2020/05/30 06:38:34 media.go:718: INFO - Preview done for DSCN5397.JPG (conversion time: 360 ms)
    2020/05/30 06:38:34 media.go:709: INFO - Creating new preview file for DSCN5398.JPG
    2020/05/30 06:38:35 media.go:718: INFO - Preview done for DSCN5398.JPG (conversion time: 357 ms)
    2020/05/30 06:38:35 media.go:709: INFO - Creating new preview file for DSCN5399.JPG
    2020/05/30 06:38:35 media.go:718: INFO - Preview done for DSCN5399.JPG (conversion time: 386 ms)
    2020/05/30 06:38:35 media.go:709: INFO - Creating new preview file for exif_rotate/180deg.jpg
    2020/05/30 06:38:36 media.go:718: INFO - Preview done for exif_rotate/180deg.jpg (conversion time: 398 ms)
    2020/05/30 06:38:36 media.go:709: INFO - Creating new preview file for exif_rotate/mirror.jpg
    2020/05/30 06:38:36 media.go:718: INFO - Preview done for exif_rotate/mirror.jpg (conversion time: 377 ms)
    2020/05/30 06:38:37 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_270deg.jpg
    2020/05/30 06:38:37 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_270deg.jpg (conversion time: 420 ms)
    2020/05/30 06:38:37 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_90deg_cw.jpg
    2020/05/30 06:38:38 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_90deg_cw.jpg (conversion time: 456 ms)
    2020/05/30 06:38:38 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_vertical.jpg
    2020/05/30 06:38:38 media.go:718: INFO - Preview done for exif_rotate/mirror_vertical.jpg (conversion time: 388 ms)
    2020/05/30 06:38:38 media.go:442: INFO - Creating new thumbnail for exif_rotate/no_exif.jpg
    2020/05/30 06:38:38 media.go:460: INFO - Thumbnail done for exif_rotate/no_exif.jpg (conversion time: 30 ms)
    2020/05/30 06:38:39 media.go:709: INFO - Creating new preview file for exif_rotate/normal.jpg
    2020/05/30 06:38:39 media.go:718: INFO - Preview done for exif_rotate/normal.jpg (conversion time: 365 ms)
    2020/05/30 06:38:39 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_270deg_cw.jpg
    2020/05/30 06:38:40 media.go:718: INFO - Preview done for exif_rotate/rotate_270deg_cw.jpg (conversion time: 397 ms)
    2020/05/30 06:38:40 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_90deg_cw.jpg
    2020/05/30 06:38:40 media.go:718: INFO - Preview done for exif_rotate/rotate_90deg_cw.jpg (conversion time: 437 ms)
    2020/05/30 06:38:40 media.go:442: INFO - Creating new thumbnail for gif.gif
    2020/05/30 06:38:40 media.go:460: INFO - Thumbnail done for gif.gif (conversion time: 117 ms)
    2020/05/30 06:38:40 media.go:709: INFO - Creating new preview file for gif.gif
    2020/05/30 06:38:41 media.go:718: INFO - Preview done for gif.gif (conversion time: 232 ms)
    2020/05/30 06:38:41 media.go:709: INFO - Creating new preview file for jpeg.jpg
    2020/05/30 06:38:41 media.go:718: INFO - Preview done for jpeg.jpg (conversion time: 365 ms)
    2020/05/30 06:38:42 media.go:709: INFO - Creating new preview file for jpeg_rotated.jpg
    2020/05/30 06:38:42 media.go:718: INFO - Preview done for jpeg_rotated.jpg (conversion time: 408 ms)
    2020/05/30 06:38:42 media.go:442: INFO - Creating new thumbnail for png.png
    2020/05/30 06:38:42 media.go:460: INFO - Thumbnail done for png.png (conversion time: 155 ms)
    2020/05/30 06:38:42 media.go:709: INFO - Creating new preview file for png.png
    2020/05/30 06:38:42 media.go:718: INFO - Preview done for png.png (conversion time: 226 ms)
    2020/05/30 06:38:42 media.go:442: INFO - Creating new thumbnail for screenshot_browser.jpg
    2020/05/30 06:38:42 media.go:460: INFO - Thumbnail done for screenshot_browser.jpg (conversion time: 12 ms)
    2020/05/30 06:38:42 media.go:442: INFO - Creating new thumbnail for screenshot_mobile.jpg
    2020/05/30 06:38:42 media.go:460: INFO - Thumbnail done for screenshot_mobile.jpg (conversion time: 9 ms)
    2020/05/30 06:38:42 media.go:442: INFO - Creating new thumbnail for screenshot_viewer.jpg
    2020/05/30 06:38:43 media.go:460: INFO - Thumbnail done for screenshot_viewer.jpg (conversion time: 9 ms)
    2020/05/30 06:38:43 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_swipe.jpg
    2020/05/30 06:38:43 media.go:460: INFO - Thumbnail done for screenshot_viewer_swipe.jpg (conversion time: 9 ms)
    2020/05/30 06:38:43 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_zoom.jpg
    2020/05/30 06:38:43 media.go:460: INFO - Thumbnail done for screenshot_viewer_zoom.jpg (conversion time: 9 ms)
    2020/05/30 06:38:43 media.go:442: INFO - Creating new thumbnail for tiff.tiff
    2020/05/30 06:38:43 media.go:460: INFO - Thumbnail done for tiff.tiff (conversion time: 63 ms)
    2020/05/30 06:38:43 media.go:870: INFO - Generating cache took 0 minutes and 27 seconds
    Number of folders: 1
    Number of images: 50
    Number of videos: 0
    Number of images with embedded EXIF: 41
    Number of generated image thumbnails: 9
    Number of generated video thumbnails: 0
    Number of generated image previews: 43
    Number of failed folders: 0
    Number of failed image thumbnails: 0
    Number of failed video thumbnails: 0
    Number of failed image previews: 0
    Number of small images not require preview: 7
  
  
## Logs from ROCK64
  
### mediaweb_arm_v5

    pi@rock64:~/performance_test$ ./mediaweb_arm_v5
    2020/05/30 07:03:17 settings.go:56: INFO - Loading configuration: mediaweb.conf
    2020/05/30 07:03:17 main_common.go:15: INFO - Version: 
    2020/05/30 07:03:17 main_common.go:16: INFO - Build time: Sat May 30 06:44:04 CEST 2020
    2020/05/30 07:03:17 main_common.go:17: INFO - Git hash: 9a22d457e7cb90f7e0504a1a78d2a7212269a102
    2020/05/30 07:03:17 media.go:49: INFO - Media path: testmedia
    2020/05/30 07:03:17 media.go:59: INFO - Cache path: tmpcache
    2020/05/30 07:03:17 media.go:64: INFO - JPEG auto rotate: true
    2020/05/30 07:03:17 media.go:65: INFO - Image preview: true  (max width/height 1280 px)
    2020/05/30 07:03:17 media.go:74: INFO - Video thumbnails supported (ffmpeg installed): true
    2020/05/30 07:03:17 webapi.go:47: INFO - Starting Web API on port :9999
    2020/05/30 07:03:17 media.go:864: INFO - Pre-generating cache (thumbnails: true, preview: true)
    2020/05/30 07:03:17 watcher.go:50: INFO - Starting media watcher
    2020/05/30 07:03:17 updater.go:48: INFO - Starting updater
    2020/05/30 07:03:19 media.go:709: INFO - Creating new preview file for DSCN5369.JPG
    2020/05/30 07:03:36 media.go:718: INFO - Preview done for DSCN5369.JPG (conversion time: 16825 ms)
    2020/05/30 07:03:37 media.go:709: INFO - Creating new preview file for DSCN5370.JPG
    2020/05/30 07:03:54 media.go:718: INFO - Preview done for DSCN5370.JPG (conversion time: 16841 ms)
    2020/05/30 07:03:55 media.go:709: INFO - Creating new preview file for DSCN5371.JPG
    2020/05/30 07:04:12 media.go:718: INFO - Preview done for DSCN5371.JPG (conversion time: 16530 ms)
    2020/05/30 07:04:13 media.go:709: INFO - Creating new preview file for DSCN5372.JPG
    2020/05/30 07:04:30 media.go:718: INFO - Preview done for DSCN5372.JPG (conversion time: 16570 ms)
    2020/05/30 07:04:32 media.go:709: INFO - Creating new preview file for DSCN5373.JPG
    2020/05/30 07:04:48 media.go:718: INFO - Preview done for DSCN5373.JPG (conversion time: 16672 ms)
    2020/05/30 07:04:50 media.go:709: INFO - Creating new preview file for DSCN5374.JPG
    2020/05/30 07:05:06 media.go:718: INFO - Preview done for DSCN5374.JPG (conversion time: 16432 ms)
    2020/05/30 07:05:07 media.go:709: INFO - Creating new preview file for DSCN5375.JPG
    2020/05/30 07:05:24 media.go:718: INFO - Preview done for DSCN5375.JPG (conversion time: 16697 ms)
    2020/05/30 07:05:26 media.go:709: INFO - Creating new preview file for DSCN5376.JPG
    2020/05/30 07:05:42 media.go:718: INFO - Preview done for DSCN5376.JPG (conversion time: 16607 ms)
    2020/05/30 07:05:44 media.go:709: INFO - Creating new preview file for DSCN5377.JPG
    2020/05/30 07:06:01 media.go:718: INFO - Preview done for DSCN5377.JPG (conversion time: 16900 ms)
    2020/05/30 07:06:02 media.go:709: INFO - Creating new preview file for DSCN5378.JPG
    2020/05/30 07:06:19 media.go:718: INFO - Preview done for DSCN5378.JPG (conversion time: 16630 ms)
    2020/05/30 07:06:21 media.go:709: INFO - Creating new preview file for DSCN5379.JPG
    2020/05/30 07:06:37 media.go:718: INFO - Preview done for DSCN5379.JPG (conversion time: 16112 ms)
    2020/05/30 07:06:38 media.go:709: INFO - Creating new preview file for DSCN5380.JPG
    2020/05/30 07:06:55 media.go:718: INFO - Preview done for DSCN5380.JPG (conversion time: 16676 ms)
    2020/05/30 07:06:56 media.go:709: INFO - Creating new preview file for DSCN5381.JPG
    2020/05/30 07:07:13 media.go:718: INFO - Preview done for DSCN5381.JPG (conversion time: 16997 ms)
    2020/05/30 07:07:15 media.go:709: INFO - Creating new preview file for DSCN5382.JPG
    2020/05/30 07:07:30 media.go:718: INFO - Preview done for DSCN5382.JPG (conversion time: 15554 ms)
    2020/05/30 07:07:32 media.go:709: INFO - Creating new preview file for DSCN5383.JPG
    2020/05/30 07:07:48 media.go:718: INFO - Preview done for DSCN5383.JPG (conversion time: 15849 ms)
    2020/05/30 07:07:49 media.go:709: INFO - Creating new preview file for DSCN5384.JPG
    2020/05/30 07:08:06 media.go:718: INFO - Preview done for DSCN5384.JPG (conversion time: 16893 ms)
    2020/05/30 07:08:08 media.go:709: INFO - Creating new preview file for DSCN5385.JPG
    2020/05/30 07:08:25 media.go:718: INFO - Preview done for DSCN5385.JPG (conversion time: 16895 ms)
    2020/05/30 07:08:26 media.go:709: INFO - Creating new preview file for DSCN5386.JPG
    2020/05/30 07:08:43 media.go:718: INFO - Preview done for DSCN5386.JPG (conversion time: 16726 ms)
    2020/05/30 07:08:45 media.go:709: INFO - Creating new preview file for DSCN5387.JPG
    2020/05/30 07:09:01 media.go:718: INFO - Preview done for DSCN5387.JPG (conversion time: 16884 ms)
    2020/05/30 07:09:03 media.go:709: INFO - Creating new preview file for DSCN5388.JPG
    2020/05/30 07:09:20 media.go:718: INFO - Preview done for DSCN5388.JPG (conversion time: 16777 ms)
    2020/05/30 07:09:21 media.go:709: INFO - Creating new preview file for DSCN5389.JPG
    2020/05/30 07:09:38 media.go:718: INFO - Preview done for DSCN5389.JPG (conversion time: 16687 ms)
    2020/05/30 07:09:39 media.go:709: INFO - Creating new preview file for DSCN5390.JPG
    2020/05/30 07:09:56 media.go:718: INFO - Preview done for DSCN5390.JPG (conversion time: 16641 ms)
    2020/05/30 07:09:58 media.go:709: INFO - Creating new preview file for DSCN5391.JPG
    2020/05/30 07:10:14 media.go:718: INFO - Preview done for DSCN5391.JPG (conversion time: 16750 ms)
    2020/05/30 07:10:16 media.go:709: INFO - Creating new preview file for DSCN5392.JPG
    2020/05/30 07:10:32 media.go:718: INFO - Preview done for DSCN5392.JPG (conversion time: 16178 ms)
    2020/05/30 07:10:33 media.go:709: INFO - Creating new preview file for DSCN5393.JPG
    2020/05/30 07:10:50 media.go:718: INFO - Preview done for DSCN5393.JPG (conversion time: 16623 ms)
    2020/05/30 07:10:51 media.go:709: INFO - Creating new preview file for DSCN5394.JPG
    2020/05/30 07:11:08 media.go:718: INFO - Preview done for DSCN5394.JPG (conversion time: 16753 ms)
    2020/05/30 07:11:10 media.go:709: INFO - Creating new preview file for DSCN5395.JPG
    2020/05/30 07:11:26 media.go:718: INFO - Preview done for DSCN5395.JPG (conversion time: 16705 ms)
    2020/05/30 07:11:28 media.go:709: INFO - Creating new preview file for DSCN5396.JPG
    2020/05/30 07:11:44 media.go:718: INFO - Preview done for DSCN5396.JPG (conversion time: 16576 ms)
    2020/05/30 07:11:46 media.go:709: INFO - Creating new preview file for DSCN5397.JPG
    2020/05/30 07:12:02 media.go:718: INFO - Preview done for DSCN5397.JPG (conversion time: 16592 ms)
    2020/05/30 07:12:04 media.go:709: INFO - Creating new preview file for DSCN5398.JPG
    2020/05/30 07:12:21 media.go:718: INFO - Preview done for DSCN5398.JPG (conversion time: 16599 ms)
    2020/05/30 07:12:22 media.go:709: INFO - Creating new preview file for DSCN5399.JPG
    2020/05/30 07:12:39 media.go:718: INFO - Preview done for DSCN5399.JPG (conversion time: 16596 ms)
    2020/05/30 07:12:40 media.go:709: INFO - Creating new preview file for exif_rotate/180deg.jpg
    2020/05/30 07:12:57 media.go:718: INFO - Preview done for exif_rotate/180deg.jpg (conversion time: 17409 ms)
    2020/05/30 07:12:59 media.go:709: INFO - Creating new preview file for exif_rotate/mirror.jpg
    2020/05/30 07:13:16 media.go:718: INFO - Preview done for exif_rotate/mirror.jpg (conversion time: 17349 ms)
    2020/05/30 07:13:18 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_270deg.jpg
    2020/05/30 07:13:36 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_270deg.jpg (conversion time: 18046 ms)
    2020/05/30 07:13:37 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_90deg_cw.jpg
    2020/05/30 07:13:55 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_90deg_cw.jpg (conversion time: 18182 ms)
    2020/05/30 07:13:57 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_vertical.jpg
    2020/05/30 07:14:14 media.go:718: INFO - Preview done for exif_rotate/mirror_vertical.jpg (conversion time: 17426 ms)
    2020/05/30 07:14:14 media.go:442: INFO - Creating new thumbnail for exif_rotate/no_exif.jpg
    2020/05/30 07:14:15 media.go:460: INFO - Thumbnail done for exif_rotate/no_exif.jpg (conversion time: 1028 ms)
    2020/05/30 07:14:17 media.go:709: INFO - Creating new preview file for exif_rotate/normal.jpg
    2020/05/30 07:14:34 media.go:718: INFO - Preview done for exif_rotate/normal.jpg (conversion time: 17472 ms)
    2020/05/30 07:14:36 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_270deg_cw.jpg
    2020/05/30 07:14:54 media.go:718: INFO - Preview done for exif_rotate/rotate_270deg_cw.jpg (conversion time: 18048 ms)
    2020/05/30 07:14:55 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_90deg_cw.jpg
    2020/05/30 07:15:13 media.go:718: INFO - Preview done for exif_rotate/rotate_90deg_cw.jpg (conversion time: 18163 ms)
    2020/05/30 07:15:13 media.go:442: INFO - Creating new thumbnail for gif.gif
    2020/05/30 07:15:21 media.go:460: INFO - Thumbnail done for gif.gif (conversion time: 7442 ms)
    2020/05/30 07:15:21 media.go:709: INFO - Creating new preview file for gif.gif
    2020/05/30 07:15:37 media.go:718: INFO - Preview done for gif.gif (conversion time: 15664 ms)
    2020/05/30 07:15:38 media.go:709: INFO - Creating new preview file for jpeg.jpg
    2020/05/30 07:15:56 media.go:718: INFO - Preview done for jpeg.jpg (conversion time: 17344 ms)
    2020/05/30 07:15:57 media.go:709: INFO - Creating new preview file for jpeg_rotated.jpg
    2020/05/30 07:16:15 media.go:718: INFO - Preview done for jpeg_rotated.jpg (conversion time: 18148 ms)
    2020/05/30 07:16:15 media.go:442: INFO - Creating new thumbnail for png.png
    2020/05/30 07:16:19 media.go:460: INFO - Thumbnail done for png.png (conversion time: 3071 ms)
    2020/05/30 07:16:19 media.go:709: INFO - Creating new preview file for png.png
    2020/05/30 07:16:27 media.go:718: INFO - Preview done for png.png (conversion time: 7544 ms)
    2020/05/30 07:16:27 media.go:442: INFO - Creating new thumbnail for screenshot_browser.jpg
    2020/05/30 07:16:28 media.go:460: INFO - Thumbnail done for screenshot_browser.jpg (conversion time: 577 ms)
    2020/05/30 07:16:28 media.go:442: INFO - Creating new thumbnail for screenshot_mobile.jpg
    2020/05/30 07:16:28 media.go:460: INFO - Thumbnail done for screenshot_mobile.jpg (conversion time: 312 ms)
    2020/05/30 07:16:28 media.go:442: INFO - Creating new thumbnail for screenshot_viewer.jpg
    2020/05/30 07:16:28 media.go:460: INFO - Thumbnail done for screenshot_viewer.jpg (conversion time: 396 ms)
    2020/05/30 07:16:28 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_swipe.jpg
    2020/05/30 07:16:29 media.go:460: INFO - Thumbnail done for screenshot_viewer_swipe.jpg (conversion time: 408 ms)
    2020/05/30 07:16:29 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_zoom.jpg
    2020/05/30 07:16:29 media.go:460: INFO - Thumbnail done for screenshot_viewer_zoom.jpg (conversion time: 351 ms)
    2020/05/30 07:16:29 media.go:442: INFO - Creating new thumbnail for tiff.tiff
    2020/05/30 07:16:31 media.go:460: INFO - Thumbnail done for tiff.tiff (conversion time: 1304 ms)
    2020/05/30 07:16:31 media.go:870: INFO - Generating cache took 13 minutes and 13 seconds
    Number of folders: 1
    Number of images: 50
    Number of videos: 0
    Number of images with embedded EXIF: 41
    Number of generated image thumbnails: 9
    Number of generated video thumbnails: 0
    Number of generated image previews: 43
    Number of failed folders: 0
    Number of failed image thumbnails: 0
    Number of failed video thumbnails: 0
    Number of failed image previews: 0
    Number of small images not require preview: 7


### mediaweb_arm_v6

    pi@rock64:~/performance_test$ ./mediaweb_arm_v6
    2020/05/30 07:21:27 settings.go:56: INFO - Loading configuration: mediaweb.conf
    2020/05/30 07:21:27 main_common.go:15: INFO - Version: 
    2020/05/30 07:21:27 main_common.go:16: INFO - Build time: Sat May 30 06:48:05 CEST 2020
    2020/05/30 07:21:27 main_common.go:17: INFO - Git hash: 9a22d457e7cb90f7e0504a1a78d2a7212269a102
    2020/05/30 07:21:27 media.go:49: INFO - Media path: testmedia
    2020/05/30 07:21:27 media.go:59: INFO - Cache path: tmpcache
    2020/05/30 07:21:27 media.go:64: INFO - JPEG auto rotate: true
    2020/05/30 07:21:27 media.go:65: INFO - Image preview: true  (max width/height 1280 px)
    2020/05/30 07:21:27 media.go:74: INFO - Video thumbnails supported (ffmpeg installed): true
    2020/05/30 07:21:27 webapi.go:47: INFO - Starting Web API on port :9999
    2020/05/30 07:21:27 watcher.go:50: INFO - Starting media watcher
    2020/05/30 07:21:27 updater.go:48: INFO - Starting updater
    2020/05/30 07:21:27 media.go:864: INFO - Pre-generating cache (thumbnails: true, preview: true)
    2020/05/30 07:21:28 media.go:709: INFO - Creating new preview file for DSCN5369.JPG
    2020/05/30 07:21:31 media.go:718: INFO - Preview done for DSCN5369.JPG (conversion time: 2869 ms)
    2020/05/30 07:21:33 media.go:709: INFO - Creating new preview file for DSCN5370.JPG
    2020/05/30 07:21:36 media.go:718: INFO - Preview done for DSCN5370.JPG (conversion time: 2875 ms)
    2020/05/30 07:21:37 media.go:709: INFO - Creating new preview file for DSCN5371.JPG
    2020/05/30 07:21:40 media.go:718: INFO - Preview done for DSCN5371.JPG (conversion time: 2622 ms)
    2020/05/30 07:21:41 media.go:709: INFO - Creating new preview file for DSCN5372.JPG
    2020/05/30 07:21:44 media.go:718: INFO - Preview done for DSCN5372.JPG (conversion time: 2627 ms)
    2020/05/30 07:21:45 media.go:709: INFO - Creating new preview file for DSCN5373.JPG
    2020/05/30 07:21:48 media.go:718: INFO - Preview done for DSCN5373.JPG (conversion time: 2730 ms)
    2020/05/30 07:21:50 media.go:709: INFO - Creating new preview file for DSCN5374.JPG
    2020/05/30 07:21:52 media.go:718: INFO - Preview done for DSCN5374.JPG (conversion time: 2614 ms)
    2020/05/30 07:21:54 media.go:709: INFO - Creating new preview file for DSCN5375.JPG
    2020/05/30 07:21:56 media.go:718: INFO - Preview done for DSCN5375.JPG (conversion time: 2727 ms)
    2020/05/30 07:21:58 media.go:709: INFO - Creating new preview file for DSCN5376.JPG
    2020/05/30 07:22:01 media.go:718: INFO - Preview done for DSCN5376.JPG (conversion time: 2708 ms)
    2020/05/30 07:22:02 media.go:709: INFO - Creating new preview file for DSCN5377.JPG
    2020/05/30 07:22:05 media.go:718: INFO - Preview done for DSCN5377.JPG (conversion time: 3003 ms)
    2020/05/30 07:22:07 media.go:709: INFO - Creating new preview file for DSCN5378.JPG
    2020/05/30 07:22:10 media.go:718: INFO - Preview done for DSCN5378.JPG (conversion time: 2699 ms)
    2020/05/30 07:22:11 media.go:709: INFO - Creating new preview file for DSCN5379.JPG
    2020/05/30 07:22:14 media.go:718: INFO - Preview done for DSCN5379.JPG (conversion time: 2674 ms)
    2020/05/30 07:22:15 media.go:709: INFO - Creating new preview file for DSCN5380.JPG
    2020/05/30 07:22:18 media.go:718: INFO - Preview done for DSCN5380.JPG (conversion time: 2671 ms)
    2020/05/30 07:22:19 media.go:709: INFO - Creating new preview file for DSCN5381.JPG
    2020/05/30 07:22:22 media.go:718: INFO - Preview done for DSCN5381.JPG (conversion time: 2842 ms)
    2020/05/30 07:22:24 media.go:709: INFO - Creating new preview file for DSCN5382.JPG
    2020/05/30 07:22:26 media.go:718: INFO - Preview done for DSCN5382.JPG (conversion time: 2629 ms)
    2020/05/30 07:22:28 media.go:709: INFO - Creating new preview file for DSCN5383.JPG
    2020/05/30 07:22:30 media.go:718: INFO - Preview done for DSCN5383.JPG (conversion time: 2646 ms)
    2020/05/30 07:22:32 media.go:709: INFO - Creating new preview file for DSCN5384.JPG
    2020/05/30 07:22:35 media.go:718: INFO - Preview done for DSCN5384.JPG (conversion time: 2882 ms)
    2020/05/30 07:22:36 media.go:709: INFO - Creating new preview file for DSCN5385.JPG
    2020/05/30 07:22:39 media.go:718: INFO - Preview done for DSCN5385.JPG (conversion time: 2802 ms)
    2020/05/30 07:22:41 media.go:709: INFO - Creating new preview file for DSCN5386.JPG
    2020/05/30 07:22:44 media.go:718: INFO - Preview done for DSCN5386.JPG (conversion time: 2792 ms)
    2020/05/30 07:22:45 media.go:709: INFO - Creating new preview file for DSCN5387.JPG
    2020/05/30 07:22:48 media.go:718: INFO - Preview done for DSCN5387.JPG (conversion time: 2808 ms)
    2020/05/30 07:22:50 media.go:709: INFO - Creating new preview file for DSCN5388.JPG
    2020/05/30 07:22:52 media.go:718: INFO - Preview done for DSCN5388.JPG (conversion time: 2799 ms)
    2020/05/30 07:22:54 media.go:709: INFO - Creating new preview file for DSCN5389.JPG
    2020/05/30 07:22:56 media.go:718: INFO - Preview done for DSCN5389.JPG (conversion time: 2683 ms)
    2020/05/30 07:22:58 media.go:709: INFO - Creating new preview file for DSCN5390.JPG
    2020/05/30 07:23:01 media.go:718: INFO - Preview done for DSCN5390.JPG (conversion time: 2708 ms)
    2020/05/30 07:23:02 media.go:709: INFO - Creating new preview file for DSCN5391.JPG
    2020/05/30 07:23:05 media.go:718: INFO - Preview done for DSCN5391.JPG (conversion time: 2791 ms)
    2020/05/30 07:23:06 media.go:709: INFO - Creating new preview file for DSCN5392.JPG
    2020/05/30 07:23:09 media.go:718: INFO - Preview done for DSCN5392.JPG (conversion time: 2694 ms)
    2020/05/30 07:23:11 media.go:709: INFO - Creating new preview file for DSCN5393.JPG
    2020/05/30 07:23:13 media.go:718: INFO - Preview done for DSCN5393.JPG (conversion time: 2643 ms)
    2020/05/30 07:23:15 media.go:709: INFO - Creating new preview file for DSCN5394.JPG
    2020/05/30 07:23:17 media.go:718: INFO - Preview done for DSCN5394.JPG (conversion time: 2702 ms)
    2020/05/30 07:23:19 media.go:709: INFO - Creating new preview file for DSCN5395.JPG
    2020/05/30 07:23:21 media.go:718: INFO - Preview done for DSCN5395.JPG (conversion time: 2678 ms)
    2020/05/30 07:23:23 media.go:709: INFO - Creating new preview file for DSCN5396.JPG
    2020/05/30 07:23:26 media.go:718: INFO - Preview done for DSCN5396.JPG (conversion time: 2659 ms)
    2020/05/30 07:23:27 media.go:709: INFO - Creating new preview file for DSCN5397.JPG
    2020/05/30 07:23:30 media.go:718: INFO - Preview done for DSCN5397.JPG (conversion time: 2647 ms)
    2020/05/30 07:23:31 media.go:709: INFO - Creating new preview file for DSCN5398.JPG
    2020/05/30 07:23:34 media.go:718: INFO - Preview done for DSCN5398.JPG (conversion time: 2663 ms)
    2020/05/30 07:23:35 media.go:709: INFO - Creating new preview file for DSCN5399.JPG
    2020/05/30 07:23:38 media.go:718: INFO - Preview done for DSCN5399.JPG (conversion time: 2675 ms)
    2020/05/30 07:23:39 media.go:709: INFO - Creating new preview file for exif_rotate/180deg.jpg
    2020/05/30 07:23:42 media.go:718: INFO - Preview done for exif_rotate/180deg.jpg (conversion time: 2590 ms)
    2020/05/30 07:23:43 media.go:709: INFO - Creating new preview file for exif_rotate/mirror.jpg
    2020/05/30 07:23:46 media.go:718: INFO - Preview done for exif_rotate/mirror.jpg (conversion time: 2603 ms)
    2020/05/30 07:23:47 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_270deg.jpg
    2020/05/30 07:23:51 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_270deg.jpg (conversion time: 3231 ms)
    2020/05/30 07:23:52 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_90deg_cw.jpg
    2020/05/30 07:23:55 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_90deg_cw.jpg (conversion time: 3261 ms)
    2020/05/30 07:23:57 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_vertical.jpg
    2020/05/30 07:23:59 media.go:718: INFO - Preview done for exif_rotate/mirror_vertical.jpg (conversion time: 2576 ms)
    2020/05/30 07:23:59 media.go:442: INFO - Creating new thumbnail for exif_rotate/no_exif.jpg
    2020/05/30 07:23:59 media.go:460: INFO - Thumbnail done for exif_rotate/no_exif.jpg (conversion time: 188 ms)
    2020/05/30 07:24:01 media.go:709: INFO - Creating new preview file for exif_rotate/normal.jpg
    2020/05/30 07:24:04 media.go:718: INFO - Preview done for exif_rotate/normal.jpg (conversion time: 2558 ms)
    2020/05/30 07:24:05 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_270deg_cw.jpg
    2020/05/30 07:24:08 media.go:718: INFO - Preview done for exif_rotate/rotate_270deg_cw.jpg (conversion time: 3268 ms)
    2020/05/30 07:24:10 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_90deg_cw.jpg
    2020/05/30 07:24:13 media.go:718: INFO - Preview done for exif_rotate/rotate_90deg_cw.jpg (conversion time: 3321 ms)
    2020/05/30 07:24:13 media.go:442: INFO - Creating new thumbnail for gif.gif
    2020/05/30 07:24:14 media.go:460: INFO - Thumbnail done for gif.gif (conversion time: 855 ms)
    2020/05/30 07:24:14 media.go:709: INFO - Creating new preview file for gif.gif
    2020/05/30 07:24:16 media.go:718: INFO - Preview done for gif.gif (conversion time: 1792 ms)
    2020/05/30 07:24:18 media.go:709: INFO - Creating new preview file for jpeg.jpg
    2020/05/30 07:24:20 media.go:718: INFO - Preview done for jpeg.jpg (conversion time: 2539 ms)
    2020/05/30 07:24:22 media.go:709: INFO - Creating new preview file for jpeg_rotated.jpg
    2020/05/30 07:24:25 media.go:718: INFO - Preview done for jpeg_rotated.jpg (conversion time: 3333 ms)
    2020/05/30 07:24:25 media.go:442: INFO - Creating new thumbnail for png.png
    2020/05/30 07:24:26 media.go:460: INFO - Thumbnail done for png.png (conversion time: 1079 ms)
    2020/05/30 07:24:27 media.go:709: INFO - Creating new preview file for png.png
    2020/05/30 07:24:29 media.go:718: INFO - Preview done for png.png (conversion time: 1820 ms)
    2020/05/30 07:24:29 media.go:442: INFO - Creating new thumbnail for screenshot_browser.jpg
    2020/05/30 07:24:29 media.go:460: INFO - Thumbnail done for screenshot_browser.jpg (conversion time: 97 ms)
    2020/05/30 07:24:29 media.go:442: INFO - Creating new thumbnail for screenshot_mobile.jpg
    2020/05/30 07:24:29 media.go:460: INFO - Thumbnail done for screenshot_mobile.jpg (conversion time: 61 ms)
    2020/05/30 07:24:29 media.go:442: INFO - Creating new thumbnail for screenshot_viewer.jpg
    2020/05/30 07:24:29 media.go:460: INFO - Thumbnail done for screenshot_viewer.jpg (conversion time: 73 ms)
    2020/05/30 07:24:29 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_swipe.jpg
    2020/05/30 07:24:29 media.go:460: INFO - Thumbnail done for screenshot_viewer_swipe.jpg (conversion time: 71 ms)
    2020/05/30 07:24:29 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_zoom.jpg
    2020/05/30 07:24:29 media.go:460: INFO - Thumbnail done for screenshot_viewer_zoom.jpg (conversion time: 72 ms)
    2020/05/30 07:24:29 media.go:442: INFO - Creating new thumbnail for tiff.tiff
    2020/05/30 07:24:30 media.go:460: INFO - Thumbnail done for tiff.tiff (conversion time: 368 ms)
    2020/05/30 07:24:30 media.go:870: INFO - Generating cache took 3 minutes and 3 seconds
    Number of folders: 1
    Number of images: 50
    Number of videos: 0
    Number of images with embedded EXIF: 41
    Number of generated image thumbnails: 9
    Number of generated video thumbnails: 0
    Number of generated image previews: 43
    Number of failed folders: 0
    Number of failed image thumbnails: 0

### mediaweb_arm_v7

    pi@rock64:~/performance_test$ ./mediaweb_arm_v7
    2020/05/30 07:30:27 settings.go:56: INFO - Loading configuration: mediaweb.conf
    2020/05/30 07:30:27 main_common.go:15: INFO - Version: 
    2020/05/30 07:30:27 main_common.go:16: INFO - Build time: Sat May 30 06:49:04 CEST 2020
    2020/05/30 07:30:27 main_common.go:17: INFO - Git hash: 9a22d457e7cb90f7e0504a1a78d2a7212269a102
    2020/05/30 07:30:27 media.go:49: INFO - Media path: testmedia
    2020/05/30 07:30:27 media.go:59: INFO - Cache path: tmpcache
    2020/05/30 07:30:27 media.go:64: INFO - JPEG auto rotate: true
    2020/05/30 07:30:27 media.go:65: INFO - Image preview: true  (max width/height 1280 px)
    2020/05/30 07:30:27 media.go:74: INFO - Video thumbnails supported (ffmpeg installed): true
    2020/05/30 07:30:27 webapi.go:47: INFO - Starting Web API on port :9999
    2020/05/30 07:30:27 watcher.go:50: INFO - Starting media watcher
    2020/05/30 07:30:27 updater.go:48: INFO - Starting updater
    2020/05/30 07:30:27 media.go:864: INFO - Pre-generating cache (thumbnails: true, preview: true)
    2020/05/30 07:30:29 media.go:709: INFO - Creating new preview file for DSCN5369.JPG
    2020/05/30 07:30:31 media.go:718: INFO - Preview done for DSCN5369.JPG (conversion time: 2756 ms)
    2020/05/30 07:30:33 media.go:709: INFO - Creating new preview file for DSCN5370.JPG
    2020/05/30 07:30:36 media.go:718: INFO - Preview done for DSCN5370.JPG (conversion time: 2749 ms)
    2020/05/30 07:30:37 media.go:709: INFO - Creating new preview file for DSCN5371.JPG
    2020/05/30 07:30:40 media.go:718: INFO - Preview done for DSCN5371.JPG (conversion time: 2519 ms)
    2020/05/30 07:30:41 media.go:709: INFO - Creating new preview file for DSCN5372.JPG
    2020/05/30 07:30:44 media.go:718: INFO - Preview done for DSCN5372.JPG (conversion time: 2545 ms)
    2020/05/30 07:30:45 media.go:709: INFO - Creating new preview file for DSCN5373.JPG
    2020/05/30 07:30:48 media.go:718: INFO - Preview done for DSCN5373.JPG (conversion time: 2596 ms)
    2020/05/30 07:30:49 media.go:709: INFO - Creating new preview file for DSCN5374.JPG
    2020/05/30 07:30:51 media.go:718: INFO - Preview done for DSCN5374.JPG (conversion time: 2473 ms)
    2020/05/30 07:30:53 media.go:709: INFO - Creating new preview file for DSCN5375.JPG
    2020/05/30 07:30:55 media.go:718: INFO - Preview done for DSCN5375.JPG (conversion time: 2594 ms)
    2020/05/30 07:30:57 media.go:709: INFO - Creating new preview file for DSCN5376.JPG
    2020/05/30 07:31:00 media.go:718: INFO - Preview done for DSCN5376.JPG (conversion time: 2596 ms)
    2020/05/30 07:31:01 media.go:709: INFO - Creating new preview file for DSCN5377.JPG
    2020/05/30 07:31:04 media.go:718: INFO - Preview done for DSCN5377.JPG (conversion time: 2860 ms)
    2020/05/30 07:31:06 media.go:709: INFO - Creating new preview file for DSCN5378.JPG
    2020/05/30 07:31:08 media.go:718: INFO - Preview done for DSCN5378.JPG (conversion time: 2570 ms)
    2020/05/30 07:31:09 media.go:709: INFO - Creating new preview file for DSCN5379.JPG
    2020/05/30 07:31:12 media.go:718: INFO - Preview done for DSCN5379.JPG (conversion time: 2607 ms)
    2020/05/30 07:31:14 media.go:709: INFO - Creating new preview file for DSCN5380.JPG
    2020/05/30 07:31:16 media.go:718: INFO - Preview done for DSCN5380.JPG (conversion time: 2580 ms)
    2020/05/30 07:31:18 media.go:709: INFO - Creating new preview file for DSCN5381.JPG
    2020/05/30 07:31:20 media.go:718: INFO - Preview done for DSCN5381.JPG (conversion time: 2696 ms)
    2020/05/30 07:31:22 media.go:709: INFO - Creating new preview file for DSCN5382.JPG
    2020/05/30 07:31:24 media.go:718: INFO - Preview done for DSCN5382.JPG (conversion time: 2532 ms)
    2020/05/30 07:31:26 media.go:709: INFO - Creating new preview file for DSCN5383.JPG
    2020/05/30 07:31:28 media.go:718: INFO - Preview done for DSCN5383.JPG (conversion time: 2550 ms)
    2020/05/30 07:31:30 media.go:709: INFO - Creating new preview file for DSCN5384.JPG
    2020/05/30 07:31:33 media.go:718: INFO - Preview done for DSCN5384.JPG (conversion time: 2789 ms)
    2020/05/30 07:31:34 media.go:709: INFO - Creating new preview file for DSCN5385.JPG
    2020/05/30 07:31:37 media.go:718: INFO - Preview done for DSCN5385.JPG (conversion time: 2698 ms)
    2020/05/30 07:31:38 media.go:709: INFO - Creating new preview file for DSCN5386.JPG
    2020/05/30 07:31:41 media.go:718: INFO - Preview done for DSCN5386.JPG (conversion time: 2735 ms)
    2020/05/30 07:31:43 media.go:709: INFO - Creating new preview file for DSCN5387.JPG
    2020/05/30 07:31:45 media.go:718: INFO - Preview done for DSCN5387.JPG (conversion time: 2700 ms)
    2020/05/30 07:31:47 media.go:709: INFO - Creating new preview file for DSCN5388.JPG
    2020/05/30 07:31:49 media.go:718: INFO - Preview done for DSCN5388.JPG (conversion time: 2678 ms)
    2020/05/30 07:31:51 media.go:709: INFO - Creating new preview file for DSCN5389.JPG
    2020/05/30 07:31:53 media.go:718: INFO - Preview done for DSCN5389.JPG (conversion time: 2549 ms)
    2020/05/30 07:31:55 media.go:709: INFO - Creating new preview file for DSCN5390.JPG
    2020/05/30 07:31:57 media.go:718: INFO - Preview done for DSCN5390.JPG (conversion time: 2566 ms)
    2020/05/30 07:31:59 media.go:709: INFO - Creating new preview file for DSCN5391.JPG
    2020/05/30 07:32:02 media.go:718: INFO - Preview done for DSCN5391.JPG (conversion time: 2696 ms)
    2020/05/30 07:32:03 media.go:709: INFO - Creating new preview file for DSCN5392.JPG
    2020/05/30 07:32:06 media.go:718: INFO - Preview done for DSCN5392.JPG (conversion time: 2586 ms)
    2020/05/30 07:32:07 media.go:709: INFO - Creating new preview file for DSCN5393.JPG
    2020/05/30 07:32:09 media.go:718: INFO - Preview done for DSCN5393.JPG (conversion time: 2563 ms)
    2020/05/30 07:32:11 media.go:709: INFO - Creating new preview file for DSCN5394.JPG
    2020/05/30 07:32:13 media.go:718: INFO - Preview done for DSCN5394.JPG (conversion time: 2600 ms)
    2020/05/30 07:32:15 media.go:709: INFO - Creating new preview file for DSCN5395.JPG
    2020/05/30 07:32:17 media.go:718: INFO - Preview done for DSCN5395.JPG (conversion time: 2562 ms)
    2020/05/30 07:32:19 media.go:709: INFO - Creating new preview file for DSCN5396.JPG
    2020/05/30 07:32:21 media.go:718: INFO - Preview done for DSCN5396.JPG (conversion time: 2591 ms)
    2020/05/30 07:32:23 media.go:709: INFO - Creating new preview file for DSCN5397.JPG
    2020/05/30 07:32:25 media.go:718: INFO - Preview done for DSCN5397.JPG (conversion time: 2560 ms)
    2020/05/30 07:32:27 media.go:709: INFO - Creating new preview file for DSCN5398.JPG
    2020/05/30 07:32:29 media.go:718: INFO - Preview done for DSCN5398.JPG (conversion time: 2547 ms)
    2020/05/30 07:32:31 media.go:709: INFO - Creating new preview file for DSCN5399.JPG
    2020/05/30 07:32:33 media.go:718: INFO - Preview done for DSCN5399.JPG (conversion time: 2560 ms)
    2020/05/30 07:32:35 media.go:709: INFO - Creating new preview file for exif_rotate/180deg.jpg
    2020/05/30 07:32:37 media.go:718: INFO - Preview done for exif_rotate/180deg.jpg (conversion time: 2501 ms)
    2020/05/30 07:32:39 media.go:709: INFO - Creating new preview file for exif_rotate/mirror.jpg
    2020/05/30 07:32:41 media.go:718: INFO - Preview done for exif_rotate/mirror.jpg (conversion time: 2507 ms)
    2020/05/30 07:32:42 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_270deg.jpg
    2020/05/30 07:32:46 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_270deg.jpg (conversion time: 3139 ms)
    2020/05/30 07:32:47 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_90deg_cw.jpg
    2020/05/30 07:32:50 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_90deg_cw.jpg (conversion time: 3204 ms)
    2020/05/30 07:32:52 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_vertical.jpg
    2020/05/30 07:32:54 media.go:718: INFO - Preview done for exif_rotate/mirror_vertical.jpg (conversion time: 2480 ms)
    2020/05/30 07:32:54 media.go:442: INFO - Creating new thumbnail for exif_rotate/no_exif.jpg
    2020/05/30 07:32:54 media.go:460: INFO - Thumbnail done for exif_rotate/no_exif.jpg (conversion time: 180 ms)
    2020/05/30 07:32:56 media.go:709: INFO - Creating new preview file for exif_rotate/normal.jpg
    2020/05/30 07:32:58 media.go:718: INFO - Preview done for exif_rotate/normal.jpg (conversion time: 2458 ms)
    2020/05/30 07:32:59 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_270deg_cw.jpg
    2020/05/30 07:33:03 media.go:718: INFO - Preview done for exif_rotate/rotate_270deg_cw.jpg (conversion time: 3143 ms)
    2020/05/30 07:33:04 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_90deg_cw.jpg
    2020/05/30 07:33:07 media.go:718: INFO - Preview done for exif_rotate/rotate_90deg_cw.jpg (conversion time: 3198 ms)
    2020/05/30 07:33:07 media.go:442: INFO - Creating new thumbnail for gif.gif
    2020/05/30 07:33:08 media.go:460: INFO - Thumbnail done for gif.gif (conversion time: 851 ms)
    2020/05/30 07:33:09 media.go:709: INFO - Creating new preview file for gif.gif
    2020/05/30 07:33:10 media.go:718: INFO - Preview done for gif.gif (conversion time: 1707 ms)
    2020/05/30 07:33:12 media.go:709: INFO - Creating new preview file for jpeg.jpg
    2020/05/30 07:33:14 media.go:718: INFO - Preview done for jpeg.jpg (conversion time: 2449 ms)
    2020/05/30 07:33:16 media.go:709: INFO - Creating new preview file for jpeg_rotated.jpg
    2020/05/30 07:33:19 media.go:718: INFO - Preview done for jpeg_rotated.jpg (conversion time: 3239 ms)
    2020/05/30 07:33:19 media.go:442: INFO - Creating new thumbnail for png.png
    2020/05/30 07:33:20 media.go:460: INFO - Thumbnail done for png.png (conversion time: 1067 ms)
    2020/05/30 07:33:21 media.go:709: INFO - Creating new preview file for png.png
    2020/05/30 07:33:23 media.go:718: INFO - Preview done for png.png (conversion time: 1777 ms)
    2020/05/30 07:33:23 media.go:442: INFO - Creating new thumbnail for screenshot_browser.jpg
    2020/05/30 07:33:23 media.go:460: INFO - Thumbnail done for screenshot_browser.jpg (conversion time: 116 ms)
    2020/05/30 07:33:23 media.go:442: INFO - Creating new thumbnail for screenshot_mobile.jpg
    2020/05/30 07:33:23 media.go:460: INFO - Thumbnail done for screenshot_mobile.jpg (conversion time: 72 ms)
    2020/05/30 07:33:23 media.go:442: INFO - Creating new thumbnail for screenshot_viewer.jpg
    2020/05/30 07:33:23 media.go:460: INFO - Thumbnail done for screenshot_viewer.jpg (conversion time: 74 ms)
    2020/05/30 07:33:23 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_swipe.jpg
    2020/05/30 07:33:23 media.go:460: INFO - Thumbnail done for screenshot_viewer_swipe.jpg (conversion time: 73 ms)
    2020/05/30 07:33:23 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_zoom.jpg
    2020/05/30 07:33:23 media.go:460: INFO - Thumbnail done for screenshot_viewer_zoom.jpg (conversion time: 55 ms)
    2020/05/30 07:33:23 media.go:442: INFO - Creating new thumbnail for tiff.tiff
    2020/05/30 07:33:24 media.go:460: INFO - Thumbnail done for tiff.tiff (conversion time: 364 ms)
    2020/05/30 07:33:24 media.go:870: INFO - Generating cache took 2 minutes and 56 seconds
    Number of folders: 1
    Number of images: 50
    Number of videos: 0
    Number of images with embedded EXIF: 41
    Number of generated image thumbnails: 9
    Number of generated video thumbnails: 0
    Number of generated image previews: 43
    Number of failed folders: 0
    Number of failed image thumbnails: 0
    Number of failed video thumbnails: 0
    Number of failed image previews: 0
    Number of small images not require preview: 7

### mediaweb_arm64

    pi@rock64:~/performance_test$ ./mediaweb_arm64
    2020/05/30 07:36:09 settings.go:56: INFO - Loading configuration: mediaweb.conf
    2020/05/30 07:36:09 main_common.go:15: INFO - Version: 
    2020/05/30 07:36:09 main_common.go:16: INFO - Build time: Sat May 30 06:50:02 CEST 2020
    2020/05/30 07:36:09 main_common.go:17: INFO - Git hash: 9a22d457e7cb90f7e0504a1a78d2a7212269a102
    2020/05/30 07:36:09 media.go:49: INFO - Media path: testmedia
    2020/05/30 07:36:09 media.go:59: INFO - Cache path: tmpcache
    2020/05/30 07:36:09 media.go:64: INFO - JPEG auto rotate: true
    2020/05/30 07:36:09 media.go:65: INFO - Image preview: true  (max width/height 1280 px)
    2020/05/30 07:36:09 media.go:74: INFO - Video thumbnails supported (ffmpeg installed): true
    2020/05/30 07:36:09 media.go:864: INFO - Pre-generating cache (thumbnails: true, preview: true)
    2020/05/30 07:36:09 watcher.go:50: INFO - Starting media watcher
    2020/05/30 07:36:09 updater.go:48: INFO - Starting updater
    2020/05/30 07:36:09 webapi.go:47: INFO - Starting Web API on port :9999
    2020/05/30 07:36:11 media.go:709: INFO - Creating new preview file for DSCN5369.JPG
    2020/05/30 07:36:13 media.go:718: INFO - Preview done for DSCN5369.JPG (conversion time: 2254 ms)
    2020/05/30 07:36:15 media.go:709: INFO - Creating new preview file for DSCN5370.JPG
    2020/05/30 07:36:17 media.go:718: INFO - Preview done for DSCN5370.JPG (conversion time: 2241 ms)
    2020/05/30 07:36:18 media.go:709: INFO - Creating new preview file for DSCN5371.JPG
    2020/05/30 07:36:20 media.go:718: INFO - Preview done for DSCN5371.JPG (conversion time: 1994 ms)
    2020/05/30 07:36:21 media.go:709: INFO - Creating new preview file for DSCN5372.JPG
    2020/05/30 07:36:23 media.go:718: INFO - Preview done for DSCN5372.JPG (conversion time: 2013 ms)
    2020/05/30 07:36:25 media.go:709: INFO - Creating new preview file for DSCN5373.JPG
    2020/05/30 07:36:27 media.go:718: INFO - Preview done for DSCN5373.JPG (conversion time: 2087 ms)
    2020/05/30 07:36:28 media.go:709: INFO - Creating new preview file for DSCN5374.JPG
    2020/05/30 07:36:30 media.go:718: INFO - Preview done for DSCN5374.JPG (conversion time: 1955 ms)
    2020/05/30 07:36:31 media.go:709: INFO - Creating new preview file for DSCN5375.JPG
    2020/05/30 07:36:34 media.go:718: INFO - Preview done for DSCN5375.JPG (conversion time: 2093 ms)
    2020/05/30 07:36:35 media.go:709: INFO - Creating new preview file for DSCN5376.JPG
    2020/05/30 07:36:37 media.go:718: INFO - Preview done for DSCN5376.JPG (conversion time: 2112 ms)
    2020/05/30 07:36:39 media.go:709: INFO - Creating new preview file for DSCN5377.JPG
    2020/05/30 07:36:41 media.go:718: INFO - Preview done for DSCN5377.JPG (conversion time: 2390 ms)
    2020/05/30 07:36:42 media.go:709: INFO - Creating new preview file for DSCN5378.JPG
    2020/05/30 07:36:44 media.go:718: INFO - Preview done for DSCN5378.JPG (conversion time: 2062 ms)
    2020/05/30 07:36:46 media.go:709: INFO - Creating new preview file for DSCN5379.JPG
    2020/05/30 07:36:48 media.go:718: INFO - Preview done for DSCN5379.JPG (conversion time: 2048 ms)
    2020/05/30 07:36:49 media.go:709: INFO - Creating new preview file for DSCN5380.JPG
    2020/05/30 07:36:51 media.go:718: INFO - Preview done for DSCN5380.JPG (conversion time: 2090 ms)
    2020/05/30 07:36:53 media.go:709: INFO - Creating new preview file for DSCN5381.JPG
    2020/05/30 07:36:55 media.go:718: INFO - Preview done for DSCN5381.JPG (conversion time: 2219 ms)
    2020/05/30 07:36:56 media.go:709: INFO - Creating new preview file for DSCN5382.JPG
    2020/05/30 07:36:58 media.go:718: INFO - Preview done for DSCN5382.JPG (conversion time: 2015 ms)
    2020/05/30 07:36:59 media.go:709: INFO - Creating new preview file for DSCN5383.JPG
    2020/05/30 07:37:02 media.go:718: INFO - Preview done for DSCN5383.JPG (conversion time: 2026 ms)
    2020/05/30 07:37:03 media.go:709: INFO - Creating new preview file for DSCN5384.JPG
    2020/05/30 07:37:05 media.go:718: INFO - Preview done for DSCN5384.JPG (conversion time: 2277 ms)
    2020/05/30 07:37:07 media.go:709: INFO - Creating new preview file for DSCN5385.JPG
    2020/05/30 07:37:09 media.go:718: INFO - Preview done for DSCN5385.JPG (conversion time: 2209 ms)
    2020/05/30 07:37:10 media.go:709: INFO - Creating new preview file for DSCN5386.JPG
    2020/05/30 07:37:13 media.go:718: INFO - Preview done for DSCN5386.JPG (conversion time: 2205 ms)
    2020/05/30 07:37:14 media.go:709: INFO - Creating new preview file for DSCN5387.JPG
    2020/05/30 07:37:16 media.go:718: INFO - Preview done for DSCN5387.JPG (conversion time: 2191 ms)
    2020/05/30 07:37:18 media.go:709: INFO - Creating new preview file for DSCN5388.JPG
    2020/05/30 07:37:20 media.go:718: INFO - Preview done for DSCN5388.JPG (conversion time: 2156 ms)
    2020/05/30 07:37:21 media.go:709: INFO - Creating new preview file for DSCN5389.JPG
    2020/05/30 07:37:23 media.go:718: INFO - Preview done for DSCN5389.JPG (conversion time: 2028 ms)
    2020/05/30 07:37:24 media.go:709: INFO - Creating new preview file for DSCN5390.JPG
    2020/05/30 07:37:26 media.go:718: INFO - Preview done for DSCN5390.JPG (conversion time: 2082 ms)
    2020/05/30 07:37:28 media.go:709: INFO - Creating new preview file for DSCN5391.JPG
    2020/05/30 07:37:30 media.go:718: INFO - Preview done for DSCN5391.JPG (conversion time: 2167 ms)
    2020/05/30 07:37:31 media.go:709: INFO - Creating new preview file for DSCN5392.JPG
    2020/05/30 07:37:33 media.go:718: INFO - Preview done for DSCN5392.JPG (conversion time: 2088 ms)
    2020/05/30 07:37:35 media.go:709: INFO - Creating new preview file for DSCN5393.JPG
    2020/05/30 07:37:37 media.go:718: INFO - Preview done for DSCN5393.JPG (conversion time: 2024 ms)
    2020/05/30 07:37:38 media.go:709: INFO - Creating new preview file for DSCN5394.JPG
    2020/05/30 07:37:40 media.go:718: INFO - Preview done for DSCN5394.JPG (conversion time: 2111 ms)
    2020/05/30 07:37:42 media.go:709: INFO - Creating new preview file for DSCN5395.JPG
    2020/05/30 07:37:44 media.go:718: INFO - Preview done for DSCN5395.JPG (conversion time: 2051 ms)
    2020/05/30 07:37:45 media.go:709: INFO - Creating new preview file for DSCN5396.JPG
    2020/05/30 07:37:47 media.go:718: INFO - Preview done for DSCN5396.JPG (conversion time: 2071 ms)
    2020/05/30 07:37:48 media.go:709: INFO - Creating new preview file for DSCN5397.JPG
    2020/05/30 07:37:50 media.go:718: INFO - Preview done for DSCN5397.JPG (conversion time: 2052 ms)
    2020/05/30 07:37:52 media.go:709: INFO - Creating new preview file for DSCN5398.JPG
    2020/05/30 07:37:54 media.go:718: INFO - Preview done for DSCN5398.JPG (conversion time: 2041 ms)
    2020/05/30 07:37:55 media.go:709: INFO - Creating new preview file for DSCN5399.JPG
    2020/05/30 07:37:57 media.go:718: INFO - Preview done for DSCN5399.JPG (conversion time: 2061 ms)
    2020/05/30 07:37:58 media.go:709: INFO - Creating new preview file for exif_rotate/180deg.jpg
    2020/05/30 07:38:00 media.go:718: INFO - Preview done for exif_rotate/180deg.jpg (conversion time: 2054 ms)
    2020/05/30 07:38:02 media.go:709: INFO - Creating new preview file for exif_rotate/mirror.jpg
    2020/05/30 07:38:04 media.go:718: INFO - Preview done for exif_rotate/mirror.jpg (conversion time: 2033 ms)
    2020/05/30 07:38:05 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_270deg.jpg
    2020/05/30 07:38:08 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_270deg.jpg (conversion time: 2699 ms)
    2020/05/30 07:38:09 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_90deg_cw.jpg
    2020/05/30 07:38:12 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_90deg_cw.jpg (conversion time: 2732 ms)
    2020/05/30 07:38:13 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_vertical.jpg
    2020/05/30 07:38:15 media.go:718: INFO - Preview done for exif_rotate/mirror_vertical.jpg (conversion time: 2011 ms)
    2020/05/30 07:38:15 media.go:442: INFO - Creating new thumbnail for exif_rotate/no_exif.jpg
    2020/05/30 07:38:15 media.go:460: INFO - Thumbnail done for exif_rotate/no_exif.jpg (conversion time: 148 ms)
    2020/05/30 07:38:17 media.go:709: INFO - Creating new preview file for exif_rotate/normal.jpg
    2020/05/30 07:38:19 media.go:718: INFO - Preview done for exif_rotate/normal.jpg (conversion time: 2008 ms)
    2020/05/30 07:38:20 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_270deg_cw.jpg
    2020/05/30 07:38:23 media.go:718: INFO - Preview done for exif_rotate/rotate_270deg_cw.jpg (conversion time: 2704 ms)
    2020/05/30 07:38:24 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_90deg_cw.jpg
    2020/05/30 07:38:27 media.go:718: INFO - Preview done for exif_rotate/rotate_90deg_cw.jpg (conversion time: 2752 ms)
    2020/05/30 07:38:27 media.go:442: INFO - Creating new thumbnail for gif.gif
    2020/05/30 07:38:27 media.go:460: INFO - Thumbnail done for gif.gif (conversion time: 694 ms)
    2020/05/30 07:38:28 media.go:709: INFO - Creating new preview file for gif.gif
    2020/05/30 07:38:29 media.go:718: INFO - Preview done for gif.gif (conversion time: 1227 ms)
    2020/05/30 07:38:30 media.go:709: INFO - Creating new preview file for jpeg.jpg
    2020/05/30 07:38:32 media.go:718: INFO - Preview done for jpeg.jpg (conversion time: 2003 ms)
    2020/05/30 07:38:34 media.go:709: INFO - Creating new preview file for jpeg_rotated.jpg
    2020/05/30 07:38:37 media.go:718: INFO - Preview done for jpeg_rotated.jpg (conversion time: 2808 ms)
    2020/05/30 07:38:37 media.go:442: INFO - Creating new thumbnail for png.png
    2020/05/30 07:38:37 media.go:460: INFO - Thumbnail done for png.png (conversion time: 844 ms)
    2020/05/30 07:38:38 media.go:709: INFO - Creating new preview file for png.png
    2020/05/30 07:38:39 media.go:718: INFO - Preview done for png.png (conversion time: 1290 ms)
    2020/05/30 07:38:39 media.go:442: INFO - Creating new thumbnail for screenshot_browser.jpg
    2020/05/30 07:38:40 media.go:460: INFO - Thumbnail done for screenshot_browser.jpg (conversion time: 68 ms)
    2020/05/30 07:38:40 media.go:442: INFO - Creating new thumbnail for screenshot_mobile.jpg
    2020/05/30 07:38:40 media.go:460: INFO - Thumbnail done for screenshot_mobile.jpg (conversion time: 46 ms)
    2020/05/30 07:38:40 media.go:442: INFO - Creating new thumbnail for screenshot_viewer.jpg
    2020/05/30 07:38:40 media.go:460: INFO - Thumbnail done for screenshot_viewer.jpg (conversion time: 46 ms)
    2020/05/30 07:38:40 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_swipe.jpg
    2020/05/30 07:38:40 media.go:460: INFO - Thumbnail done for screenshot_viewer_swipe.jpg (conversion time: 47 ms)
    2020/05/30 07:38:40 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_zoom.jpg
    2020/05/30 07:38:40 media.go:460: INFO - Thumbnail done for screenshot_viewer_zoom.jpg (conversion time: 40 ms)
    2020/05/30 07:38:40 media.go:442: INFO - Creating new thumbnail for tiff.tiff
    2020/05/30 07:38:40 media.go:460: INFO - Thumbnail done for tiff.tiff (conversion time: 304 ms)
    2020/05/30 07:38:40 media.go:870: INFO - Generating cache took 2 minutes and 31 seconds
    Number of folders: 1
    Number of images: 50
    Number of videos: 0
    Number of images with embedded EXIF: 41
    Number of generated image thumbnails: 9
    Number of generated video thumbnails: 0
    Number of generated image previews: 43
    Number of failed folders: 0
    Number of failed image thumbnails: 0
    Number of failed video thumbnails: 0
    Number of failed image previews: 0
    Number of small images not require preview: 7


## Banana Pi BPI-M1

### mediaweb_arm_v5

    pi@BananaPi2:~/performance_test$ ./mediaweb_arm_v5
    2020/05/30 06:57:34 settings.go:56: INFO - Loading configuration: mediaweb.conf
    2020/05/30 06:57:34 main_common.go:15: INFO - Version: 
    2020/05/30 06:57:34 main_common.go:16: INFO - Build time: Sat May 30 06:44:04 CEST 2020
    2020/05/30 06:57:34 main_common.go:17: INFO - Git hash: 9a22d457e7cb90f7e0504a1a78d2a7212269a102
    2020/05/30 06:57:34 media.go:49: INFO - Media path: testmedia
    2020/05/30 06:57:34 media.go:59: INFO - Cache path: tmpcache
    2020/05/30 06:57:34 media.go:64: INFO - JPEG auto rotate: true
    2020/05/30 06:57:34 media.go:65: INFO - Image preview: true  (max width/height 1280 px)
    2020/05/30 06:57:34 media.go:74: INFO - Video thumbnails supported (ffmpeg installed): false
    2020/05/30 06:57:34 webapi.go:47: INFO - Starting Web API on port :9999
    2020/05/30 06:57:34 media.go:864: INFO - Pre-generating cache (thumbnails: true, preview: true)
    2020/05/30 06:57:34 watcher.go:50: INFO - Starting media watcher
    2020/05/30 06:57:34 updater.go:48: INFO - Starting updater
    2020/05/30 06:57:37 media.go:709: INFO - Creating new preview file for DSCN5369.JPG
    2020/05/30 06:58:28 media.go:718: INFO - Preview done for DSCN5369.JPG (conversion time: 51037 ms)
    2020/05/30 06:58:30 media.go:709: INFO - Creating new preview file for DSCN5370.JPG
    2020/05/30 06:59:20 media.go:718: INFO - Preview done for DSCN5370.JPG (conversion time: 50359 ms)
    2020/05/30 06:59:23 media.go:709: INFO - Creating new preview file for DSCN5371.JPG
    2020/05/30 07:00:12 media.go:718: INFO - Preview done for DSCN5371.JPG (conversion time: 49800 ms)
    2020/05/30 07:00:15 media.go:709: INFO - Creating new preview file for DSCN5372.JPG
    2020/05/30 07:01:04 media.go:718: INFO - Preview done for DSCN5372.JPG (conversion time: 49762 ms)
    2020/05/30 07:01:07 media.go:709: INFO - Creating new preview file for DSCN5373.JPG
    2020/05/30 07:01:57 media.go:718: INFO - Preview done for DSCN5373.JPG (conversion time: 49987 ms)
    2020/05/30 07:01:59 media.go:709: INFO - Creating new preview file for DSCN5374.JPG
    2020/05/30 07:02:48 media.go:718: INFO - Preview done for DSCN5374.JPG (conversion time: 49764 ms)
    2020/05/30 07:02:51 media.go:709: INFO - Creating new preview file for DSCN5375.JPG
    2020/05/30 07:03:41 media.go:718: INFO - Preview done for DSCN5375.JPG (conversion time: 50033 ms)
    2020/05/30 07:03:43 media.go:709: INFO - Creating new preview file for DSCN5376.JPG
    2020/05/30 07:04:33 media.go:718: INFO - Preview done for DSCN5376.JPG (conversion time: 49905 ms)
    2020/05/30 07:04:36 media.go:709: INFO - Creating new preview file for DSCN5377.JPG
    2020/05/30 07:05:26 media.go:718: INFO - Preview done for DSCN5377.JPG (conversion time: 50688 ms)
    2020/05/30 07:05:28 media.go:709: INFO - Creating new preview file for DSCN5378.JPG
    2020/05/30 07:06:19 media.go:718: INFO - Preview done for DSCN5378.JPG (conversion time: 50219 ms)
    2020/05/30 07:06:21 media.go:709: INFO - Creating new preview file for DSCN5379.JPG
    2020/05/30 07:07:09 media.go:718: INFO - Preview done for DSCN5379.JPG (conversion time: 48115 ms)
    2020/05/30 07:07:11 media.go:709: INFO - Creating new preview file for DSCN5380.JPG
    2020/05/30 07:08:02 media.go:718: INFO - Preview done for DSCN5380.JPG (conversion time: 50258 ms)
    2020/05/30 07:08:04 media.go:709: INFO - Creating new preview file for DSCN5381.JPG
    2020/05/30 07:08:55 media.go:718: INFO - Preview done for DSCN5381.JPG (conversion time: 50654 ms)
    2020/05/30 07:08:57 media.go:709: INFO - Creating new preview file for DSCN5382.JPG
    2020/05/30 07:09:43 media.go:718: INFO - Preview done for DSCN5382.JPG (conversion time: 46244 ms)
    2020/05/30 07:09:45 media.go:709: INFO - Creating new preview file for DSCN5383.JPG
    2020/05/30 07:10:33 media.go:718: INFO - Preview done for DSCN5383.JPG (conversion time: 47618 ms)
    2020/05/30 07:10:35 media.go:709: INFO - Creating new preview file for DSCN5384.JPG
    2020/05/30 07:11:26 media.go:718: INFO - Preview done for DSCN5384.JPG (conversion time: 50658 ms)
    2020/05/30 07:11:28 media.go:709: INFO - Creating new preview file for DSCN5385.JPG
    2020/05/30 07:12:19 media.go:718: INFO - Preview done for DSCN5385.JPG (conversion time: 50300 ms)
    2020/05/30 07:12:21 media.go:709: INFO - Creating new preview file for DSCN5386.JPG
    2020/05/30 07:13:11 media.go:718: INFO - Preview done for DSCN5386.JPG (conversion time: 49978 ms)
    2020/05/30 07:13:13 media.go:709: INFO - Creating new preview file for DSCN5387.JPG
    2020/05/30 07:14:04 media.go:718: INFO - Preview done for DSCN5387.JPG (conversion time: 50592 ms)
    2020/05/30 07:14:06 media.go:709: INFO - Creating new preview file for DSCN5388.JPG
    2020/05/30 07:14:57 media.go:718: INFO - Preview done for DSCN5388.JPG (conversion time: 50517 ms)
    2020/05/30 07:14:59 media.go:709: INFO - Creating new preview file for DSCN5389.JPG
    2020/05/30 07:15:49 media.go:718: INFO - Preview done for DSCN5389.JPG (conversion time: 50029 ms)
    2020/05/30 07:15:51 media.go:709: INFO - Creating new preview file for DSCN5390.JPG
    2020/05/30 07:16:41 media.go:718: INFO - Preview done for DSCN5390.JPG (conversion time: 49852 ms)
    2020/05/30 07:16:43 media.go:709: INFO - Creating new preview file for DSCN5391.JPG
    2020/05/30 07:17:34 media.go:718: INFO - Preview done for DSCN5391.JPG (conversion time: 50222 ms)
    2020/05/30 07:17:36 media.go:709: INFO - Creating new preview file for DSCN5392.JPG
    2020/05/30 07:18:24 media.go:718: INFO - Preview done for DSCN5392.JPG (conversion time: 48411 ms)
    2020/05/30 07:18:26 media.go:709: INFO - Creating new preview file for DSCN5393.JPG
    2020/05/30 07:19:16 media.go:718: INFO - Preview done for DSCN5393.JPG (conversion time: 49862 ms)
    2020/05/30 07:19:18 media.go:709: INFO - Creating new preview file for DSCN5394.JPG
    2020/05/30 07:20:08 media.go:718: INFO - Preview done for DSCN5394.JPG (conversion time: 50017 ms)
    2020/05/30 07:20:11 media.go:709: INFO - Creating new preview file for DSCN5395.JPG
    2020/05/30 07:21:01 media.go:718: INFO - Preview done for DSCN5395.JPG (conversion time: 50046 ms)
    2020/05/30 07:21:03 media.go:709: INFO - Creating new preview file for DSCN5396.JPG
    2020/05/30 07:21:52 media.go:718: INFO - Preview done for DSCN5396.JPG (conversion time: 49495 ms)
    2020/05/30 07:21:55 media.go:709: INFO - Creating new preview file for DSCN5397.JPG
    2020/05/30 07:22:45 media.go:718: INFO - Preview done for DSCN5397.JPG (conversion time: 49978 ms)
    2020/05/30 07:22:47 media.go:709: INFO - Creating new preview file for DSCN5398.JPG
    2020/05/30 07:23:37 media.go:718: INFO - Preview done for DSCN5398.JPG (conversion time: 49903 ms)
    2020/05/30 07:23:39 media.go:709: INFO - Creating new preview file for DSCN5399.JPG
    2020/05/30 07:24:29 media.go:718: INFO - Preview done for DSCN5399.JPG (conversion time: 50196 ms)
    2020/05/30 07:24:31 media.go:709: INFO - Creating new preview file for exif_rotate/180deg.jpg
    2020/05/30 07:25:24 media.go:718: INFO - Preview done for exif_rotate/180deg.jpg (conversion time: 52468 ms)
    2020/05/30 07:25:26 media.go:709: INFO - Creating new preview file for exif_rotate/mirror.jpg
    2020/05/30 07:26:18 media.go:718: INFO - Preview done for exif_rotate/mirror.jpg (conversion time: 52529 ms)
    2020/05/30 07:26:21 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_270deg.jpg
    2020/05/30 07:27:15 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_270deg.jpg (conversion time: 54038 ms)
    2020/05/30 07:27:17 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_90deg_cw.jpg
    2020/05/30 07:28:11 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_90deg_cw.jpg (conversion time: 54153 ms)
    2020/05/30 07:28:13 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_vertical.jpg
    2020/05/30 07:29:06 media.go:718: INFO - Preview done for exif_rotate/mirror_vertical.jpg (conversion time: 52582 ms)
    2020/05/30 07:29:06 media.go:442: INFO - Creating new thumbnail for exif_rotate/no_exif.jpg
    2020/05/30 07:29:09 media.go:460: INFO - Thumbnail done for exif_rotate/no_exif.jpg (conversion time: 2985 ms)
    2020/05/30 07:29:11 media.go:709: INFO - Creating new preview file for exif_rotate/normal.jpg
    2020/05/30 07:30:04 media.go:718: INFO - Preview done for exif_rotate/normal.jpg (conversion time: 52647 ms)
    2020/05/30 07:30:06 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_270deg_cw.jpg
    2020/05/30 07:31:00 media.go:718: INFO - Preview done for exif_rotate/rotate_270deg_cw.jpg (conversion time: 53872 ms)
    2020/05/30 07:31:02 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_90deg_cw.jpg
    2020/05/30 07:31:56 media.go:718: INFO - Preview done for exif_rotate/rotate_90deg_cw.jpg (conversion time: 54025 ms)
    2020/05/30 07:31:56 media.go:442: INFO - Creating new thumbnail for gif.gif
    2020/05/30 07:32:18 media.go:460: INFO - Thumbnail done for gif.gif (conversion time: 22039 ms)
    2020/05/30 07:32:19 media.go:709: INFO - Creating new preview file for gif.gif
    2020/05/30 07:33:07 media.go:718: INFO - Preview done for gif.gif (conversion time: 48350 ms)
    2020/05/30 07:33:09 media.go:709: INFO - Creating new preview file for jpeg.jpg
    2020/05/30 07:34:02 media.go:718: INFO - Preview done for jpeg.jpg (conversion time: 52321 ms)
    2020/05/30 07:34:04 media.go:709: INFO - Creating new preview file for jpeg_rotated.jpg
    2020/05/30 07:34:58 media.go:718: INFO - Preview done for jpeg_rotated.jpg (conversion time: 53852 ms)
    2020/05/30 07:34:58 media.go:442: INFO - Creating new thumbnail for png.png
    2020/05/30 07:35:06 media.go:460: INFO - Thumbnail done for png.png (conversion time: 8380 ms)
    2020/05/30 07:35:08 media.go:709: INFO - Creating new preview file for png.png
    2020/05/30 07:35:30 media.go:718: INFO - Preview done for png.png (conversion time: 21899 ms)
    2020/05/30 07:35:30 media.go:442: INFO - Creating new thumbnail for screenshot_browser.jpg
    2020/05/30 07:35:31 media.go:460: INFO - Thumbnail done for screenshot_browser.jpg (conversion time: 1615 ms)
    2020/05/30 07:35:31 media.go:442: INFO - Creating new thumbnail for screenshot_mobile.jpg
    2020/05/30 07:35:32 media.go:460: INFO - Thumbnail done for screenshot_mobile.jpg (conversion time: 874 ms)
    2020/05/30 07:35:32 media.go:442: INFO - Creating new thumbnail for screenshot_viewer.jpg
    2020/05/30 07:35:33 media.go:460: INFO - Thumbnail done for screenshot_viewer.jpg (conversion time: 1161 ms)
    2020/05/30 07:35:34 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_swipe.jpg
    2020/05/30 07:35:35 media.go:460: INFO - Thumbnail done for screenshot_viewer_swipe.jpg (conversion time: 1203 ms)
    2020/05/30 07:35:35 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_zoom.jpg
    2020/05/30 07:35:36 media.go:460: INFO - Thumbnail done for screenshot_viewer_zoom.jpg (conversion time: 1047 ms)
    2020/05/30 07:35:36 media.go:442: INFO - Creating new thumbnail for tiff.tiff
    2020/05/30 07:35:39 media.go:460: INFO - Thumbnail done for tiff.tiff (conversion time: 3505 ms)
    2020/05/30 07:35:40 media.go:870: INFO - Generating cache took 38 minutes and 5 seconds
    Number of folders: 1
    Number of images: 50
    Number of videos: 0
    Number of images with embedded EXIF: 41
    Number of generated image thumbnails: 9
    Number of generated video thumbnails: 0
    Number of generated image previews: 43
    Number of failed folders: 0
    Number of failed image thumbnails: 0
    Number of failed video thumbnails: 0
    Number of failed image previews: 0
    Number of small images not require preview: 7


### mediaweb_arm_v6

    pi@BananaPi2:~/performance_test$ ./mediaweb_arm_v6
    2020/05/30 07:38:57 settings.go:56: INFO - Loading configuration: mediaweb.conf
    2020/05/30 07:38:57 main_common.go:15: INFO - Version: 
    2020/05/30 07:38:57 main_common.go:16: INFO - Build time: Sat May 30 06:48:05 CEST 2020
    2020/05/30 07:38:57 main_common.go:17: INFO - Git hash: 9a22d457e7cb90f7e0504a1a78d2a7212269a102
    2020/05/30 07:38:57 media.go:49: INFO - Media path: testmedia
    2020/05/30 07:38:57 media.go:59: INFO - Cache path: tmpcache
    2020/05/30 07:38:57 media.go:64: INFO - JPEG auto rotate: true
    2020/05/30 07:38:57 media.go:65: INFO - Image preview: true  (max width/height 1280 px)
    2020/05/30 07:38:57 media.go:74: INFO - Video thumbnails supported (ffmpeg installed): false
    2020/05/30 07:38:57 webapi.go:47: INFO - Starting Web API on port :9999
    2020/05/30 07:38:57 media.go:864: INFO - Pre-generating cache (thumbnails: true, preview: true)
    2020/05/30 07:38:57 watcher.go:50: INFO - Starting media watcher
    2020/05/30 07:38:57 updater.go:48: INFO - Starting updater
    2020/05/30 07:39:01 media.go:709: INFO - Creating new preview file for DSCN5369.JPG
    2020/05/30 07:39:06 media.go:718: INFO - Preview done for DSCN5369.JPG (conversion time: 5618 ms)
    2020/05/30 07:39:09 media.go:709: INFO - Creating new preview file for DSCN5370.JPG
    2020/05/30 07:39:14 media.go:718: INFO - Preview done for DSCN5370.JPG (conversion time: 5581 ms)
    2020/05/30 07:39:16 media.go:709: INFO - Creating new preview file for DSCN5371.JPG
    2020/05/30 07:39:22 media.go:718: INFO - Preview done for DSCN5371.JPG (conversion time: 5168 ms)
    2020/05/30 07:39:24 media.go:709: INFO - Creating new preview file for DSCN5372.JPG
    2020/05/30 07:39:29 media.go:718: INFO - Preview done for DSCN5372.JPG (conversion time: 5236 ms)
    2020/05/30 07:39:31 media.go:709: INFO - Creating new preview file for DSCN5373.JPG
    2020/05/30 07:39:36 media.go:718: INFO - Preview done for DSCN5373.JPG (conversion time: 5271 ms)
    2020/05/30 07:39:38 media.go:709: INFO - Creating new preview file for DSCN5374.JPG
    2020/05/30 07:39:44 media.go:718: INFO - Preview done for DSCN5374.JPG (conversion time: 5139 ms)
    2020/05/30 07:39:46 media.go:709: INFO - Creating new preview file for DSCN5375.JPG
    2020/05/30 07:39:51 media.go:718: INFO - Preview done for DSCN5375.JPG (conversion time: 5329 ms)
    2020/05/30 07:39:53 media.go:709: INFO - Creating new preview file for DSCN5376.JPG
    2020/05/30 07:39:59 media.go:718: INFO - Preview done for DSCN5376.JPG (conversion time: 5316 ms)
    2020/05/30 07:40:01 media.go:709: INFO - Creating new preview file for DSCN5377.JPG
    2020/05/30 07:40:07 media.go:718: INFO - Preview done for DSCN5377.JPG (conversion time: 5802 ms)
    2020/05/30 07:40:09 media.go:709: INFO - Creating new preview file for DSCN5378.JPG
    2020/05/30 07:40:15 media.go:718: INFO - Preview done for DSCN5378.JPG (conversion time: 5305 ms)
    2020/05/30 07:40:17 media.go:709: INFO - Creating new preview file for DSCN5379.JPG
    2020/05/30 07:40:22 media.go:718: INFO - Preview done for DSCN5379.JPG (conversion time: 5275 ms)
    2020/05/30 07:40:24 media.go:709: INFO - Creating new preview file for DSCN5380.JPG
    2020/05/30 07:40:30 media.go:718: INFO - Preview done for DSCN5380.JPG (conversion time: 5309 ms)
    2020/05/30 07:40:32 media.go:709: INFO - Creating new preview file for DSCN5381.JPG
    2020/05/30 07:40:38 media.go:718: INFO - Preview done for DSCN5381.JPG (conversion time: 5476 ms)
    2020/05/30 07:40:40 media.go:709: INFO - Creating new preview file for DSCN5382.JPG
    2020/05/30 07:40:45 media.go:718: INFO - Preview done for DSCN5382.JPG (conversion time: 5197 ms)
    2020/05/30 07:40:47 media.go:709: INFO - Creating new preview file for DSCN5383.JPG
    2020/05/30 07:40:52 media.go:718: INFO - Preview done for DSCN5383.JPG (conversion time: 5237 ms)
    2020/05/30 07:40:55 media.go:709: INFO - Creating new preview file for DSCN5384.JPG
    2020/05/30 07:41:00 media.go:718: INFO - Preview done for DSCN5384.JPG (conversion time: 5550 ms)
    2020/05/30 07:41:03 media.go:709: INFO - Creating new preview file for DSCN5385.JPG
    2020/05/30 07:41:08 media.go:718: INFO - Preview done for DSCN5385.JPG (conversion time: 5541 ms)
    2020/05/30 07:41:10 media.go:709: INFO - Creating new preview file for DSCN5386.JPG
    2020/05/30 07:41:16 media.go:718: INFO - Preview done for DSCN5386.JPG (conversion time: 5524 ms)
    2020/05/30 07:41:18 media.go:709: INFO - Creating new preview file for DSCN5387.JPG
    2020/05/30 07:41:24 media.go:718: INFO - Preview done for DSCN5387.JPG (conversion time: 5509 ms)
    2020/05/30 07:41:26 media.go:709: INFO - Creating new preview file for DSCN5388.JPG
    2020/05/30 07:41:32 media.go:718: INFO - Preview done for DSCN5388.JPG (conversion time: 5395 ms)
    2020/05/30 07:41:34 media.go:709: INFO - Creating new preview file for DSCN5389.JPG
    2020/05/30 07:41:39 media.go:718: INFO - Preview done for DSCN5389.JPG (conversion time: 5276 ms)
    2020/05/30 07:41:41 media.go:709: INFO - Creating new preview file for DSCN5390.JPG
    2020/05/30 07:41:47 media.go:718: INFO - Preview done for DSCN5390.JPG (conversion time: 5258 ms)
    2020/05/30 07:41:49 media.go:709: INFO - Creating new preview file for DSCN5391.JPG
    2020/05/30 07:41:54 media.go:718: INFO - Preview done for DSCN5391.JPG (conversion time: 5437 ms)
    2020/05/30 07:41:56 media.go:709: INFO - Creating new preview file for DSCN5392.JPG
    2020/05/30 07:42:02 media.go:718: INFO - Preview done for DSCN5392.JPG (conversion time: 5321 ms)
    2020/05/30 07:42:04 media.go:709: INFO - Creating new preview file for DSCN5393.JPG
    2020/05/30 07:42:09 media.go:718: INFO - Preview done for DSCN5393.JPG (conversion time: 5247 ms)
    2020/05/30 07:42:11 media.go:709: INFO - Creating new preview file for DSCN5394.JPG
    2020/05/30 07:42:17 media.go:718: INFO - Preview done for DSCN5394.JPG (conversion time: 5387 ms)
    2020/05/30 07:42:19 media.go:709: INFO - Creating new preview file for DSCN5395.JPG
    2020/05/30 07:42:24 media.go:718: INFO - Preview done for DSCN5395.JPG (conversion time: 5266 ms)
    2020/05/30 07:42:26 media.go:709: INFO - Creating new preview file for DSCN5396.JPG
    2020/05/30 07:42:32 media.go:718: INFO - Preview done for DSCN5396.JPG (conversion time: 5256 ms)
    2020/05/30 07:42:34 media.go:709: INFO - Creating new preview file for DSCN5397.JPG
    2020/05/30 07:42:39 media.go:718: INFO - Preview done for DSCN5397.JPG (conversion time: 5241 ms)
    2020/05/30 07:42:41 media.go:709: INFO - Creating new preview file for DSCN5398.JPG
    2020/05/30 07:42:46 media.go:718: INFO - Preview done for DSCN5398.JPG (conversion time: 5295 ms)
    2020/05/30 07:42:49 media.go:709: INFO - Creating new preview file for DSCN5399.JPG
    2020/05/30 07:42:54 media.go:718: INFO - Preview done for DSCN5399.JPG (conversion time: 5215 ms)
    2020/05/30 07:42:56 media.go:709: INFO - Creating new preview file for exif_rotate/180deg.jpg
    2020/05/30 07:43:01 media.go:718: INFO - Preview done for exif_rotate/180deg.jpg (conversion time: 5374 ms)
    2020/05/30 07:43:03 media.go:709: INFO - Creating new preview file for exif_rotate/mirror.jpg
    2020/05/30 07:43:09 media.go:718: INFO - Preview done for exif_rotate/mirror.jpg (conversion time: 5308 ms)
    2020/05/30 07:43:11 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_270deg.jpg
    2020/05/30 07:43:18 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_270deg.jpg (conversion time: 6785 ms)
    2020/05/30 07:43:20 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_90deg_cw.jpg
    2020/05/30 07:43:27 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_90deg_cw.jpg (conversion time: 6930 ms)
    2020/05/30 07:43:29 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_vertical.jpg
    2020/05/30 07:43:34 media.go:718: INFO - Preview done for exif_rotate/mirror_vertical.jpg (conversion time: 5307 ms)
    2020/05/30 07:43:34 media.go:442: INFO - Creating new thumbnail for exif_rotate/no_exif.jpg
    2020/05/30 07:43:35 media.go:460: INFO - Thumbnail done for exif_rotate/no_exif.jpg (conversion time: 360 ms)
    2020/05/30 07:43:37 media.go:709: INFO - Creating new preview file for exif_rotate/normal.jpg
    2020/05/30 07:43:42 media.go:718: INFO - Preview done for exif_rotate/normal.jpg (conversion time: 5170 ms)
    2020/05/30 07:43:44 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_270deg_cw.jpg
    2020/05/30 07:43:51 media.go:718: INFO - Preview done for exif_rotate/rotate_270deg_cw.jpg (conversion time: 6773 ms)
    2020/05/30 07:43:53 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_90deg_cw.jpg
    2020/05/30 07:44:00 media.go:718: INFO - Preview done for exif_rotate/rotate_90deg_cw.jpg (conversion time: 6950 ms)
    2020/05/30 07:44:00 media.go:442: INFO - Creating new thumbnail for gif.gif
    2020/05/30 07:44:02 media.go:460: INFO - Thumbnail done for gif.gif (conversion time: 1709 ms)
    2020/05/30 07:44:03 media.go:709: INFO - Creating new preview file for gif.gif
    2020/05/30 07:44:06 media.go:718: INFO - Preview done for gif.gif (conversion time: 3827 ms)
    2020/05/30 07:44:09 media.go:709: INFO - Creating new preview file for jpeg.jpg
    2020/05/30 07:44:14 media.go:718: INFO - Preview done for jpeg.jpg (conversion time: 5028 ms)
    2020/05/30 07:44:16 media.go:709: INFO - Creating new preview file for jpeg_rotated.jpg
    2020/05/30 07:44:23 media.go:718: INFO - Preview done for jpeg_rotated.jpg (conversion time: 6947 ms)
    2020/05/30 07:44:23 media.go:442: INFO - Creating new thumbnail for png.png
    2020/05/30 07:44:25 media.go:460: INFO - Thumbnail done for png.png (conversion time: 1849 ms)
    2020/05/30 07:44:26 media.go:709: INFO - Creating new preview file for png.png
    2020/05/30 07:44:30 media.go:718: INFO - Preview done for png.png (conversion time: 3485 ms)
    2020/05/30 07:44:30 media.go:442: INFO - Creating new thumbnail for screenshot_browser.jpg
    2020/05/30 07:44:30 media.go:460: INFO - Thumbnail done for screenshot_browser.jpg (conversion time: 187 ms)
    2020/05/30 07:44:30 media.go:442: INFO - Creating new thumbnail for screenshot_mobile.jpg
    2020/05/30 07:44:30 media.go:460: INFO - Thumbnail done for screenshot_mobile.jpg (conversion time: 119 ms)
    2020/05/30 07:44:30 media.go:442: INFO - Creating new thumbnail for screenshot_viewer.jpg
    2020/05/30 07:44:30 media.go:460: INFO - Thumbnail done for screenshot_viewer.jpg (conversion time: 130 ms)
    2020/05/30 07:44:30 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_swipe.jpg
    2020/05/30 07:44:30 media.go:460: INFO - Thumbnail done for screenshot_viewer_swipe.jpg (conversion time: 140 ms)
    2020/05/30 07:44:30 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_zoom.jpg
    2020/05/30 07:44:30 media.go:460: INFO - Thumbnail done for screenshot_viewer_zoom.jpg (conversion time: 118 ms)
    2020/05/30 07:44:31 media.go:442: INFO - Creating new thumbnail for tiff.tiff
    2020/05/30 07:44:31 media.go:460: INFO - Thumbnail done for tiff.tiff (conversion time: 639 ms)
    2020/05/30 07:44:32 media.go:870: INFO - Generating cache took 5 minutes and 34 seconds
    Number of folders: 1
    Number of images: 50
    Number of videos: 0
    Number of images with embedded EXIF: 41
    Number of generated image thumbnails: 9
    Number of generated video thumbnails: 0
    Number of generated image previews: 43
    Number of failed folders: 0
    Number of failed image thumbnails: 0
    Number of failed video thumbnails: 0
    Number of failed image previews: 0
    Number of small images not require preview: 7


### mediaweb_arm_v7

    pi@BananaPi2:~/performance_test$ ./mediaweb_arm_v7
    2020/05/30 07:46:34 settings.go:56: INFO - Loading configuration: mediaweb.conf
    2020/05/30 07:46:34 main_common.go:15: INFO - Version: 
    2020/05/30 07:46:34 main_common.go:16: INFO - Build time: Sat May 30 06:49:04 CEST 2020
    2020/05/30 07:46:34 main_common.go:17: INFO - Git hash: 9a22d457e7cb90f7e0504a1a78d2a7212269a102
    2020/05/30 07:46:34 media.go:49: INFO - Media path: testmedia
    2020/05/30 07:46:34 media.go:59: INFO - Cache path: tmpcache
    2020/05/30 07:46:34 media.go:64: INFO - JPEG auto rotate: true
    2020/05/30 07:46:34 media.go:65: INFO - Image preview: true  (max width/height 1280 px)
    2020/05/30 07:46:34 media.go:74: INFO - Video thumbnails supported (ffmpeg installed): false
    2020/05/30 07:46:34 webapi.go:47: INFO - Starting Web API on port :9999
    2020/05/30 07:46:34 media.go:864: INFO - Pre-generating cache (thumbnails: true, preview: true)
    2020/05/30 07:46:34 watcher.go:50: INFO - Starting media watcher
    2020/05/30 07:46:34 updater.go:48: INFO - Starting updater
    2020/05/30 07:46:36 media.go:709: INFO - Creating new preview file for DSCN5369.JPG
    2020/05/30 07:46:42 media.go:718: INFO - Preview done for DSCN5369.JPG (conversion time: 6171 ms)
    2020/05/30 07:46:45 media.go:709: INFO - Creating new preview file for DSCN5370.JPG
    2020/05/30 07:46:50 media.go:718: INFO - Preview done for DSCN5370.JPG (conversion time: 5455 ms)
    2020/05/30 07:46:52 media.go:709: INFO - Creating new preview file for DSCN5371.JPG
    2020/05/30 07:46:57 media.go:718: INFO - Preview done for DSCN5371.JPG (conversion time: 5099 ms)
    2020/05/30 07:47:00 media.go:709: INFO - Creating new preview file for DSCN5372.JPG
    2020/05/30 07:47:05 media.go:718: INFO - Preview done for DSCN5372.JPG (conversion time: 5113 ms)
    2020/05/30 07:47:07 media.go:709: INFO - Creating new preview file for DSCN5373.JPG
    2020/05/30 07:47:12 media.go:718: INFO - Preview done for DSCN5373.JPG (conversion time: 5172 ms)
    2020/05/30 07:47:14 media.go:709: INFO - Creating new preview file for DSCN5374.JPG
    2020/05/30 07:47:19 media.go:718: INFO - Preview done for DSCN5374.JPG (conversion time: 4978 ms)
    2020/05/30 07:47:21 media.go:709: INFO - Creating new preview file for DSCN5375.JPG
    2020/05/30 07:47:26 media.go:718: INFO - Preview done for DSCN5375.JPG (conversion time: 5173 ms)
    2020/05/30 07:47:28 media.go:709: INFO - Creating new preview file for DSCN5376.JPG
    2020/05/30 07:47:34 media.go:718: INFO - Preview done for DSCN5376.JPG (conversion time: 5219 ms)
    2020/05/30 07:47:36 media.go:709: INFO - Creating new preview file for DSCN5377.JPG
    2020/05/30 07:47:42 media.go:718: INFO - Preview done for DSCN5377.JPG (conversion time: 5645 ms)
    2020/05/30 07:47:44 media.go:709: INFO - Creating new preview file for DSCN5378.JPG
    2020/05/30 07:47:49 media.go:718: INFO - Preview done for DSCN5378.JPG (conversion time: 5165 ms)
    2020/05/30 07:47:51 media.go:709: INFO - Creating new preview file for DSCN5379.JPG
    2020/05/30 07:47:56 media.go:718: INFO - Preview done for DSCN5379.JPG (conversion time: 5139 ms)
    2020/05/30 07:47:59 media.go:709: INFO - Creating new preview file for DSCN5380.JPG
    2020/05/30 07:48:04 media.go:718: INFO - Preview done for DSCN5380.JPG (conversion time: 5227 ms)
    2020/05/30 07:48:06 media.go:709: INFO - Creating new preview file for DSCN5381.JPG
    2020/05/30 07:48:12 media.go:718: INFO - Preview done for DSCN5381.JPG (conversion time: 5351 ms)
    2020/05/30 07:48:14 media.go:709: INFO - Creating new preview file for DSCN5382.JPG
    2020/05/30 07:48:19 media.go:718: INFO - Preview done for DSCN5382.JPG (conversion time: 5101 ms)
    2020/05/30 07:48:21 media.go:709: INFO - Creating new preview file for DSCN5383.JPG
    2020/05/30 07:48:26 media.go:718: INFO - Preview done for DSCN5383.JPG (conversion time: 5071 ms)
    2020/05/30 07:48:28 media.go:709: INFO - Creating new preview file for DSCN5384.JPG
    2020/05/30 07:48:34 media.go:718: INFO - Preview done for DSCN5384.JPG (conversion time: 5509 ms)
    2020/05/30 07:48:36 media.go:709: INFO - Creating new preview file for DSCN5385.JPG
    2020/05/30 07:48:41 media.go:718: INFO - Preview done for DSCN5385.JPG (conversion time: 5348 ms)
    2020/05/30 07:48:43 media.go:709: INFO - Creating new preview file for DSCN5386.JPG
    2020/05/30 07:48:49 media.go:718: INFO - Preview done for DSCN5386.JPG (conversion time: 5402 ms)
    2020/05/30 07:48:51 media.go:709: INFO - Creating new preview file for DSCN5387.JPG
    2020/05/30 07:48:57 media.go:718: INFO - Preview done for DSCN5387.JPG (conversion time: 5378 ms)
    2020/05/30 07:48:59 media.go:709: INFO - Creating new preview file for DSCN5388.JPG
    2020/05/30 07:49:04 media.go:718: INFO - Preview done for DSCN5388.JPG (conversion time: 5332 ms)
    2020/05/30 07:49:06 media.go:709: INFO - Creating new preview file for DSCN5389.JPG
    2020/05/30 07:49:11 media.go:718: INFO - Preview done for DSCN5389.JPG (conversion time: 5159 ms)
    2020/05/30 07:49:14 media.go:709: INFO - Creating new preview file for DSCN5390.JPG
    2020/05/30 07:49:19 media.go:718: INFO - Preview done for DSCN5390.JPG (conversion time: 5160 ms)
    2020/05/30 07:49:21 media.go:709: INFO - Creating new preview file for DSCN5391.JPG
    2020/05/30 07:49:26 media.go:718: INFO - Preview done for DSCN5391.JPG (conversion time: 5328 ms)
    2020/05/30 07:49:28 media.go:709: INFO - Creating new preview file for DSCN5392.JPG
    2020/05/30 07:49:34 media.go:718: INFO - Preview done for DSCN5392.JPG (conversion time: 5199 ms)
    2020/05/30 07:49:36 media.go:709: INFO - Creating new preview file for DSCN5393.JPG
    2020/05/30 07:49:41 media.go:718: INFO - Preview done for DSCN5393.JPG (conversion time: 5094 ms)
    2020/05/30 07:49:43 media.go:709: INFO - Creating new preview file for DSCN5394.JPG
    2020/05/30 07:49:48 media.go:718: INFO - Preview done for DSCN5394.JPG (conversion time: 5180 ms)
    2020/05/30 07:49:50 media.go:709: INFO - Creating new preview file for DSCN5395.JPG
    2020/05/30 07:49:55 media.go:718: INFO - Preview done for DSCN5395.JPG (conversion time: 5174 ms)
    2020/05/30 07:49:57 media.go:709: INFO - Creating new preview file for DSCN5396.JPG
    2020/05/30 07:50:03 media.go:718: INFO - Preview done for DSCN5396.JPG (conversion time: 5180 ms)
    2020/05/30 07:50:05 media.go:709: INFO - Creating new preview file for DSCN5397.JPG
    2020/05/30 07:50:10 media.go:718: INFO - Preview done for DSCN5397.JPG (conversion time: 5144 ms)
    2020/05/30 07:50:12 media.go:709: INFO - Creating new preview file for DSCN5398.JPG
    2020/05/30 07:50:17 media.go:718: INFO - Preview done for DSCN5398.JPG (conversion time: 5184 ms)
    2020/05/30 07:50:19 media.go:709: INFO - Creating new preview file for DSCN5399.JPG
    2020/05/30 07:50:24 media.go:718: INFO - Preview done for DSCN5399.JPG (conversion time: 5168 ms)
    2020/05/30 07:50:26 media.go:709: INFO - Creating new preview file for exif_rotate/180deg.jpg
    2020/05/30 07:50:32 media.go:718: INFO - Preview done for exif_rotate/180deg.jpg (conversion time: 5254 ms)
    2020/05/30 07:50:34 media.go:709: INFO - Creating new preview file for exif_rotate/mirror.jpg
    2020/05/30 07:50:39 media.go:718: INFO - Preview done for exif_rotate/mirror.jpg (conversion time: 5190 ms)
    2020/05/30 07:50:41 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_270deg.jpg
    2020/05/30 07:50:48 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_270deg.jpg (conversion time: 6870 ms)
    2020/05/30 07:50:50 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_rotate_90deg_cw.jpg
    2020/05/30 07:50:57 media.go:718: INFO - Preview done for exif_rotate/mirror_rotate_90deg_cw.jpg (conversion time: 6755 ms)
    2020/05/30 07:50:59 media.go:709: INFO - Creating new preview file for exif_rotate/mirror_vertical.jpg
    2020/05/30 07:51:04 media.go:718: INFO - Preview done for exif_rotate/mirror_vertical.jpg (conversion time: 5223 ms)
    2020/05/30 07:51:04 media.go:442: INFO - Creating new thumbnail for exif_rotate/no_exif.jpg
    2020/05/30 07:51:04 media.go:460: INFO - Thumbnail done for exif_rotate/no_exif.jpg (conversion time: 352 ms)
    2020/05/30 07:51:07 media.go:709: INFO - Creating new preview file for exif_rotate/normal.jpg
    2020/05/30 07:51:12 media.go:718: INFO - Preview done for exif_rotate/normal.jpg (conversion time: 5033 ms)
    2020/05/30 07:51:14 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_270deg_cw.jpg
    2020/05/30 07:51:20 media.go:718: INFO - Preview done for exif_rotate/rotate_270deg_cw.jpg (conversion time: 6576 ms)
    2020/05/30 07:51:22 media.go:709: INFO - Creating new preview file for exif_rotate/rotate_90deg_cw.jpg
    2020/05/30 07:51:29 media.go:718: INFO - Preview done for exif_rotate/rotate_90deg_cw.jpg (conversion time: 6966 ms)
    2020/05/30 07:51:29 media.go:442: INFO - Creating new thumbnail for gif.gif
    2020/05/30 07:51:31 media.go:460: INFO - Thumbnail done for gif.gif (conversion time: 1688 ms)
    2020/05/30 07:51:32 media.go:709: INFO - Creating new preview file for gif.gif
    2020/05/30 07:51:36 media.go:718: INFO - Preview done for gif.gif (conversion time: 3809 ms)
    2020/05/30 07:51:38 media.go:709: INFO - Creating new preview file for jpeg.jpg
    2020/05/30 07:51:43 media.go:718: INFO - Preview done for jpeg.jpg (conversion time: 4959 ms)
    2020/05/30 07:51:45 media.go:709: INFO - Creating new preview file for jpeg_rotated.jpg
    2020/05/30 07:51:52 media.go:718: INFO - Preview done for jpeg_rotated.jpg (conversion time: 6899 ms)
    2020/05/30 07:51:52 media.go:442: INFO - Creating new thumbnail for png.png
    2020/05/30 07:51:54 media.go:460: INFO - Thumbnail done for png.png (conversion time: 1854 ms)
    2020/05/30 07:51:55 media.go:709: INFO - Creating new preview file for png.png
    2020/05/30 07:51:59 media.go:718: INFO - Preview done for png.png (conversion time: 3387 ms)
    2020/05/30 07:51:59 media.go:442: INFO - Creating new thumbnail for screenshot_browser.jpg
    2020/05/30 07:51:59 media.go:460: INFO - Thumbnail done for screenshot_browser.jpg (conversion time: 189 ms)
    2020/05/30 07:51:59 media.go:442: INFO - Creating new thumbnail for screenshot_mobile.jpg
    2020/05/30 07:51:59 media.go:460: INFO - Thumbnail done for screenshot_mobile.jpg (conversion time: 112 ms)
    2020/05/30 07:51:59 media.go:442: INFO - Creating new thumbnail for screenshot_viewer.jpg
    2020/05/30 07:51:59 media.go:460: INFO - Thumbnail done for screenshot_viewer.jpg (conversion time: 128 ms)
    2020/05/30 07:51:59 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_swipe.jpg
    2020/05/30 07:51:59 media.go:460: INFO - Thumbnail done for screenshot_viewer_swipe.jpg (conversion time: 144 ms)
    2020/05/30 07:51:59 media.go:442: INFO - Creating new thumbnail for screenshot_viewer_zoom.jpg
    2020/05/30 07:51:59 media.go:460: INFO - Thumbnail done for screenshot_viewer_zoom.jpg (conversion time: 112 ms)
    2020/05/30 07:51:59 media.go:442: INFO - Creating new thumbnail for tiff.tiff
    2020/05/30 07:52:00 media.go:460: INFO - Thumbnail done for tiff.tiff (conversion time: 624 ms)
    2020/05/30 07:52:01 media.go:870: INFO - Generating cache took 5 minutes and 26 seconds
    Number of folders: 1
    Number of images: 50
    Number of videos: 0
    Number of images with embedded EXIF: 41
    Number of generated image thumbnails: 9
    Number of generated video thumbnails: 0
    Number of generated image previews: 43
    Number of failed folders: 0
    Number of failed image thumbnails: 0
    Number of failed video thumbnails: 0
    Number of failed image previews: 0
    Number of small images not require preview: 7