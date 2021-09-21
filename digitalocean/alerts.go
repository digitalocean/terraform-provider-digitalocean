package digitalocean

func expandEmail(config []interface{}) []string {
	if len(config) == 0 {
		return nil
	}
	emailList := make([]string, len(config))

	for i, v := range config {
		emailList[i] = v.(string)
	}

	return emailList
}

func flattenEmail(emails []string) []string {
	if len(emails) == 0 {
		return nil
	}

	flattenedEmails := make([]string, 0)
	for _, v := range emails {
		if v != "" {
			flattenedEmails = append(flattenedEmails, v)
		}
	}

	return flattenedEmails
}
