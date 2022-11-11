package main

import (
	"algogram/app"
	"algogram/errores"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func login(input []string, usuarioLoggeado *app.Usuario, listaUsuarios app.ListaUsuarios) (string, error) {
	var err error
	if *usuarioLoggeado != nil {
		return "", new(errores.UsuarioLoggeado)
	}
	usuario := strings.Join(input[1:], " ")
	*usuarioLoggeado, err = listaUsuarios.BuscarUsuario(usuario)
	if err != nil {
		return "", err
	}
	return usuario, nil
}

func logout(usuarioLoggeado *app.Usuario) error {
	if *usuarioLoggeado == nil {
		return new(errores.NoLoggeado)
	}
	*usuarioLoggeado = nil
	return nil
}

func publicar(input []string, usuarioLoggeado app.Usuario, listaPosts app.ListaPosts, listaUsuarios app.ListaUsuarios) error {
	if usuarioLoggeado == nil {
		return new(errores.NoLoggeado)
	}
	texto := strings.Join(input[1:], " ")
	uid, usuario := usuarioLoggeado.VerUsuario()
	post := listaPosts.GuardarPost(texto, uid, usuario)
	listaUsuarios.GuardarPost(post)
	return nil
}

func verSiguientePost(usuarioLoggeado app.Usuario) (app.Post, error) {
	if usuarioLoggeado == nil {
		return nil, new(errores.NoLoggeadoONoHayPosts)
	}
	post, newError := usuarioLoggeado.VerProximoPost()
	if newError != nil {
		return nil, newError
	}
	return post, nil
}

func likearPost(input []string, usuarioLoggeado app.Usuario, listaPosts app.ListaPosts) error {
	if usuarioLoggeado == nil {
		return new(errores.NoLoggeadoOPostInexistente)
	}
	id, err := strconv.Atoi(input[1])
	if err != nil {
		return new(errores.ErrorParametros)
	}
	_, usuario := usuarioLoggeado.VerUsuario()
	err = listaPosts.LikearPost(id, usuario)
	if err != nil {
		return err
	}
	return nil
}

func mostrarLikes(input []string, listaPosts app.ListaPosts) ([]string, error) {
	id, err := strconv.Atoi(input[1])
	if err != nil {
		return []string{}, new(errores.ErrorParametros)
	}
	likes, newError := listaPosts.MostrarLikes(id)
	if newError != nil || len(likes) == 0 {
		return []string{}, newError
	}
	return likes, nil
}

func App(usuarioLoggeado app.Usuario, listaUsuarios app.ListaUsuarios, listaPosts app.ListaPosts) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		input := strings.Split(scanner.Text(), " ")
		action := input[0]

		switch action {

		case "login":
			usuario, newError := login(input, &usuarioLoggeado, listaUsuarios)
			if newError != nil {
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			fmt.Println("Hola", usuario)

		case "logout":
			newError := logout(&usuarioLoggeado)
			if newError != nil {
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			fmt.Println("Adios")

		case "publicar":
			newError := publicar(input, usuarioLoggeado, listaPosts, listaUsuarios)
			if newError != nil {
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			fmt.Println("Post publicado")

		case "ver_siguiente_feed":
			post, newError := verSiguientePost(usuarioLoggeado)
			if newError != nil {
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			texto, id, usuario, _ := post.VerPost()
			fmt.Println("Post ID", id)
			fmt.Println(usuario, "dijo:", texto)
			fmt.Println("Likes:", post.Likes())

		case "likear_post":
			newError := likearPost(input, usuarioLoggeado, listaPosts)
			if newError != nil {
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			fmt.Println("Post likeado")

		case "mostrar_likes":
			likes, newError := mostrarLikes(input, listaPosts)
			if newError != nil {
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			fmt.Println("El post tiene", len(likes), "likes:")
			for _, usuario := range likes {
				fmt.Printf("\t%s\n", usuario)
			}

		default:
			fmt.Fprintln(os.Stdout, "Comando incorrecto")
		}
	}
}
