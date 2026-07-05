export function normalizeTourImages(value: unknown): string[] {
  if (Array.isArray(value)) {
    return value.filter(
      (item): item is string => typeof item === 'string' && item.trim().length > 0,
    )
  }

  if (typeof value === 'string') {
    const trimmed = value.trim()
    if (!trimmed) return []

    try {
      const parsed = JSON.parse(trimmed)
      if (Array.isArray(parsed)) {
        return parsed.filter(
          (item): item is string => typeof item === 'string' && item.trim().length > 0,
        )
      }
      if (typeof parsed === 'string' && parsed.trim()) {
        return [parsed.trim()]
      }
    } catch {
      return trimmed.startsWith('/') || trimmed.startsWith('http') || trimmed.startsWith('data:')
        ? [trimmed]
        : []
    }
  }

  return []
}
