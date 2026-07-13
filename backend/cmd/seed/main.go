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
			Avatar:    "/uploads/carlos.jpg",
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
			Avatar:    "/uploads/sofia.jpg",
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
			Title:        "Expedición al Volcán Barú",
			Desc:         "Asciende al punto más alto de Panamá. Si el clima lo permite, podrás ver ambos océanos desde la cima al amanecer.",
			Loc:          "Chiriquí, Panamá",
			Price:        1200.00,
			Duration:     720,
			Diff:         "challenging",
			MeetingPoint: "Entrada del Parque Nacional Volcán Barú, Boquete.",
			WhatIncluded: []string{"Guía certificado", "Bebida caliente", "Snacks nutritivos", "Entrada al parque"},
			Images: []string{
				"/uploads/volcanBaru.jpg",
			},
		},
		{
			GuideName:    "sofia_guia",
			CatSlug:      "aventura",
			Title:        "Aventura Guna Yala y Snorkel",
			Desc:         "Explora las paradisíacas Islas de San Blas. Disfruta de playas de arena blanca y snorkel en arrecifes de coral y barcos hundidos.",
			Loc:          "Guna Yala, Panamá",
			Price:        2500.00,
			Duration:     600,
			Diff:         "easy",
			MeetingPoint: "Puerto de Cartí.",
			WhatIncluded: []string{"Transporte en lancha", "Almuerzo tradicional Guna", "Equipo de snorkel"},
			Images: []string{
				"/uploads/islaPerro.jpg",
			},
		},
		{
			GuideName:    "carlos_guia",
			CatSlug:      "cultural",
			Title:        "Historia de Panamá Viejo al Atardecer",
			Desc:         "Camina entre las ruinas del primer asentamiento europeo en la costa pacífica de América, con vistas increíbles a la ciudad moderna.",
			Loc:          "Ciudad de Panamá",
			Price:        450.00,
			Duration:     180,
			Diff:         "easy",
			MeetingPoint: "Centro de Visitantes de Panamá Viejo.",
			WhatIncluded: []string{"Entrada al sitio arqueológico", "Guía histórico", "Botella de agua"},
			Images: []string{
				"/uploads/panamaViejo.jpg",
			},
		},
		{
			GuideName:    "sofia_guia",
			CatSlug:      "gastronomia",
			Title:        "Cena Tradicional y Show de Diablicos",
			Desc:         "Degusta la auténtica comida panameña en el Casco Antiguo mientras disfrutas de un espectáculo en vivo de Diablicos Sucios.",
			Loc:          "Casco Antiguo, Panamá",
			Price:        1100.00,
			Duration:     150,
			Diff:         "easy",
			MeetingPoint: "Restaurante Diablicos, Casco Antiguo.",
			WhatIncluded: []string{"Cena completa (plato fuerte y bebida)", "Espectáculo folclórico", "Explicación de las danzas"},
			Images: []string{
				"/uploads/restauranteDiablicos.jpg",
			},
		},
		{
			GuideName:    "carlos_guia",
			CatSlug:      "cultural",
			Title:        "Exploración de Fuertes y Piratas",
			Desc:         "Visita las ruinas de Portobelo y el Fuerte San Lorenzo. Descubre las historias de piratas y corsarios que atacaron estas costas.",
			Loc:          "Colón, Panamá",
			Price:        950.00,
			Duration:     360,
			Diff:         "moderate",
			MeetingPoint: "Plaza principal de Portobelo.",
			WhatIncluded: []string{"Transporte entre fuertes", "Guía experto en piratería", "Almuerzo ligero"},
			Images: []string{
				"/uploads/ruinasPortobelo.jpg",
			},
		},
		{
			GuideName:    "sofia_guia",
			CatSlug:      "senderismo",
			Title:        "Cascada El Gran Tife",
			Desc:         "Aventura en la selva profunda de la cordillera panameña hasta llegar a una de las cascadas más impresionantes y escondidas.",
			Loc:          "Parque Nacional Omar Torrijos",
			Price:        850.00,
			Duration:     420,
			Diff:         "challenging",
			MeetingPoint: "Entrada del Parque Nacional en El Copé.",
			WhatIncluded: []string{"Guía local", "Almuerzo estilo picnic", "Equipo de primeros auxilios"},
			Images: []string{
				"/uploads/granTife.jpg",
			},
		},
	}
	_, _ = pool.Exec(ctx, "DELETE FROM tours")

	for _, t := range tours {
		guideID := userIDs[t.GuideName]
		catID := catIDs[t.CatSlug]
		whatIncJSON, _ := json.Marshal(t.WhatIncluded)
		imagesJSON, _ := json.Marshal(t.Images)

		var tourID int
		err := pool.QueryRow(ctx, `
			INSERT INTO tours (guide_id, category_id, title, description, location, price_per_person, duration_minutes, difficulty, what_included, meeting_point, images, is_published)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, true)
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

	// 4. Seed Reviews to give guides consistent ratings
	slog.Info("seeding reviews...")
	travelerID := userIDs["juan_traveler"]
	
	rowsTours, err := pool.Query(ctx, "SELECT id, guide_id FROM tours")
	if err == nil {
		var tourInfos []struct{ id, guideID int }
		for rowsTours.Next() {
			var tInfo struct{ id, guideID int }
			_ = rowsTours.Scan(&tInfo.id, &tInfo.guideID)
			tourInfos = append(tourInfos, tInfo)
		}
		rowsTours.Close()

		for _, t := range tourInfos {
			rating := 5
			if t.guideID == userIDs["sofia_guia"] {
				rating = 4 // Make Sofia consistent
			} else {
				rating = 5 // Make Carlos consistent
			}
			
			_, err := pool.Exec(ctx, `
				INSERT INTO reviews (user_id, tour_id, rating, comment)
				VALUES ($1, $2, $3, 'Excelente experiencia, muy recomendable')
				ON CONFLICT (user_id, tour_id) DO NOTHING
			`, travelerID, t.id, rating)
			if err != nil {
				slog.Error("insert review", "tour_id", t.id, "err", err)
			}
		}
	} else {
		slog.Error("failed to query tours for reviews", "err", err)
	}

	slog.Info("Database seeded successfully!")
}
