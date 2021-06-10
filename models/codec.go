package models

type Codec struct {
	Encoder string `yaml:"encoder" default-value:"openssl aes-256-cbc -e -pbkdf2 -pass file:#{password-file}"`
	Decoder string `yaml:"decoder" default-value:"openssl aes-256-cbc -d -pbkdf2 -pass file:#{password-file}"`
}
