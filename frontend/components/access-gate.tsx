"use client"

import type React from "react"

import { useState } from "react"
import Image from "next/image"
import { Button } from "@/components/ui/button"
import { Spinner } from "@/components/ui/spinner"
import { useLanguage } from "@/contexts/language-context"

interface AccessGateProps {
  children: React.ReactNode
}

export function AccessGate({ children }: AccessGateProps) {
  const [isLoading, setIsLoading] = useState(false)
  const [hasAccess, setHasAccess] = useState(true)
  const { t, language } = useLanguage()

  
  return <>{children}</>
}
