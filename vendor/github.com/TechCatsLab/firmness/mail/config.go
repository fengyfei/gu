/*
 * Revision History:
 *     Initial: 2018/05/26        Li Zebang
 */

package mail

// Config a basic structure of configurations
type Config struct {
	From        string      `json:"from"`
	To          string      `json:"to"`
	Host        string      `json:"host"`
	Port        string      `json:"port"`
	Credentials Credentials `json:"credentials"`
}

// Credentials used for logging into the email account of sender
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
