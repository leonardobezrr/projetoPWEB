package repository

import (
	db "easytrady-backend/api/DB"
	models "easytrady-backend/api/Models"
	"fmt"
	"log"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

func GetUsuarios() ([]models.Usuarios, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM usuarios;")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var usuarios []models.Usuarios
	for rows.Next() {
		var usuario models.Usuarios
		if err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.Senha); err != nil {
			log.Fatal(err)
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func InsertUsuario(usuario models.Usuarios) (id int, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	// Verificação de e-mail
	err = checkmail.ValidateFormat(usuario.Email)
	if err != nil {
		fmt.Println("Formato de e-mail inválido:")
		return 0, err
	}

	// Verificar se a senha contém pelo menos 4 caracteres
	if len(usuario.Senha) < 4 {
		errMsg := "A senha deve conter pelo menos 4 caracteres."
		fmt.Println(errMsg)
		return 0, fmt.Errorf(errMsg)
	}	

	// Criptografar senha
	hashedSenha, err := hashSenha(usuario.Senha)
	if err != nil {
		fmt.Println("Erro ao criar hash da senha:", err)
		return 0, err
	}

	sql := `INSERT INTO usuarios (nome, email, senha) VALUES ($1, $2, $3) RETURNING id`

	err = conn.QueryRow(sql, usuario.Nome, usuario.Email, hashedSenha).Scan(&id)
	if err != nil {
		errMsg := fmt.Sprintf("Erro ao inserir usuário no banco de dados: %s", err)
		fmt.Println(errMsg)
		return 0, fmt.Errorf(errMsg)
	}

	return
}

func hashSenha(senha string) (string, error) {
	hashedSenha, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashedSenha), nil
}

func UpdateUsuario(usuario models.Usuarios) error {
	conn, err := db.OpenConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	sql := `UPDATE usuarios SET nome=$1, email=$2, senha=$3 WHERE id=$4`

	_, err = conn.Exec(sql, usuario.Nome, usuario.Email, usuario.Senha, usuario.ID)
	if err != nil {
		fmt.Println("Erro ao atualizar usuário no banco de dados:", err)
		return err
	}
	return nil
}
