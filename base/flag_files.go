package base

type FlagFiles []string

func (flagFiles *FlagFiles) String() string {
	if flagFiles == nil || len(*flagFiles) == 0 {
		return ""
	}

	return flagFiles.String()
}

func (flagFiles *FlagFiles) Set(value string) error {
	if value != "" {
		*flagFiles = append(*flagFiles, value)
	}

	return nil
}
