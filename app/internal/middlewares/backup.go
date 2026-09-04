package middlewares

import (
	"agenda-app/app/config"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	PostgresConnString = "postgresql://%s:%s@%s:%s/%s" //  postgresql://usuario:password@localhost:5432/mi_base_datos
)

// Mapa de rutas exactas que queremos ignorar (Búsqueda O(1))
var excludedEndpoints = map[string]bool{
	"/login/": true,
}

var (
	dbConnStr       string
	backupDirectory string
	maxBackups      int
)

// Estructura para gestionar el estado del "1 Activo, 1 Pendiente"
type BackupScheduler struct {
	mu          sync.Mutex
	isExecuting bool // Indica si pg_dump está corriendo actualmente
	hasPending  bool // Indica si ya hay un respaldo esperando en fila
}

var scheduler = &BackupScheduler{}

// Trigger intenta programar un nuevo respaldo siguiendo las reglas de descarte
func (s *BackupScheduler) Trigger() {
	s.mu.Lock()

	// Caso 1: No hay nadie ejecutando. Empezamos de inmediato.
	if !s.isExecuting {
		s.isExecuting = true
		s.mu.Unlock()
		go s.runLoop()
		return
	}

	// Caso 2: Ya hay uno ejecutando, pero NO hay nadie pendiente. Anotamos el pendiente.
	if !s.hasPending {
		s.hasPending = true
		log.Println("[Backup Status] Servidor ocupado. Respaldo agendado como PENDIENTE.")
		s.mu.Unlock()
		return
	}

	// Caso 3: Ya hay uno ejecutando Y además ya hay uno pendiente. Ignoramos los nuevos.
	log.Println("[Backup Status] Ocupado y con pendiente en cola. Petición IGNORADA para proteger el servidor.")
	s.mu.Unlock()
}

// runLoop ejecuta el respaldo activo y revisa si quedó alguno pendiente al terminar
func (s *BackupScheduler) runLoop() {
	for {
		log.Println("[Backup Worker] Iniciando ejecución de pg_dump...")
		if err := executeLocalBackup(dbConnStr, backupDirectory); err != nil {
			log.Printf("[Backup Error] Falló el respaldo: %v", err)
		} else {
			log.Println("[Backup Exitoso] Respaldo completado.")
		}

		// Al terminar, revisamos bajo llave si hay un pendiente esperando
		s.mu.Lock()
		if s.hasPending {
			// El que estaba pendiente ahora pasa a ejecutarse
			s.hasPending = false
			s.mu.Unlock()
			// Pequeña pausa de cortesía para el disco antes de volver a empezar
			time.Sleep(1 * time.Second)
			continue // Repite el ciclo for
		}

		// Si no había nadie pendiente, apagamos el motor
		s.isExecuting = false
		s.mu.Unlock()
		break
	}
}

func BackupLocalMiddleware(cfg config.DB, backupDir string) gin.HandlerFunc {
	dbConnStr = fmt.Sprintf(PostgresConnString, cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DBName)
	backupDirectory = backupDir
	maxBackups = 10

	if err := os.MkdirAll(backupDirectory, 0755); err != nil {
		log.Fatalf("[backup] Error crítico al crear el directorio de respaldos: %v", err)
	}

	// Inicializamos el planificador de respaldos
	scheduler.mu.Lock()
	scheduler.isExecuting = false
	scheduler.hasPending = false
	scheduler.mu.Unlock()

	return func(c *gin.Context) {
		c.Next()

		method := c.Request.Method
		if method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH" {
			currentPath := c.Request.URL.Path
			if excludedEndpoints[currentPath] {
				// Si la ruta está excluida, dejamos que la petición continúe
				// pero terminamos el middleware aquí para omitir el respaldo por completo.
				c.Next()
				return
			}

			go scheduler.Trigger()
		}
	}
}

func executeLocalBackup(pgConnString, backupDirectory string) error {
	fileName := fmt.Sprintf("backup_%s.sql", time.Now().Format("2006-01-02_15-04-05"))
	outputPath := filepath.Join(backupDirectory, fileName)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("[backup] no se pudo crear el archivo: %w", err)
	}
	defer outputFile.Close()

	cmd := exec.Command("pg_dump", pgConnString)
	cmd.Stdout = outputFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		_ = os.Remove(outputPath)
		return fmt.Errorf("[backup] error en pg_dump: %w", err)
	}

	log.Println("[Backup Retention] Iniciando depuración de respaldos antiguos...")
	if err := cleanOldBackups(); err != nil {
		log.Printf("[Backup Retention Error] No se pudieron limpiar archivos antiguos: %v", err)
		// No retornamos error aquí para no marcar el respaldo actual como fallido
	}

	return nil
}

// cleanOldBackups busca en la carpeta de respaldos y elimina los excedentes
func cleanOldBackups() error {
	// Leer todos los archivos del directorio de respaldos
	files, err := os.ReadDir(backupDirectory)
	if err != nil {
		return fmt.Errorf("error al leer el directorio: %w", err)
	}

	// Filtrar solo archivos con extensión .sql (por si hay carpetas u otros elementos)
	var backupFiles []os.FileInfo
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".sql" {
			info, err := f.Info()
			if err == nil {
				backupFiles = append(backupFiles, info)
			}
		}
	}

	// Si hay 10 o menos archivos, no hay nada que limpiar
	if len(backupFiles) <= maxBackups {
		return nil
	}

	// Ordenar los archivos por su fecha de modificación: los más NUEVOS primero
	sort.Slice(backupFiles, func(i, j int) bool {
		return backupFiles[i].ModTime().After(backupFiles[j].ModTime())
	})

	// Los elementos desde el índice MaxBackupRetention (10) en adelante son los antiguos y deben borrarse
	for i := maxBackups; i < len(backupFiles); i++ {
		oldFilePath := filepath.Join(backupDirectory, backupFiles[i].Name())
		log.Printf("[Backup Retention] Eliminando respaldo obsoleto: %s", backupFiles[i].Name())

		if err := os.Remove(oldFilePath); err != nil {
			log.Printf("[Backup Retention Error] No se pudo borrar %s: %v", oldFilePath, err)
		}
	}

	return nil
}
