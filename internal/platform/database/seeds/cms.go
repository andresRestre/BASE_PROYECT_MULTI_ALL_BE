package seeds

import (
	"log"
	"gorm.io/gorm"
	cmsDomain "multicliente-backend/internal/features/cms/domain"
)

func SeedCMS(db *gorm.DB) error {
	// 1. Seed Texts
	texts := []cmsDomain.LandingText{
		{
			ID:       1,
			Key:      "hero_title",
			Section:  "hero",
			TextES:   "Soluciones Empresariales Multicliente de Última Generación",
			TextEN:   "Next-Generation Multi-Tenant Business Solutions",
			TextFR:   "Solutions D'entreprise Multiclient de Nouvelle Génération",
		},
		{
			ID:       2,
			Key:      "hero_subtitle",
			Section:  "hero",
			TextES:   "Transforme la gestión de su empresa con nuestra plataforma desacoplada, segura, responsive y multiplataforma.",
			TextEN:   "Transform your company management with our decoupled, secure, responsive multi-platform system.",
			TextFR:   "Transformez la gestion de votre entreprise avec notre plateforme sécurisée et multi-plateforme.",
		},
		{
			ID:       3,
			Key:      "news_section_title",
			Section:  "news",
			TextES:   "Últimas Noticias y Novedades",
			TextEN:   "Latest News & Updates",
			TextFR:   "Dernières Nouvelles et Mises à Jour",
		},
		{
			ID:       4,
			Key:      "banner_section_title",
			Section:  "banners",
			TextES:   "Galería de Banners e Innovación",
			TextEN:   "Innovation & Banner Gallery",
			TextFR:   "Galerie de Bannières et d'Innovation",
		},
	}

	for _, t := range texts {
		var existing cmsDomain.LandingText
		if err := db.Where("key = ?", t.Key).First(&existing).Error; err != nil {
			db.Create(&t)
		}
	}

	// 2. Seed News
	newsItems := []cmsDomain.LandingNews{
		{
			ID:          1,
			TitleES:     "Lanzamiento de la Plataforma Multicliente v2.0",
			TitleEN:     "Official Launch of Multi-Tenant Platform v2.0",
			TitleFR:     "Lancement Officiel de la Plateforme Multiclient v2.0",
			ContentES:   "Presentamos la nueva versión con control de acceso por roles dinámicos, soporte multilingüe y gestión avanzada de inventarios.",
			ContentEN:   "Introducing the new version with dynamic role access control, multi-language support, and advanced inventory management.",
			ContentFR:   "Présentation de la nouvelle version avec contrôle d'accès par rôles dynamiques et gestion avancée des inventaires.",
			ImageURL:    "https://images.unsplash.com/photo-1519389950473-47ba0277781c?w=800",
			IsPublished: true,
		},
		{
			ID:          2,
			TitleES:     "Optimización de Rendimiento y Arquitectura SSR",
			TitleEN:     "Performance Optimization & SSR Architecture",
			TitleFR:     "Optimisation des Performances et Architecture SSR",
			ContentES:   "La nueva landing page en React SSR garantiza tiempos de carga ultrarrápidos e indexación perfecta en motores de búsqueda.",
			ContentEN:   "The new React SSR landing page guarantees ultra-fast load times and seamless search engine indexing.",
			ContentFR:   "La nouvelle page d'atterrissage React SSR garantit des temps de chargement ultra-rapides et un référencement parfait.",
			ImageURL:    "https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800",
			IsPublished: true,
		},
	}

	for _, n := range newsItems {
		var existing cmsDomain.LandingNews
		if err := db.Where("id = ?", n.ID).First(&existing).Error; err != nil {
			db.Create(&n)
		}
	}

	// 3. Seed Banners
	banners := []cmsDomain.LandingBanner{
		{
			ID:        1,
			Title:     "Infraestructura Cloud Escalable",
			Subtitle:  "Seguridad y rendimiento garantizados 24/7",
			ImageURL:  "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=1200",
			LinkURL:   "https://google.com",
			SortOrder: 1,
			IsActive:  true,
		},
		{
			ID:        2,
			Title:     "Gestión Multiplataforma en Tiempo Real",
			Subtitle:  "Acceda desde Android, Windows Desktop o Web",
			ImageURL:  "https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=1200",
			LinkURL:   "",
			SortOrder: 2,
			IsActive:  true,
		},
	}

	for _, b := range banners {
		var existing cmsDomain.LandingBanner
		if err := db.Where("id = ?", b.ID).First(&existing).Error; err != nil {
			db.Create(&b)
		}
	}

	log.Println("✅ CMS Data Seeded")
	return nil
}
