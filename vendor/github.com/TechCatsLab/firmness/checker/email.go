/*
 * Revision History:
 *     Initial: 2018/05/26        Li Zebang
 */

package checker

import (
	"regexp"
)

// IsEmail return ture if email is valid.
func IsEmail(email string) bool {
	rgx := regexp.MustCompile(`^\w+((-\w+)|(\.\w+))*\@[A-Za-z0-9]+((\.|-)[A-Za-z0-9]+)*\.[A-Za-z0-9]+$`)
	return rgx.MatchString(email)
}
