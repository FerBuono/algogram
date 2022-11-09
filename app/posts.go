package app

type Post interface {
	// VerPost devuelve la información del post indicado
	VerPost() (string, int, string, int)

	// MostrarLikes devuelve una lista con los usuarios que likearon el post
	MostrarLikes() []string

	// Likear agrega un nuevo usuario al diccionario de likes
	Likear(usuario string)

	// Likes devuelve la cantida de likes del post
	Likes() int
}

type ListaPosts interface {
	// GuardarPost agrega un post a la lista
	GuardarPost(texto string, uid int, usuario string) Post

	// LikearPost agrega un nuevo usuario al diccionario de likes del post buscado por id
	LikearPost(id int, usuario string) error

	// MostrarLikes devuelve un arreglo con los usuarios que likearon un post
	MostrarLikes(id int) ([]string, error)
}
