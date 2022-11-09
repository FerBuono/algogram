package app

type Usuario interface {
	// VerUsuario devuelve la información del usuario (id y nombre)
	VerUsuario() (int, string)

	// AgregarPost agrega el último post publicado a la cola del usuario
	AgregarPost(post Post)

	// VerProximoPost devuelve el siguiente post en el feed de acuerdo a la funcion de afinidad
	VerProximoPost() (Post, error)
}

type ListaUsuarios interface {
	// GuardarUsuario agrega un nuevo usuario a la lista
	GuardarUsuario(nombre string, id int) Usuario

	// BuscarUsuario busca el usuario pedido por nombre en la lista
	BuscarUsuario(nombre string) (Usuario, error)

	// GuardarPost guarda el último post publicado en la cola de posts de cada usuario
	GuardarPost(post Post)
}
