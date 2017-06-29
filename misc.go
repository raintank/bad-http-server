package main

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// closestRatio returns whether you should pick a, if you wonder whether
// adding 1 to a or 1 to b will result in a a/(a+b) ratio closer to ratio.
func closestRatio(ratio, a, b float64) bool {
	aPlusOne := (a + 1) / (a + b + 1)
	bPlusOne := a / (a + b + 1)
	return math.Abs(aPlusOne-ratio) < math.Abs(bPlusOne-ratio)
}

// for urls of form /<type>/<ratio>, return key (which is `type/ratio`) and ratio
func parseKeyRatio(path, format string) (string, int, error) {
	if strings.Count(path, "/") != 2 {
		return "", 0, fmt.Errorf("bad format. expected %q", format)
	}
	i := strings.LastIndex(path, "/")
	ratioStr := path[i+1:] // may be "" but strconv.Atoi will catch that
	badRatio, err := strconv.Atoi(ratioStr)
	if err != nil || badRatio < 0 || badRatio > 100 {
		return "", 0, errors.New("bad ratio (should be a percentage between 0 and 100, inclusive)")
	}
	return path[1:], badRatio, nil
}

// for urls of form /type/key or /type/key/ or /type/key/ratio, return key (which is `/type/key`) and ratio if specified
func parseDynamicKeyRatio(path, format string) (string, int, error) {
	badRatio := -1
	cnt := strings.Count(path, "/")
	if cnt == 2 {
		// key must not be ""
		if path[len(path)-1] == '/' {
			return "", 0, fmt.Errorf("bad format. expected %q", format)
		}
		return path[1:], -1, nil
	} else if cnt == 3 {
		i := strings.LastIndex(path, "/")
		ratioStr := path[i+1:]
		if ratioStr != "" {
			var err error
			badRatio, err = strconv.Atoi(ratioStr)
			if err != nil || badRatio < 0 || badRatio > 100 {
				return "", 0, errors.New("bad ratio (should be a percentage between 0 and 100, inclusive)")
			}
		}
		return path[1:i], badRatio, nil
	}
	return "", 0, fmt.Errorf("bad format. expected %q", format)
}
