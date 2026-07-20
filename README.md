# ⚙️ Backend REST API (Go / Golang)

Servidor backend desacoplado desarrollado en Go para la plataforma **Proyecto Base Multicliente**.

---

## 🛠️ Prerrequisitos

- **Go (Golang)**: `1.22.x` o superior
- **PostgreSQL**: `14.x` o superior

---

## 🚀 Instalación e Inicio Rápido

### 1. Instalar Dependencias

```bash
cd backend
go mod download
```

### 2. Configurar Variables de Entorno (`.env`)

Crea un archivo `.env` en la raíz de `backend/` con la siguiente estructura:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_contraseña
DB_NAME=proyectobase_db
JWT_SECRET=tu_clave_secreta_super_segura
```

### 3. Ejecutar el Servidor

```bash
go run cmd/main.go
```

El servidor quedará disponible en: `http://localhost:8080/api`

---

## 📁 Estructura del Backend

- `cmd/main.go`: Punto de entrada del servidor Go.
- `internal/`: Controladores, modelos, servicios y repositorios organizados por dominio.
- `i18n/`: Archivos JSON de mensajes de respuesta en múltiples idiomas (`es.json`, `en.json`, `fr.json`).