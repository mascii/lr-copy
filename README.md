# lr-copy
`lr-copy` is a Go application that makes copies of JPEG/RAW files divided into directories by date.
The date format of the destination directory is `YYYYY/YYYYY-MM-DD`, which is similar to Adobe Lightroom.

```
$ tree /Volumes/SDCARD/DCIM/100OLYMP
/Volumes/SDCARD/DCIM/100OLYMP
├── P2010001.JPG
├── P2020002.JPG
├── P2030003.JPG
├── P2030004.JPG
└── P2040005.JPG

$ lr-copy --src /Volumes/SDCARD/DCIM/100OLYMP --dst /Users/Bob/Pictures/
/Volumes/SDCARD/DCIM/100OLYMP/P2010001.JPG => /Users/Bob/Pictures/2024/2024-02-01/P2010001.JPG
/Volumes/SDCARD/DCIM/100OLYMP/P2020002.JPG => /Users/Bob/Pictures/2024/2024-02-02/P2020002.JPG
/Volumes/SDCARD/DCIM/100OLYMP/P2030003.JPG => /Users/Bob/Pictures/2024/2024-02-03/P2030003.JPG
/Volumes/SDCARD/DCIM/100OLYMP/P2030004.JPG => /Users/Bob/Pictures/2024/2024-02-03/P2030004.JPG
/Volumes/SDCARD/DCIM/100OLYMP/P2040005.JPG => /Users/Bob/Pictures/2024/2024-02-04/P2040005.JPG
-----------------------------------
Continue with the process? (y/N): y
Copied /Volumes/SDCARD/DCIM/100OLYMP/P2010001.JPG to /Users/Bob/Pictures/2024/2024-02-01/P2010001.JPG
Copied /Volumes/SDCARD/DCIM/100OLYMP/P2020002.JPG to /Users/Bob/Pictures/2024/2024-02-02/P2020002.JPG
Copied /Volumes/SDCARD/DCIM/100OLYMP/P2030003.JPG to /Users/Bob/Pictures/2024/2024-02-03/P2030003.JPG
Copied /Volumes/SDCARD/DCIM/100OLYMP/P2030004.JPG to /Users/Bob/Pictures/2024/2024-02-03/P2030004.JPG
Copied /Volumes/SDCARD/DCIM/100OLYMP/P2040005.JPG to /Users/Bob/Pictures/2024/2024-02-04/P2040005.JPG
-----------------------------------
Copied: 5, Skipped: 0, Failed: 0

$ tree /Users/Bob/Pictures/
/Users/Bob/Pictures/
└── 2024
    ├── 2024-02-01
    │   └── P2010001.JPG
    ├── 2024-02-02
    │   └── P2020002.JPG
    ├── 2024-02-03
    │   ├── P2030003.JPG
    │   └── P2030004.JPG
    └── 2024-02-04
        └── P2040005.JPG
```

## Usage
```
lr-copy -src=<source_directory_path> -dst=<destination_directory_path> [-overwrite] [-separate=false]
```

### Flags
- `-src`: Source directory path. This is a required flag.
- `-dst`: Destination directory base path. This is a required flag.
- `-overwrite`: Overwrite existing files. This is an optional flag. By default, it is set to false.
- `-separate`: Separate directory excepting JPEG by file type (e.g. ORF, ARW, etc.). This is an optional flag. By default, it is set to true.
