package handlers

import (
	"fmt"
	"log"
	"net/http"

	"codeberg.org/haulproject/haul/db"
	"codeberg.org/haulproject/haul/types"
	"github.com/labstack/echo/v4"
)

// Create

func HandleV1AssemblyCreate(c echo.Context) error {
	var assembly types.Assembly
	err := c.Bind(&assembly)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Bad request",
		})
	}

	result, err := db.CreateAssembly(assembly)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	message := fmt.Sprintf("Inserted assembly with _id: %s", result.InsertedID)

	log.Println(message)

	return c.JSON(http.StatusOK, map[string]string{
		"message": message,
	})

}

func HandleV1KitCreate(c echo.Context) error {
	var kit types.Kit
	err := c.Bind(&kit)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Bad request",
		})
	}

	result, err := db.CreateKit(kit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	message := fmt.Sprintf("Inserted document with _id: %s", result.InsertedID)

	log.Println(message)

	return c.JSON(http.StatusOK, map[string]string{
		"message": message,
	})

}

// List

func HandleV1AssemblyList(c echo.Context) error {
	assemblies, err := db.ReadAll("assemblies")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, assemblies)
}

func HandleV1KitList(c echo.Context) error {
	kits, err := db.ReadAll("kits")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, kits)
}
