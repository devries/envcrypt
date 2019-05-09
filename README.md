# envcrypt
Envelope encryption pipe commands using Cloud KMS

## Introduction
This package creates two commands described below:

- `pencrypt`: An envelope encryption pipe which creates a random AES256 encryption key,
  encrypts that key using Google Cloud KMS, and then encrypts the input message using
  a GCM cipher with a random 12 byte nonce. The encrypted message is output in JSON
  format with the Cloud KMS encrypted key and the encrypted input file.
  
- `pdecrypt`: This command reverses the operation, using Cloud KMS to decrypt the AES256 key,
  then decrypting the corresponding message.
  
By default each command reads from STDIN and writes to STDOUT, but it is possible to use the "-i"
and "-o" flags to read and write from output files.
