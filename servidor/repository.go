package servidor

type Repository interface {
	BuscarServidor(matricula int64, c *Servidor) error
	ListarServidores(c *Servidores) error
}
