package servidor

type Repository interface {
	BuscarServidor(matricula int64) (*Servidor, error)
	ListarServidores() (*Servidores, error)
}
