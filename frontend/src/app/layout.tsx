import type { Metadata } from "next";

import "./globals.css";

import AuthProvider from "@/src/components/auth/authProvider";

export const metadata: Metadata = {
  title: "Todo App",
  description: "Todo application using Next.js and Go",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {

  return (
    <html lang="en">
      <body>
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}