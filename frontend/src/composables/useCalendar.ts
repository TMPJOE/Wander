export interface CalendarEvent {
  title: string
  description?: string
  location?: string
  start: Date | string
  end: Date | string
}

export function useCalendar() {
  const downloadIcs = (event: CalendarEvent) => {
    const startDate = new Date(event.start)
    const endDate = new Date(event.end)

    const formatDate = (date: Date) => {
      return date.toISOString().replace(/-|:|\.\d+/g, '')
    }

    const startStr = formatDate(startDate)
    const endStr = formatDate(endDate)

    const icsContent = [
      'BEGIN:VCALENDAR',
      'VERSION:2.0',
      'PRODID:-//Wander App//NONSGML Tour Booking//EN',
      'BEGIN:VEVENT',
      `SUMMARY:${event.title}`,
      `DESCRIPTION:${event.description || 'Tour booked via Wander'}`,
      `LOCATION:${event.location || ''}`,
      `DTSTART:${startStr}`,
      `DTEND:${endStr}`,
      'STATUS:CONFIRMED',
      'END:VEVENT',
      'END:VCALENDAR',
    ].join('\r\n')

    const blob = new Blob([icsContent], { type: 'text/calendar;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `${event.title.replace(/[^a-z0-9]/gi, '_').toLowerCase()}.ics`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }

  return { downloadIcs }
}
