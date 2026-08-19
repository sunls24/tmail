import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx"
import { Button } from "@/components/ui/button.tsx"
import { type language, useTranslations } from "@/i18n/ui"
import { ABORT_SAFE } from "@/lib/constant.ts"
import { fetchError } from "@/lib/fetch-error.ts"
import type { Attachment, Envelope } from "@/lib/types.ts"
import { apiFetch, fmtDate, fmtString } from "@/lib/utils.ts"
import { Download, Minimize2, Paperclip, RotateCw } from "lucide-react"
import React, { useEffect, useRef, useState } from "react"

function Detail({
  children,
  envelope,
  lang,
}: {
  children: React.ReactElement
  envelope: Envelope
  lang: string
}) {
  const [loading, setLoading] = useState(false)
  const [open, setOpen] = useState(false)
  const controller = useRef<AbortController>(null)
  const [attachments, setAttachments] = useState<Attachment[]>([])
  const [content, setContent] = useState("")
  const [iframeHeight, setIframeHeight] = useState<number | null>(null)
  const iframeCleanup = useRef<(() => void) | null>(null)

  const t = useTranslations(lang as language)

  useEffect(() => {
    return () => iframeCleanup.current?.()
  }, [])

  function onOpenChange(open: boolean) {
    setOpen(open)
    if (!open) {
      iframeCleanup.current?.()
      iframeCleanup.current = null
      controller.current?.abort(ABORT_SAFE)
      controller.current = null
      setAttachments([])
      setContent("")
      setLoading(false)
      setIframeHeight(null)
      return
    }

    setIframeHeight(null)
    setLoading(true)
    controller.current?.abort(ABORT_SAFE)
    const currentController = new AbortController()
    controller.current = currentController
    apiFetch<MailDetail>(
      `/api/fetch/${envelope.id}?to=${encodeURIComponent(envelope.to)}`,
      {
        signal: currentController.signal,
      }
    )
      .then((res) => {
        if (controller.current !== currentController) {
          return
        }
        setAttachments(res.attachments)
        setContent(res.content)
      })
      .catch((error) => {
        if (controller.current === currentController) {
          fetchError(error)
        }
      })
      .finally(() => {
        if (controller.current === currentController) {
          setLoading(false)
        }
      })
  }

  function onIframeLoad(event: React.SyntheticEvent<HTMLIFrameElement>) {
    iframeCleanup.current?.()
    iframeCleanup.current = null

    const frameDocument = event.currentTarget.contentDocument
    if (!frameDocument) {
      return
    }

    for (const anchor of frameDocument.querySelectorAll("a[href]")) {
      const href = anchor.getAttribute("href")
      if (!href) {
        continue
      }

      try {
        const url = new URL(href, frameDocument.baseURI)
        if (!["http:", "https:", "mailto:", "tel:"].includes(url.protocol)) {
          anchor.removeAttribute("href")
          continue
        }
      } catch {
        anchor.removeAttribute("href")
        continue
      }

      anchor.setAttribute("target", "_blank")
      anchor.setAttribute("rel", "noopener noreferrer")
    }

    function onFrameKeyDown(keyboardEvent: KeyboardEvent) {
      if (keyboardEvent.key === "Escape") {
        keyboardEvent.preventDefault()
        onOpenChange(false)
      }
    }

    frameDocument.addEventListener("keydown", onFrameKeyDown)

    const updateHeight = () => {
      const height = Math.max(
        frameDocument.documentElement.scrollHeight,
        frameDocument.body?.scrollHeight ?? 0
      )
      setIframeHeight((currentHeight) =>
        currentHeight === height ? currentHeight : height
      )
    }

    const observer = new ResizeObserver(updateHeight)
    observer.observe(frameDocument.documentElement)
    if (frameDocument.body) {
      observer.observe(frameDocument.body)
    }
    updateHeight()

    iframeCleanup.current = () => {
      observer.disconnect()
      frameDocument.removeEventListener("keydown", onFrameKeyDown)
    }
  }

  function onDownload(id: string) {
    window.open(
      `/api/download/${encodeURIComponent(id)}?to=${encodeURIComponent(envelope.to)}`,
      "_blank",
      "noopener,noreferrer"
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger render={children} />
      <DialogContent
        showCloseButton={false}
        className="flex max-h-[calc(100dvh-2rem)] max-w-[calc(100%-2rem)] flex-col sm:max-h-[calc(100dvh-2rem)] sm:max-w-3xl"
      >
        <DialogHeader className="relative">
          <DialogTitle>{envelope.subject}</DialogTitle>
          <DialogDescription className="flex flex-col justify-between sm:flex-row">
            <span>{envelope.from}</span>
            <span>{fmtDate(envelope.created_at)}</span>
          </DialogDescription>
          <DialogClose
            render={<Button variant="ghost" size="icon" />}
            className="absolute -top-1 -right-1"
          >
            <Minimize2 />
          </DialogClose>
        </DialogHeader>
        {attachments.length > 0 && (
          <div className="flex flex-wrap gap-1">
            {attachments.map((a) => (
              <button
                type="button"
                aria-label={fmtString(t("download"), a.filename)}
                className="bg-secondary text-muted-foreground hover:text-foreground group flex items-center gap-1 rounded-sm border px-1.5 py-1 text-sm hover:cursor-pointer hover:shadow-xs"
                key={a.id}
                onClick={() => onDownload(a.id)}
              >
                <Download
                  className="animate-in fade-in hidden duration-500 group-hover:block"
                  size={16}
                />
                <Paperclip className="group-hover:hidden" size={16} />
                {a.filename}
              </button>
            ))}
          </div>
        )}
        <div className="flex min-h-0 flex-1 overflow-hidden border-t pt-4">
          {loading && (
            <div className="text-muted-foreground flex h-16 w-full shrink-0 items-center justify-center gap-1">
              <RotateCw className="animate-spin" size={18} />
              <span>{t("mailLoading")}</span>
            </div>
          )}
          {content && (
            <iframe
              title={envelope.subject}
              sandbox="allow-popups allow-popups-to-escape-sandbox allow-same-origin"
              referrerPolicy="no-referrer"
              srcDoc={content}
              onLoad={onIframeLoad}
              style={iframeHeight ? { height: `${iframeHeight}px` } : undefined}
              className="max-h-[calc(100dvh-8rem)] min-h-0 w-full border-0 sm:max-h-[calc(100dvh-8rem)]"
            />
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}

export default Detail

type MailDetail = {
  content: string
  attachments: Attachment[]
}
