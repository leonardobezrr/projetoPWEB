package handler

import (
	models "easytrady-backend/api/Models"
	repository "easytrady-backend/api/Repository"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	var usuario models.Usuarios
	if err := c.Bind(&usuario); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao decodificar dados do usuário"})
	}

	usuarios, err := repository.GetUsuarios()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao obter usuários do banco de dados"})
	}

	autenticado := false
	for _, u := range usuarios {
		if u.Email == usuario.Email && u.Senha == usuario.Senha {
			autenticado = true
			break
		}
	}

	if autenticado {
		return c.JSON(http.StatusOK, map[string]string{"message": "Login bem-sucedido para o usuário: " + usuario.Email})
	}

	return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Credenciais inválidas"})
}

func GetAllUsuarios(c echo.Context) error {
	usuarios, err := repository.GetUsuarios()
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao obter usuários do banco de dados"})
	}
	return c.JSON(http.StatusOK, usuarios)
}

func PostUsuario(c echo.Context) error {
	usuario := models.Usuarios{}
	err := c.Bind(&usuario)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid request payload")
	}

	id, err := repository.InsertUsuario(usuario)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Erro ao inserir usuário no banco de dados")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": fmt.Sprintf("Usuário inserido com sucesso ID: %d", id),
	})
}

func UpdateUsuario(c echo.Context) error {
	var usuario models.Usuarios
	err := c.Bind(&usuario)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid request payload")
	}

	id := c.Param("id")

	usuario.ID = id

	err = repository.UpdateUsuario(usuario)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Erro ao atualizar usuário no banco de dados")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("Usuário atualizado com sucesso ID: %s", usuario.ID),
	})
}