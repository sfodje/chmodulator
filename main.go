package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var argument, owner, group, other string
var permissionsKey = map[string]int{"r": 4, "w": 2, "x": 1, "-": 0}

func init() {
	if len(os.Args) == 2 && !strings.Contains(os.Args[1], "=") {
		argument = os.Args[1]
		os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
	}
	// Usage overrides flag.Usage()
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n%s [permission e.g. -rwxrwxr-x or 775]\nor\n%s <flags>\n",
			os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&owner, "owner", "---", "owner level permissions")
	flag.StringVar(&group, "group", "---", "group level permissions")
	flag.StringVar(&other, "other", "---", "other level permissions")
}

func main() {
	flag.Parse()
	if argument == "" && owner == "---" && group == "---" && other == "---" {
		flag.Usage()
		return
	}
	var permStr = "-" + owner + group + other
	if argument != "" {
		permStr = argument
	}

	if matched, _ := regexp.MatchString("^\\d+$", permStr); matched {
		output, err := convertUIntStr(permStr)
		checkError(err, "")
		fmt.Println(output)
		return
	}

	if matched, _ := regexp.MatchString("^[-rwx]+$", permStr); matched {
		output, err := convertRWXStr(permStr)
		checkError(err, "")
		fmt.Println(output)
		return
	}

	fmt.Println("Invalid argument!")
	flag.Usage()
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatal(err, msg)
	}
	return
}

func reduce(arr []string) string {
	var output string

	for _, v := range arr {
		output += strings.Replace(v, "-", "", -1)
	}
	return output
}

func convertRWXStr(str string) (string, error) {
	var permArray = strings.Split(str, "")
	if len(permArray) < 9 {
		return "", fmt.Errorf("Invalid argument")
	}
	if len(permArray) > 9 {
		permArray = permArray[1:]
	}
	str = "-" + strings.Join(permArray[:9], "")

	var value int
	var output string
	for i, v := range permArray {
		value += permissionsKey[v]
		if (i+1)%3 == 0 {
			output += fmt.Sprintf("%d", value)
			value = 0
		}
	}

	var ownerPerms = permArray[0:3]
	var groupPerms = permArray[3:6]
	var otherPerms = permArray[6:9]

	return fmt.Sprintf("%s\n%s\nOwner: %s\nGroup: %s\nOther: %s\n",
		str, "0"+output, reduce(ownerPerms), reduce(groupPerms), reduce(otherPerms)), nil

}

func convertUIntStr(str string) (string, error) {
	var permArray = strings.Split(str, "")
	if len(permArray) < 3 {
		return "", fmt.Errorf("Invalid argument")
	}
	if len(permArray) > 3 {
		permArray = permArray[1:]
	}
	str = "0" + strings.Join(permArray[:3], "")
	var output = []string{"-"}
	for _, v := range permArray {
		i, err := strconv.Atoi(v)
		if err != nil {
			return "", err
		}
		if i -= permissionsKey["r"]; i >= 0 {
			output = append(output, "r")
		} else {
			i += permissionsKey["r"]
			output = append(output, "-")
		}

		if i -= permissionsKey["w"]; i >= 0 {
			output = append(output, "w")
		} else {
			i += permissionsKey["w"]
			output = append(output, "-")
		}

		if i -= permissionsKey["x"]; i >= 0 {
			output = append(output, "x")
		} else {
			i += permissionsKey["x"]
			output = append(output, "-")
		}
	}

	var ownerPerms = output[1:][0:3]
	var groupPerms = output[1:][3:6]
	var otherPerms = output[1:][6:9]

	return fmt.Sprintf("%s\n%s\nOwner: %s\nGroup: %s\nOther: %s\n",
		strings.Join(output, ""), str, reduce(ownerPerms), reduce(groupPerms), reduce(otherPerms)), nil
}
