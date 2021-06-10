package constants

const ArgPasswordFile = `password-file`

const EncoderDefaultCommand = `/usr/bin/openssl aes-256-cbc -e -pbkdf2 -pass file:#{password-file}`
const DecoderDefaultCommand = `/usr/bin/openssl aes-256-cbc -d -pbkdf2 -pass file:#{password-file}`
const DefaultEditor = `/usr/bin/vim`

const GoogleDriveFolder = `.gdrive`
const DatabaseFile = `database`

const TrashFolderName = `.Trash`

const ConfigDatabaseKey = `config`

const BackUpFileExtension = `.gdrive-back`

const GDriveIgnoreFile = `.gdrive-ignore`

const AnyoneWithLink = `anyoneWithLink`
