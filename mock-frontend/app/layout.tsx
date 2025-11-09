import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Mock Frontend - Chat Testing',
  description: 'Testing interface for events-server and chat-consumer',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}

