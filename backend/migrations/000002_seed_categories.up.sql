INSERT INTO categories (name, slug, icon, description) VALUES
('Gastronomía', 'gastronomia', 'utensils', 'Descubre los mejores sabores y platillos locales'),
('Cultura', 'cultura', 'landmark', 'Sumérgete en la inmensa riqueza cultural de la región'),
('Historia', 'historia', 'monument', 'Explora sitios históricos y aprende sobre el pasado'),
('Aventura', 'aventura', 'mountain', 'Explora paisajes increíbles con actividades al aire libre'),
('Vida Nocturna', 'vida-nocturna', 'glass-water', 'Vive la ciudad de noche, bares, clubes y más'),
('Naturaleza', 'naturaleza', 'tree-pine', 'Conecta con la flora y fauna de la región'),
('Fotografía', 'fotografia', 'camera', 'Tours enfocados en capturar las mejores vistas')
ON CONFLICT (slug) DO NOTHING;
