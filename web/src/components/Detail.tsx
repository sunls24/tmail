import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog.tsx"
import { Button } from "@/components/ui/button.tsx"
import { type language, useTranslations } from "@/i18n/ui"
import { ABORT_SAFE } from "@/lib/constant.ts"
import type { Attachment, Envelope } from "@/lib/types.ts"
import { apiFetch, fetchError, fmtDate, fmtString } from "@/lib/utils.ts"
import * as AlertDialogPrimitive from "@radix-ui/react-alert-dialog"
import { clsx } from "clsx"
import { Download, Minimize2, Paperclip, RotateCw } from "lucide-react"
import React, { useRef, useState } from "react"

function Detail({
  children,
  envelope,
  lang,
}: {
  children: React.ReactNode
  envelope: Envelope
  lang: string
}) {
  const [loading, setLoading] = useState(false)
  const [expanded, setExpanded] = useState(false)
  const controller = useRef<AbortController>(null)
  const [attachments, setAttachments] = useState<Attachment[]>([])
  const [content, setContent] = useState("")

  const t = useTranslations(lang as language)

  function onOpenChange(open: boolean) {
    if (!open) {
      controller.current?.abort(ABORT_SAFE)
      controller.current = null
      setAttachments([])
      setContent("")
      setLoading(false)
      return
    }

    setExpanded(false)
    setLoading(true)
    controller.current?.abort(ABORT_SAFE)
    const currentController = new AbortController()
    controller.current = currentController
    apiFetch<MailDetail>("/api/fetch/" + envelope.id, {
      signal: currentController.signal,
    })
      .then((res) => {
        if (controller.current !== currentController) {
          return
        }
        setAttachments(res.attachments)
        setContent(res.content)
        setExpanded(true)
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

  function onDownload(id: string) {
    window.open(
      `/api/download/${encodeURIComponent(id)}`,
      "_blank",
      "noopener,noreferrer"
    )
  }

  return (
    <AlertDialog onOpenChange={onOpenChange}>
      <AlertDialogTrigger asChild>{children}</AlertDialogTrigger>
      <AlertDialogContent
        className={clsx(
          "flex flex-col sm:max-w-4xl",
          expanded
            ? "h-[85dvh] max-h-[85dvh] sm:h-[min(75dvh,48rem)]"
            : "max-h-11/12"
        )}
      >
        <AlertDialogHeader className="relative">
          <AlertDialogTitle>{envelope.subject}</AlertDialogTitle>
          <AlertDialogDescription className="flex flex-col justify-between sm:flex-row">
            <span>{envelope.from}</span>
            <span>{fmtDate(envelope.created_at)}</span>
          </AlertDialogDescription>
          <AlertDialogPrimitive.Cancel
            asChild
            className="absolute -top-1 -right-1"
          >
            <Button variant="ghost" size="icon">
              <Minimize2 />
            </Button>
          </AlertDialogPrimitive.Cancel>
        </AlertDialogHeader>
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
              sandbox=""
              referrerPolicy="no-referrer"
              srcDoc={content}
              className="h-full w-full border-0"
            />
          )}
        </div>
      </AlertDialogContent>
    </AlertDialog>
  )
}

export default Detail

type MailDetail = {
  content: string
  attachments: Attachment[]
}
