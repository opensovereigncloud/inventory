package file

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func ToString(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read file %s", path)
	}

	contentsString := string(contents)
	trimmedString := strings.TrimSpace(contentsString)

	return trimmedString, nil
}

func ToInt(path string) (int, error) {
	fileString, err := ToString(path)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to read string from file %s ", path)
	}

	num, err := strconv.Atoi(fileString)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to convert %s file to int", fileString)
	}

	return num, nil
}

func ToUint64(path string) (uint64, error) {
	fileString, err := ToString(path)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to read string from file %s ", path)
	}

	num, err := strconv.ParseUint(fileString, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to convert %s file to int", fileString)
	}

	return num, nil
}

func ToBool(path string) (bool, error) {
	num, err := ToInt(path)
	if err != nil {
		return false, errors.Wrapf(err, "unable to read int from file %s ", path)
	}

	return num == 1, nil
}
