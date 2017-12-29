package cdtacctwebscraper

import (
	"fmt"
	"regexp"
	"strings"
)

func (ac AccountConfig) VariableSubstitute(commands Commands) (new Commands, err error) {
	re := regexp.MustCompile("{[A-Za-z]*}")

	for _, command := range commands {
		newCommand := Command{Instruction: command.Instruction}
		if command.Commands != nil {
			newCommand.Commands, err = ac.VariableSubstitute(command.Commands)
			if err != nil {
				return
			}
		}
		for _, param := range command.Params {
			tags := re.FindAllString(param, -1)
			for _, tag := range tags {
				switch strings.ToLower(tag) {
				case "{username}":
					param = strings.Replace(param, tag, ac.Username, 1)
				case "{password}":
					param = strings.Replace(param, tag, ac.Password, 1)
				case "{accountname}":
					param = strings.Replace(param, tag, ac.AccountName, 1)
				case "{accountid}":
					param = strings.Replace(param, tag, ac.AccountId, 1)
				case "{businessid}":
					param = strings.Replace(param, tag, ac.BusinessId, 1)
				case "{compactname}":
					param = strings.Replace(param, tag, ac.CompactName, 1)
				case "{url}":
					param = strings.Replace(param, tag, ac.Url, 1)
				case "{href}":
					// substitute later, during run tasks
				default:
					err = fmt.Errorf("Unrecognized tag %s", tag)
					return
				}
			}
			newCommand.Params = append(newCommand.Params, param)
		}
		new = append(new, newCommand)
	}
	return
}

func HrefSubstitute(commands Commands, href string) (new Commands, err error) {
	re := regexp.MustCompile("{[A-Za-z]*}")

	for _, command := range commands {
		newCommand := Command{Instruction: command.Instruction}
		for _, param := range command.Params {
			tags := re.FindAllString(param, -1)
			for _, tag := range tags {
				if strings.ToLower(tag) == "{href}" {
					param = strings.Replace(param, tag, href, 1)
				} else {
					err = fmt.Errorf("Unrecognized tag %s", tag)
					return
				}
			}
			newCommand.Params = append(newCommand.Params, param)
		}
		new = append(new, newCommand)
	}
	return
}
