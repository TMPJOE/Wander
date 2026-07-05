export interface CategoryPalette {
  primary: string
  primaryDark: string
  primaryLight: string
  primary50: string
  primary100: string
  rgb: string
}

const DEFAULT_PALETTE: CategoryPalette = {
  primary: '#d05342',
  primaryDark: '#b04130',
  primaryLight: '#e8755f',
  primary50: '#fdf2f0',
  primary100: '#faded8',
  rgb: '208, 83, 66',
}

const CATEGORY_PALETTES: Record<string, CategoryPalette> = {
  gastronomia: {
    primary: '#d63031', // Rojo tomate vibrante
    primaryDark: '#a82425',
    primaryLight: '#e85d5d',
    primary50: '#fde8e8',
    primary100: '#f9c4c4',
    rgb: '214, 48, 49',
  },
  cultura: {
    primary: '#b8860b', // Dorado antiguo
    primaryDark: '#8b6508',
    primaryLight: '#d4a017',
    primary50: '#fdf8eb',
    primary100: '#f5e6c8',
    rgb: '184, 134, 11',
  },
  historia: {
    primary: '#6c3483', // Púrpura histórico
    primaryDark: '#4a235a',
    primaryLight: '#8e44ad',
    primary50: '#f5eef8',
    primary100: '#e8daef',
    rgb: '108, 52, 131',
  },
  aventura: {
    primary: '#e67e22', // Naranja aventura
    primaryDark: '#b86419',
    primaryLight: '#f39c12',
    primary50: '#fef5e7',
    primary100: '#fde0b5',
    rgb: '230, 126, 34',
  },
  'vida-nocturna': {
    primary: '#8e44ad', // Púrpura nocturno
    primaryDark: '#6c3483',
    primaryLight: '#a569bd',
    primary50: '#f5eef8',
    primary100: '#e8daef',
    rgb: '142, 68, 173',
  },
  naturaleza: {
    primary: '#27ae60', // Verde naturaleza intenso
    primaryDark: '#1e8449',
    primaryLight: '#2ecc71',
    primary50: '#eafaf1',
    primary100: '#abebc6',
    rgb: '39, 174, 96',
  },
  fotografia: {
    primary: '#2980b9', // Azul profundo
    primaryDark: '#1f618d',
    primaryLight: '#3498db',
    primary50: '#ebf5fb',
    primary100: '#aed6f1',
    rgb: '41, 128, 185',
  },
}

export function getCategoryPalette(slug: string): CategoryPalette {
  return CATEGORY_PALETTES[slug] || DEFAULT_PALETTE
}

const THEME_VARS = [
  '--color-primary',
  '--color-primary-dark',
  '--color-primary-light',
  '--color-primary-50',
  '--color-primary-100',
  '--color-primary-rgb',
]

function getThemeRoot(root?: HTMLElement): HTMLElement | null {
  if (root && root.style) return root
  if (typeof document !== 'undefined' && document.documentElement?.style) {
    return document.documentElement
  }
  return null
}

export function applyCategoryTheme(slug: string, root?: HTMLElement): void {
  const target = getThemeRoot(root)
  if (!target) return

  const palette = getCategoryPalette(slug)
  target.style.setProperty('--color-primary', palette.primary)
  target.style.setProperty('--color-primary-dark', palette.primaryDark)
  target.style.setProperty('--color-primary-light', palette.primaryLight)
  target.style.setProperty('--color-primary-50', palette.primary50)
  target.style.setProperty('--color-primary-100', palette.primary100)
  target.style.setProperty('--color-primary-rgb', palette.rgb)
}

export function clearCategoryTheme(root?: HTMLElement): void {
  const target = getThemeRoot(root)
  if (!target) return
  THEME_VARS.forEach((v) => target.style.removeProperty(v))
}
