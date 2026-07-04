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
    primary: '#e07a3c',
    primaryDark: '#b85e23',
    primaryLight: '#ee9a6a',
    primary50: '#fdf3ee',
    primary100: '#fbe1d2',
    rgb: '224, 122, 60',
  },
  'cultura-historia': {
    primary: '#7c5a3e',
    primaryDark: '#5d4329',
    primaryLight: '#9a7a58',
    primary50: '#faf6f2',
    primary100: '#efe3d8',
    rgb: '124, 90, 62',
  },
  aventura: {
    primary: '#2c7a4b',
    primaryDark: '#1f5a37',
    primaryLight: '#479a6a',
    primary50: '#eff7f2',
    primary100: '#d6ebde',
    rgb: '44, 122, 75',
  },
  'vida-nocturna': {
    primary: '#5b3a9e',
    primaryDark: '#422873',
    primaryLight: '#7d5cc4',
    primary50: '#f4f1fb',
    primary100: '#e2d8f5',
    rgb: '91, 58, 158',
  },
  naturaleza: {
    primary: '#3e8d6f',
    primaryDark: '#2b6650',
    primaryLight: '#5fa88c',
    primary50: '#eef8f4',
    primary100: '#d4ecdf',
    rgb: '62, 141, 111',
  },
  fotografia: {
    primary: '#2f7e9e',
    primaryDark: '#225e77',
    primaryLight: '#55a0bd',
    primary50: '#f0f8fb',
    primary100: '#d6ecf3',
    rgb: '47, 126, 158',
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
