package auth

import (
	"backend/database"
	"backend/database/models"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx, db *database.DB) error {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	email := loginData.Email
	password := loginData.Password

	userID, err := models.GetUserIDByEmail(db.Session, email)
	if err != nil {
		if err == gocql.ErrNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials."})
		}
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An error occurred."})
	}

	user, err := models.GetUserByID(db, userID)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An error occurred."})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials."})
	}

	session, err := user.CreateSession(db, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create session."})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{"message": "Login successful.", "user": user})
}

func Register(c *fiber.Ctx, db *database.DB) error {
	var user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	name := user.Name
	email := user.Email
	password := user.Password

	var existingUser models.User
	if err := db.Session.Query(`SELECT id FROM users WHERE name = ? ALLOW FILTERING`, name).Scan(&existingUser.ID); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	_, err := models.GetUserIDByEmail(db.Session, email)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already in use."})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password."})
	}

	newUser := models.User{
		ID:       gocql.TimeUUID(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := newUser.Save(db); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user."})
	}

	userByEmail := models.UserByEmail{
		Email:  email,
		UserID: newUser.ID,
	}
	if err := db.Session.Query(`INSERT INTO users_by_email (email, user_id) VALUES (?, ?)`, userByEmail.Email, userByEmail.UserID).Exec(); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user email mapping."})
	}

	session, err := newUser.CreateSession(db, c.IP(), c.Get("User-Agent"))
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create session."})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{"message": "Registration successful", "user": newUser})
}
