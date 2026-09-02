// internal/services/order_service.go
package services

import (
	"agenda-app/app/internal/errorhandling"
	"agenda-app/app/internal/models"
	"agenda-app/app/internal/repository"
	"agenda-app/app/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/datatypes"
)

type UserService interface {
	RegisterNewUser(requestor, user models.User) (string, error)
	ListUsers(requestor, query models.User) ([]models.User, error)
	UpdateUser(requestor, model models.User) error
	GetUserByEmail(email string) (models.User, error)
	GetUserById(id uint) (models.User, error)
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

type userService struct {
	repo repository.UserRepository
}

func CheckRole(requestorRole, userRole models.UserRole) error {
	switch requestorRole {
	case models.UserRoleUser:
		return fmt.Errorf("%w: usuario con rol 'User' no puede realizar esta accion", errorhandling.ErrForbidden)
	case models.UserRoleAdmin:
		newUserIsAdminOrUser := userRole == models.UserRoleAdmin || userRole == models.UserRoleUser
		if !newUserIsAdminOrUser {
			return fmt.Errorf("%w: administradores solo pueden ejecutar acciones sobre los roles 'User' o 'Admin", errorhandling.ErrForbidden)
		}

		return nil
	case models.UserRoleSuperAdmin:
		return nil
	default:
		return fmt.Errorf("%w: rol '%s' no reconocido", errorhandling.ErrInternal, requestorRole)
	}
}

func (s *userService) GetUserById(id uint) (models.User, error) {
	return s.repo.ReadByID(id)
}

func (s *userService) RegisterNewUser(requestor, newUser models.User) (string, error) {
	// Valida permisos, el superadmin no se valida porque puede crear cualquier rol
	switch requestor.Role {
	case models.UserRoleUser:
		return "", fmt.Errorf("%w: usuario con rol 'User' no puede crear nuevos usuarios", errorhandling.ErrForbidden)
	case models.UserRoleAdmin:
		newUserIsAdminOrUser := newUser.Role == models.UserRoleAdmin || newUser.Role == models.UserRoleUser
		if !newUserIsAdminOrUser {
			return "", fmt.Errorf("%w: administradores solo pueden crear usuarios con rol 'User' o 'Admin", errorhandling.ErrForbidden)
		}
	}

	// Valida que el correo electronico no exista
	_, err := s.repo.ReadByQuery(models.User{Email: newUser.Email})
	if err == nil {
		return "", fmt.Errorf("%w: el correo electronico '%s' ya esta registrado", errorhandling.ErrDuplicatedError, newUser.Email)
	}

	// Hashea la contraseña
	newUser.PasswordHash, err = utils.HashPassword(newUser.Password)
	if err != nil {
		return "", fmt.Errorf("%w: error al encryptar password: %v", errorhandling.ErrInternal)
	}

	newUser.RequiresPasswordUpdate = true
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	return newUser.Password, s.repo.Create(newUser)
}

func (s *userService) ListUsers(requestor, query models.User) ([]models.User, error) {
	// Valida permisos
	if query.Role != "" {
		if err := CheckRole(requestor.Role, query.Role); err != nil {
			return nil, err
		}
	}

	// Preparar query con las propiedades permitidas para la busqueda
	query = models.User{
		UserFullName: query.UserFullName,
		Email:        query.Email,
		Role:         query.Role,
		Status:       query.Status,
	}

	return s.repo.ReadByQuery(query)
}

func (s *userService) UpdateUser(requestor, model models.User) error {
	// Valida que el usuario exista
	userInDB, err := s.repo.ReadByID(model.ID)
	if errors.Is(err, errorhandling.ErrNotFoundError) {
		return fmt.Errorf("%w: el usuario con id '%d' no existe", err, model.ID)
	}
	if err != nil {
		return fmt.Errorf("%w: error intentando encontrar usuario a modificar: %v", errorhandling.ErrInternal, err)
	}

	// Valida permisos
	if err := CheckRole(requestor.Role, userInDB.Role); err != nil {
		return err
	}

	// Preparar modelo con las propiedades permitidas para modificar
	// - Prepara propiedades base
	updates := models.User{
		UserFullName: model.UserFullName,
		Email:        model.Email,
		Phone:        model.Phone,
		Role:         model.Role,
		Status:       model.Status,
		UpdatedAt:    time.Now(),
	}

	// - Preparar pws si se modifica
	if model.Password != "" {
		updates.Password, err = utils.HashPassword(model.Password)
		if err != nil {
			return fmt.Errorf("%w: error al encryptar password: %v", errorhandling.ErrInternal, err)
		}

		updates.RequiresPasswordUpdate = true
	}

	// - Preparar change history
	{
		// Convertir string a arreglo de users
		currentChangeLog := []models.User{}
		err = json.Unmarshal([]byte(userInDB.ChangeHistory), &currentChangeLog)
		if err != nil {
			return fmt.Errorf("%w: error al deserializar change log: %v", errorhandling.ErrInternal, err)
		}

		// Agregar user al arreglo
		userInDB.StoredInChangeLogAt = time.Now()
		currentChangeLog = append(currentChangeLog, userInDB)

		// Convertir de nuevo a string
		newChangeLog, err := json.Marshal(currentChangeLog)
		if err != nil {
			return fmt.Errorf("%w: error al serializar change log: %v", errorhandling.ErrInternal, err)
		}

		// Actualizar change log
		updates.ChangeHistory = datatypes.JSON(newChangeLog)
	}

	return s.repo.Update(updates)
}

func (s *userService) GetUserByEmail(email string) (models.User, error) {
	results, err := s.repo.ReadByQuery(models.User{Email: email})
	if err != nil {
		return models.User{}, fmt.Errorf("%w: error al recuperar usuario por email", err)
	}

	if len(results) == 0 {
		return models.User{}, errorhandling.ErrNotFoundError
	}

	return results[0], nil
}
