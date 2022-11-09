package main

import (
	"algogram/app"
	"algogram/errores"
	"fmt"
	"os"
)

func main() {
	var newError error
	var args = os.Args[1:]
	if len(args) != 1 {
		newError = new(errores.ErrorParametros)
		fmt.Fprintln(os.Stdout, newError.Error())
		os.Exit(0)
	}

	listaUsuarios := guardarUsuarios(abrirArchivo(args[0]))
	listaPosts := app.CrearListaDePosts()

	var usuarioLoggeado app.Usuario

	App(usuarioLoggeado, listaUsuarios, listaPosts)
}
