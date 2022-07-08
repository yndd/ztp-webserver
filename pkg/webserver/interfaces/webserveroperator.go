package interfaces

type WebserverOperations interface {
	Run(port int, storageFolder string)
	SetKubeConfig(string)
}
