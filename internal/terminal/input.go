package terminal

import (
	"bufio"
	"os"
	"strings"
)

func InputString() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	text, err := reader.ReadString('\n')
	if err != nil {

		return "", err
	}

	text = strings.TrimSpace(text)
	return text, nil
}

func InputConfirm() (string, bool, error) {
	res, err := InputString()
	if err != nil {
		return "", false, err
	}
	res = strings.ToLower(res)
	yes := false
	if res == "да" || res == "д" || res == "y" || res == "yes" {
		yes = true
	}

	return res, yes, nil
}
