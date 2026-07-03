package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"wander/backend/internal/config"
	"wander/backend/internal/models"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	slog.Info("starting db seeding...")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "err", err)
		os.Exit(1)
	}

	pool, err := config.ConnectDB(cfg)
	if err != nil {
		slog.Error("connect db", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	ctx := context.Background()

	// 1. Seed Categories
	slog.Info("seeding categories...")
	categories := []models.Category{
		{Name: "Senderismo", Slug: "senderismo", Icon: "mountain", Description: "Rutas a pie por senderos de montaña y naturaleza.", SortOrder: 1},
		{Name: "Aventura", Slug: "aventura", Icon: "zap", Description: "Deportes de aventura, escalada, rafting y experiencias extremas.", SortOrder: 2},
		{Name: "Gastronomía", Slug: "gastronomia", Icon: "utensils", Description: "Rutas culinarias, catas de vino y degustaciones locales.", SortOrder: 3},
		{Name: "Cultural", Slug: "cultural", Icon: "book-open", Description: "Tours históricos, visitas a museos y monumentos.", SortOrder: 4},
		{Name: "Fotografía", Slug: "fotografia", Icon: "camera", Description: "Talleres fotográficos al aire libre en paisajes únicos.", SortOrder: 5},
	}

	for _, c := range categories {
		_, err := pool.Exec(ctx, `
			INSERT INTO categories (name, slug, icon, description, sort_order)
			VALUES ($1, $2, $3, $4, $5) ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description
		`, c.Name, c.Slug, c.Icon, c.Description, c.SortOrder)
		if err != nil {
			slog.Error("insert category", "slug", c.Slug, "err", err)
		}
	}

	// 2. Seed Users (Travelers and Guides)
	slog.Info("seeding users...")
	passHash, _ := bcrypt.GenerateFromPassword([]byte("wander123"), bcrypt.DefaultCost)

	users := []struct {
		Email     string
		Username  string
		FirstName string
		LastName  string
		Role      string
		Bio       string
		Phone     string
		Avatar    string
		Lang      []string
	}{
		{
			Email:     "carlos.guia@wander.local",
			Username:  "carlos_guia",
			FirstName: "Carlos",
			LastName:  "Mendoza",
			Role:      "guide",
			Bio:       "Guía certificado de montaña con más de 10 años explorando la Sierra Madre. Amante de la fauna y la botánica local.",
			Phone:     "+52 81 1234 5678",
			Avatar:    "https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&w=300&q=80",
			Lang:      []string{"es", "en"},
		},
		{
			Email:     "sofia.aventura@wander.local",
			Username:  "sofia_guia",
			FirstName: "Sofia",
			LastName:  "Valdez",
			Role:      "guide",
			Bio:       "Especialista en escalada en roca y deportes de aventura. Te enseñaré los secretos mejor guardados del cañón.",
			Phone:     "+52 81 8765 4321",
			Avatar:    "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&w=300&q=80",
			Lang:      []string{"es", "en", "fr"},
		},
		{
			Email:     "juan.viajero@wander.local",
			Username:  "juan_traveler",
			FirstName: "Juan",
			LastName:  "Pérez",
			Role:      "traveler",
			Bio:       "Explorador de fin de semana buscando nuevas aventuras.",
			Phone:     "+52 81 0000 1111",
			Avatar:    "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=300&q=80",
			Lang:      []string{"es"},
		},
	}

	userIDs := make(map[string]int)
	for _, u := range users {
		var id int
		err := pool.QueryRow(ctx, `
			INSERT INTO users (email, username, password_hash, first_name, last_name, role, bio, phone, avatar_url, languages)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (email) DO UPDATE SET bio = EXCLUDED.bio, phone = EXCLUDED.phone, avatar_url = EXCLUDED.avatar_url
			RETURNING id
		`, u.Email, u.Username, string(passHash), u.FirstName, u.LastName, u.Role, u.Bio, u.Phone, u.Avatar, u.Lang).Scan(&id)
		if err != nil {
			slog.Error("insert user", "email", u.Email, "err", err)
			_ = pool.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", u.Email).Scan(&id)
		}
		userIDs[u.Username] = id
	}

	// 3. Seed Tours
	slog.Info("seeding tours...")
	catIDs := make(map[string]int)
	rows, _ := pool.Query(ctx, "SELECT slug, id FROM categories")
	for rows.Next() {
		var slug string
		var id int
		_ = rows.Scan(&slug, &id)
		catIDs[slug] = id
	}
	rows.Close()

	tours := []struct {
		GuideName    string
		CatSlug      string
		Title        string
		Desc         string
		Loc          string
		Price        float64
		Duration     int
		Diff         string
		MeetingPoint string
		WhatIncluded []string
		Images       []string
	}{
		{
			GuideName:    "carlos_guia",
			CatSlug:      "senderismo",
			Title:        "Ruta del Cañón de Chipinque",
			Desc:         "Una caminata de senderismo moderado por los senderos más hermosos del parque Chipinque. Avistamiento de coatíes y miradores panorámicos espectaculares.",
			Loc:          "Monterrey, NL",
			Price:        450.00,
			Duration:     180,
			Diff:         "moderate",
			MeetingPoint: "Entrada del Parque Chipinque, Caseta de cobro.",
			WhatIncluded: []string{"Entrada al parque", "Snacks nutritivos", "Botella de agua", "Bastones de senderismo"},
			Images: []string{
				"https://images.unsplash.com/photo-1551632811-561732d1e306?auto=format&fit=crop&w=800&q=80",
				"https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?auto=format&fit=crop&w=800&q=80",
			},
		},
		{
			GuideName:    "sofia_guia",
			CatSlug:      "aventura",
			Title:        "Escalada de Aventura en Potrero Chico",
			Desc:         "Desafía tus límites escalando las legendarias paredes de piedra caliza de Potrero Chico. Apto para principiantes con buena condición física.",
			Loc:          "Hidalgo, NL",
			Price:        1200.00,
			Duration:     300,
			Diff:         "challenging",
			MeetingPoint: "La Posada en Potrero Chico.",
			WhatIncluded: []string{"Equipo completo de seguridad (arnés, casco, cuerdas)", "Instrucción certificada", "Seguro de accidentes básico"},
			Images: []string{
				"https://images.unsplash.com/photo-1522163182402-834f871fd851?auto=format&fit=crop&w=800&q=80",
			},
		},
		{
			GuideName:    "carlos_guia",
			CatSlug:      "gastronomia",
			Title:        "Tour de Tacos y Cantinas Tradicionales",
			Desc:         "Explora la rica cultura culinaria del norte visitando cantinas con historia y degustando los mejores tacos de la ciudad de la mano de un experto.",
			Loc:          "Monterrey Centro",
			Price:        650.00,
			Duration:     150,
			Diff:         "easy",
			MeetingPoint: "Arco de la Independencia, Av. Pino Suárez.",
			WhatIncluded: []string{"Degustación en 4 paradas gastronómicas", "1 bebida incluida por parada", "Narrativa histórica local"},
			Images: []string{
				"https://images.unsplash.com/photo-1565299585323-38d6b0865b47?auto=format&fit=crop&w=800&q=80",
			},
		},
	}

	for _, t := range tours {
		guideID := userIDs[t.GuideName]
		catID := catIDs[t.CatSlug]
		whatIncJSON, _ := json.Marshal(t.WhatIncluded)
		imagesJSON, _ := json.Marshal(t.Images)

		var tourID int
		err := pool.QueryRow(ctx, `
			INSERT INTO tours (guide_id, category_id, title, description, location, price_per_person, duration_minutes, difficulty, what_included, meeting_point, images, is_published)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, true)
			ON CONFLICT (title, guide_id) DO UPDATE SET price_per_person = EXCLUDED.price_per_person
			RETURNING id
		`, guideID, catID, t.Title, t.Desc, t.Loc, t.Price, t.Duration, t.Diff, whatIncJSON, t.MeetingPoint, imagesJSON).Scan(&tourID)

		if err != nil {
			slog.Error("insert tour", "title", t.Title, "err", err)
			continue
		}

		slog.Info("seeding tour schedules...", "tour_id", tourID)
		times := []time.Time{
			time.Now().AddDate(0, 0, 2).Round(time.Hour),
			time.Now().AddDate(0, 0, 5).Round(time.Hour),
			time.Now().AddDate(0, 0, 9).Round(time.Hour),
		}

		for _, st := range times {
			_, err := pool.Exec(ctx, `
				INSERT INTO tour_schedules (tour_id, start_time, end_time, available_spots, is_active)
				VALUES ($1, $2, $3, $4, true)
			`, tourID, st, st.Add(time.Duration(t.Duration)*time.Minute), 8)
			if err != nil {
				slog.Error("insert schedule", "tour_id", tourID, "err", err)
			}
		}
	}

	slog.Info("Database seeded successfully!")
}
