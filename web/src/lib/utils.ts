import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"
import { ABORT_SAFE } from "@/lib/constant.ts"
import { toast } from "sonner"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function fmtDate(str: string) {
  return new Date(str).toLocaleString()
}

export function fmtFrom(str: string) {
  const match = str.match(/^(.+?)\s*<(.+?)>$/)
  return match ? match[1].replace(/^"|"$/g, "") : str
}

type ApiEnvelope<T> = {
  code: number
  message: string
  data?: T
}

async function readBody<T>(res: Response) {
  const text = await res.text()
  if (!text) {
    return undefined
  }
  try {
    return JSON.parse(text) as T
  } catch {
    throw new Error(text)
  }
}

export async function unwrapApi<T>(res: Response) {
  if (res.status === 204) {
    return undefined as T
  }

  const body = await readBody<ApiEnvelope<T>>(res)
  if (!res.ok) {
    throw new Error(body?.message || res.statusText)
  }
  if (body?.code !== 0) {
    throw new Error(body?.message || "request failed")
  }

  return body.data as T
}

export async function apiFetch<T>(input: RequestInfo | URL, init?: RequestInit) {
  return unwrapApi<T>(await fetch(input, init))
}

function randomStr(length: number) {
  const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
  let result = ""
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}

export function randomAddress(domain: string) {
  return `${randomStr(8)}@${domain}`
}

export function fetchError(e: any) {
  if (e === ABORT_SAFE || e.name === "AbortError") {
    return
  }
  toast.error(e.message ?? e)
}

export function fmtString(template: string, ...values: any[]) {
  return template.replace(/{(\d+)}/g, (match, index) => values[index] ?? match)
}
