package files

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func FilesForPatterns(patterns []string) ([]string, error) {
	var tmp = make(map[string]struct{})

	for _, p := range patterns {
		files, err := filepath.Glob(p)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			tmp[f] = struct{}{}
		}
	}

	var res = make([]string, 0, len(tmp))
	for k := range tmp {
		res = append(res, k)
	}

	return res, nil
}

// git file hash
func GitHash(content []byte) string {
	// git object id: sha1(b'blob %d\0%s' % (len(b), b)).hexdigest()
	h := sha1.New()
	h.Write([]byte("blob "))
	h.Write([]byte(strconv.FormatInt(int64(len(content)), 10)))
	h.Write([]byte{0})
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}

func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	cnt, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return GitHash(cnt), nil
}
