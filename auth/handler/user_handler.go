package handler

import (
	"net/http"
	"time"

	//"strconv"

	"github.com/gofiber/fiber/v3"
	jtoken "github.com/golang-jwt/jwt/v5"
	middlewares "github.com/metalpoch/olt-blueprint/auth/Middlewares"
	"github.com/metalpoch/olt-blueprint/auth/constants"
	"github.com/metalpoch/olt-blueprint/auth/model"
	"github.com/metalpoch/olt-blueprint/auth/usecase"
)

type UserHandler struct {
	Usecase usecase.UserUsecase
}

func (hdlr UserHandler) Get(ctx fiber.Ctx) error {

	res, err := hdlr.Usecase.Get()
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(res)
}

func (hdlr UserHandler) Create(ctx fiber.Ctx) error {
	newUser := new(model.NewUser)
	if err := ctx.Bind().JSON(newUser); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error al insertar": err.Error()})
	}

	res, err := hdlr.Usecase.Create(newUser)
	if err != nil {
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(res)
}

func (hdlr UserHandler) GetValue(ctx fiber.Ctx) error {
	clave := ctx.Params("clave")
	valor := ctx.Params("valor")
	if clave == "clave" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "clave invalida"})
	}
	if valor == "valor" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "valor invalida"})
	}

	res, err := hdlr.Usecase.GetValue(clave, valor)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error al obtner los valores": err.Error()})
	}
	return ctx.JSON(res)
}

func (hdlr UserHandler) DeleteName(ctx fiber.Ctx) error {
	name := ctx.Params("p00")
	if name == "p00" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "fullname invalida"})
	}

	res, err := hdlr.Usecase.DeleteName(name)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error al obtner los valores": err.Error()})
	}
	return ctx.JSON(res)
}
func (hdlr UserHandler) ChangePassword(ctx fiber.Ctx) error {
	newUser := new(model.NewUser)
	if err := ctx.Bind().JSON(newUser); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.ChangePassword(newUser)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error al cambiar contraseña": err.Error()})
	}
	return ctx.JSON(res)
}
func (hdlr UserHandler) Login(ctx fiber.Ctx) error {
	//Optiene las credenciales de el body
	loginRequest := new(model.LoginRequest)

	if err := ctx.Bind().JSON(loginRequest); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//Pregunta a la db por el usuario con esas credenciales
	user, err := hdlr.Usecase.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error al inicar sesion": err.Error()})
	}
	//Creacion de los credenciales del token
	claims := jtoken.MapClaims{
		"ID":       user.Id,
		"name":     user.Fullname,
		"email":    user.Email,
		"p00":      user.P00,
		"admin":    user.IsAdmin,
		"exp":      time.Now().Add(time.Minute * 30).Unix(),
		"password": user.Password,
	}
	//Creacion del token
	temp := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	token, err := temp.SignedString([]byte(constants.SALT))

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error al generar el token": err.Error(),
		})
	}
	return ctx.JSON(model.LoginResponse{
		Token: token,
	})

}

func (hdlr UserHandler) ReadToken(ctx fiber.Ctx) error {
	loginResponse := new(model.LoginResponse)
	if err := ctx.Bind().JSON(loginResponse); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error Mal bindeo": err.Error(),
		})
	}

	token, err := middlewares.TokenCheck(loginResponse.Token)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error No autorizado": err.Error(),
		})
	}

	claims, ok := token.Claims.(*jtoken.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudieron analizar las credenciales",
		})
	}

	user := model.NewUser{
		Id:       (*claims)["ID"].(string),
		Email:    (*claims)["email"].(string),
		P00:      uint((*claims)["p00"].(float64)),
		IsAdmin:  (*claims)["admin"].(bool),
		Fullname: (*claims)["name"].(string),
		Password: (*claims)["password"].(string),
	}

	return ctx.JSON(user)
}
