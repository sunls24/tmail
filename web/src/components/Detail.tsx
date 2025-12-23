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
import { fetchError, fmtDate } from "@/lib/utils.ts"
import * as AlertDialogPrimitive from "@radix-ui/react-alert-dialog"
import { Download, Minimize2, Paperclip, RotateCw } from "lucide-react"
import React, { useMemo, useRef, useState } from "react"

function Detail({
  children,
  envelope,
  lang,
}: {
  children: React.ReactNode
  envelope: Envelope
  lang: string
}) {
  const divRef = useRef<HTMLDivElement>(null)
  const [loading, setLoading] = useState(false)
  const controller = useRef<AbortController>(null)
  const [attachments, setAttachments] = useState<Attachment[]>([])

  const t = useMemo(() => useTranslations(lang as language), [])

  function onOpenChange(open: boolean) {
    if (open) {
      setLoading(true)
      controller.current = new AbortController()
      fetch("/api/fetch/" + envelope.id, { signal: controller.current.signal })
        .then((res) => res.json())
        .then((res) => {
          setAttachments(res.attachments)
          divRef.current!.attachShadow({ mode: "open" }).innerHTML = res.content
        })
        .catch(fetchError)
        .finally(() => setLoading(false))
      return
    }
    setAttachments([])
    controller.current!.abort(ABORT_SAFE)
  }

  function onDownload(id: string) {
    window.open(`/api/download/${id}`, "_blank")
  }

  return (
    <AlertDialog onOpenChange={onOpenChange}>
      <AlertDialogTrigger asChild>{children}</AlertDialogTrigger>
      <AlertDialogContent className="flex max-h-11/12 flex-col sm:max-w-4xl">
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
              <div
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
              </div>
            ))}
          </div>
        )}
        <div ref={divRef} className="flex-1 overflow-auto border-t pt-4">
          {loading && (
            <div className="text-muted-foreground flex h-6.5 items-center justify-center gap-1">
              <RotateCw className="animate-spin" size={18} />
              <span className="">{t("mailLoading")}</span>
            </div>
          )}
        </div>
      </AlertDialogContent>
    </AlertDialog>
  )
}

export default Detail
