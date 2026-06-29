package localization

import (
	"strings"
)

// Translate checks the requested language and returns the translated version of the error message if found.
func Translate(lang string, message string) string {
	lang = strings.ToLower(lang)
	
	// Determine language
	if strings.HasPrefix(lang, "es") {
		return translateES(message)
	} else if strings.HasPrefix(lang, "fr") {
		return translateFR(message)
	}
	
	// Default to English (original)
	return message
}

func translateES(msg string) string {
	switch strings.ToLower(strings.TrimSpace(msg)) {
	case "email already exists":
		return "El correo electrónico ya existe"
	case "failed to hash password":
		return "Error al procesar la contraseña"
	case "user not found":
		return "Usuario no encontrado"
	case "role not found":
		return "Rol no encontrado"
	case "menu not found":
		return "Menú no encontrado"
	case "company not found":
		return "Empresa no encontrada"
	case "category not found":
		return "Categoría no encontrada"
	case "article not found":
		return "Artículo no encontrado"
	case "invalid credentials":
		return "Credenciales inválidas"
	case "user account is inactive":
		return "La cuenta de usuario está inactiva"
	case "failed to generate token":
		return "Error al generar el token de acceso"
	case "x-company-id header is required":
		return "Se requiere la cabecera X-Company-ID"
	case "invalid x-company-id format":
		return "Formato inválido de la cabecera X-Company-ID"
	case "invalid user id":
		return "ID de usuario inválido"
	case "invalid role id":
		return "ID de rol inválido"
	case "invalid menu id":
		return "ID de menú inválido"
	case "invalid company id":
		return "ID de empresa inválido"
	case "invalid category id":
		return "ID de categoría inválido"
	case "invalid article id":
		return "ID de artículo inválido"
	case "access denied: insufficient permissions":
		return "Acceso denegado: permisos insuficientes"
	case "access denied: role not found in context":
		return "Acceso denegado: rol no encontrado en el contexto"
	case "access denied: invalid role format":
		return "Acceso denegado: formato de rol inválido"
	case "authorization header is required":
		return "Se requiere la cabecera de autorización"
	case "invalid authorization header format, expected: bearer <token>":
		return "Formato de autorización inválido, se esperaba: Bearer <token>"
	case "invalid or expired token":
		return "Token inválido o expirado"
	case "invalid token claims":
		return "Atributos del token inválidos"
	case "role information not found in token":
		return "Información de rol no encontrada en el token"
	case "user id not found in token":
		return "ID de usuario no encontrado en el token"
	case "invalid user id in session":
		return "ID de usuario inválido en la sesión"
	case "unsupported user id format":
		return "Formato de ID de usuario no soportado"
	case "failed to verify company permissions":
		return "Error al verificar los permisos de la empresa"
	case "access denied to this company":
		return "Acceso denegado a esta empresa"
	default:
		return msg
	}
}

func translateFR(msg string) string {
	switch strings.ToLower(strings.TrimSpace(msg)) {
	case "email already exists":
		return "L'adresse email existe déjà"
	case "failed to hash password":
		return "Échec du traitement du mot de passe"
	case "user not found":
		return "Utilisateur non trouvé"
	case "role not found":
		return "Rôle non trouvé"
	case "menu not found":
		return "Menu non trouvé"
	case "company not found":
		return "Entreprise non trouvée"
	case "category not found":
		return "Catégorie non trouvée"
	case "article not found":
		return "Article non trouvé"
	case "invalid credentials":
		return "Identifiants invalides"
	case "user account is inactive":
		return "Le compte de l'utilisateur est inactif"
	case "failed to generate token":
		return "Échec de la génération du jeton d'accès"
	case "x-company-id header is required":
		return "L'en-tête X-Company-ID est requis"
	case "invalid x-company-id format":
		return "Format X-Company-ID invalide"
	case "invalid user id":
		return "ID d'utilisateur invalide"
	case "invalid role id":
		return "ID de rôle invalide"
	case "invalid menu id":
		return "ID de menu invalide"
	case "invalid company id":
		return "ID d'entreprise invalide"
	case "invalid category id":
		return "ID de catégorie invalide"
	case "invalid article id":
		return "ID d'article invalide"
	case "access denied: insufficient permissions":
		return "Accès refusé: permissions insuffisantes"
	case "access denied: role not found in context":
		return "Accès refusé: rôle non trouvé dans le contexte"
	case "access denied: invalid role format":
		return "Accès refusé: format de rôle invalide"
	case "authorization header is required":
		return "L'en-tête d'autorisation est requis"
	case "invalid authorization header format, expected: bearer <token>":
		return "Format d'en-tête d'autorisation invalide, attendu: Bearer <token>"
	case "invalid or expired token":
		return "Jeton invalide ou expiré"
	case "invalid token claims":
		return "Revendications de jeton invalides"
	case "role information not found in token":
		return "Informations sur le rôle introuvables dans le jeton"
	case "user id not found in token":
		return "ID d'utilisateur introuvable dans le jeton"
	case "invalid user id in session":
		return "ID d'utilisateur invalide dans la session"
	case "unsupported user id format":
		return "Format d'ID d'utilisateur non pris en charge"
	case "failed to verify company permissions":
		return "Échec de la vérification des autorisations de l'entreprise"
	case "access denied to this company":
		return "Accès refusé à cette entreprise"
	default:
		return msg
	}
}
