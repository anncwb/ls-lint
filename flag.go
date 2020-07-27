package main

type flagSlice []string

func (fs *flagSlice) String() string {
	if fs == nil || len(*fs) == 0 {
		return ""
	}

	return fs.String()
}

func (fs *flagSlice) Set(value string) error {
	if value != "" {
		*fs = append(*fs, value)
	}

	return nil
}
