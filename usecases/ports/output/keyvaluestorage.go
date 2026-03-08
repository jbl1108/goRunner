package output

type KeyValueStorage interface {

	/* Open a connection to the secure storage.
	Returns an error if the operation fails.
	*/
	Open() error
	/* Close the secure storage connection.
	 */
	Close() error
	/*Get data using the provided key and token.
	Returns the exported data as a byte slice or an error if the operation fails.
	*/
	Get(key string) ([]byte, error)
	/* Set data using the provided key and token.
	Returns an error if the operation fails.
	*/
	Set(key string, data []byte) error
}
