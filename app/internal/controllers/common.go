package controllers

import (
	"agenda-app/app/internal/models"
	"agenda-app/app/internal/utils"
	"agenda-app/app/internal/utils/sessions"
	"fmt"

	"github.com/gin-gonic/gin"
)

func getUserFromSession(c *gin.Context) (user models.User, err error) {
	var tokenDeSesion string

	// Obtenemos los datos del usuario desde la sesion
	tokenDeSesion, err = c.Cookie(utils.SessionCookieHeader)
	if err != nil { // El error significa que la cookie no existe o expiró
		return user, fmt.Errorf("No se pudo obtener el token de sesión: %w", err)
	}

	data, ok := sessions.ObtenerSessionDataPorSesion(tokenDeSesion)
	if !ok {
		return user, fmt.Errorf("No se pudo obtener los datos de sesión de la cache de sesiones: %w", err)
	}

	user, ok = data.(models.User)
	if !ok {
		return user, fmt.Errorf("No se pudo obtener los datos por error interno: %w", err)
	}

	return user, nil
}

func generatePages(totalRecords, limit, currentPage int) ([]int, int) {
	if totalRecords == 0 {
		return []int{}, 0
	}

	totalPages := int(totalRecords) / limit
	if int(totalRecords)%limit != 0 || totalPages == 0 {
		totalPages++ // Si sobran registros, añade una página extra
	}

	// Generate pages
	var pages []int
	if totalPages <= 7 {
		// Generate from 1 to totalPages
		pages = make([]int, totalPages)
		for i := range totalPages {
			pages[i] = i + 1
		}
	} else {
		pages = make([]int, 7)
		startingPage := currentPage - 3
		endingPage := currentPage + 3

		if endingPage > totalPages { // Verify if we can add more pages >>
			endingPage = totalPages
			startingPage = totalPages - 6
		}

		if startingPage < 1 { // Verify if we can add more pages <<
			startingPage = 1
			endingPage = 7
		}

		// Generate from startingPage to endingPage
		for i := 0; i < 7; i++ {
			pages[i] = startingPage + i
		}
	}

	return pages, totalPages
}
