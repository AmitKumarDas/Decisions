```go
func getMultipathDisk(path string) (string, error) {
	// Follow link to destination directory
	devicePath, err := os.Readlink(path)
	if err != nil {
		debug.Printf("Failed reading link for multipath disk: %s -- error: %s\n", path, err.Error())
		return "", err
	}
	sdevice := filepath.Base(devicePath)
	// If destination directory is already identified as a multipath device,
	// just return its path
	if strings.HasPrefix(sdevice, "dm-") {
		return path, nil
	}
	// Fallback to iterating through all the entries under /sys/block/dm-* and
	// check to see if any have an entry under /sys/block/dm-*/slaves matching
	// the device the symlink was pointing at
	dmPaths, _ := filepath.Glob("/sys/block/dm-*")
	for _, dmPath := range dmPaths {
		sdevices, _ := filepath.Glob(filepath.Join(dmPath, "slaves", "*"))
		for _, spath := range sdevices {
			s := filepath.Base(spath)
			if sdevice == s {
				// We've found a matching entry, return the path for the
				// dm-* device it was found under
				p := filepath.Join("/dev", filepath.Base(dmPath))
				debug.Printf("Found matching multipath device: %s under dm-* device path %s", sdevice, dmPath)
				return p, nil
			}
		}
	}
	return "", fmt.Errorf("Couldn't find dm-* path for path: %s, found non dm-* path: %s", path, devicePath)
}
```

ref - https://github.com/kubernetes-csi/csi-lib-iscsi/
