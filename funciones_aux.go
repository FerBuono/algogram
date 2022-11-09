package main

import (
	"algogram/app"
	"algogram/errores"
	"bufio"
	"fmt"
	"os"
)

func abrirArchivo(archivo string) *os.File {
	file, err := os.Open(archivo)
	if err != nil {
		newError := new(errores.ErrorLeerArchivo)
		fmt.Fprintln(os.Stdout, newError.Error())
		os.Exit(0)
	}
	return file
}

func guardarUsuarios(usuarios *os.File) app.ListaUsuarios {
	listaUsuarios := app.CrearListaDeUsuarios()
	id := 0
	scannerUsuarios := bufio.NewScanner(usuarios)
	for scannerUsuarios.Scan() {
		nombre := scannerUsuarios.Text()
		listaUsuarios.GuardarUsuario(nombre, id)
		id++
	}
	return listaUsuarios
}
