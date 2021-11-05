// Este pacote provê erros personalizados para a aplicação
package tools

// Erro emitido quando a entidade não foi encontrada na base
type EntityNotFoundError struct{}

func (m *EntityNotFoundError) Error() string {
	return "Entidade não encontrada"
}

// Erro emitido quando ocorre um erro ao recuper a(s) entidades na base de dados
type DataBaseError struct {
	Err error
}

func (m *DataBaseError) Error() string {
	return "Erro ao recuperar a(s) entidade na base de dados"
}
