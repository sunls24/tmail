import React, { useState } from "react"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx"
import { $history, updateAddress } from "@/lib/store/store.ts"
import { toast } from "@/components/ui/toast"
import { CircleArrowRight, Frown, Mail, Trash } from "lucide-react"
import { Button } from "@/components/ui/button.tsx"
import { clsx } from "clsx"
import { type language, useTranslations } from "@/i18n/ui"
import { fmtString } from "@/lib/utils.ts"

function History({
  children,
  lang,
}: {
  children: React.ReactElement
  lang: string
}) {
  const [history, setHistory] = useState<string[]>([])

  const t = useTranslations(lang as language)

  function onOpenChange(open: boolean) {
    if (open) {
      setHistory($history.get())
    }
  }

  function onConfirm() {
    $history.set([])
    toast.add({ title: t("clearHistoryTip"), type: "success" })
  }

  function onSwitch(value: string) {
    updateAddress(value)
    toast.add({ title: t("changeNew") + " " + value, type: "success" })
  }

  function onDelete(i: number) {
    const next = history.filter((_, index) => index !== i)
    $history.set(next)
    setHistory(next)
  }

  return (
    <Dialog onOpenChange={onOpenChange}>
      <DialogTrigger render={children} />
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{t("history")}</DialogTitle>
          <DialogDescription>
            {fmtString(t("historyTotal"), history.length)}
          </DialogDescription>
        </DialogHeader>
        <div className="flex max-h-96 flex-col gap-2 overflow-y-auto">
          {history.length == 0 && (
            <div className="bg-secondary text-muted-foreground flex items-center gap-1 rounded-sm px-3 py-2">
              <Frown size={20} />
              {t("nothing")}
            </div>
          )}
          {history.map((v, i) => (
            <div
              className={clsx(
                "flex items-center gap-1",
                history.length >= 8 && "mr-0.5"
              )}
              key={v}
            >
              <DialogClose
                render={
                  <button
                    type="button"
                    className="group bg-sidebar text-muted-foreground hover:text-foreground hover:bg-secondary flex flex-1 items-center rounded-sm border px-3 py-2 text-left transition-colors hover:cursor-pointer"
                  />
                }
                onClick={() => onSwitch(v)}
              >
                <Mail size={16} className="mr-2" />
                {v}
                <div className="flex-1" />
                <div className="hidden items-center gap-2 group-hover:flex">
                  <span className="text-sm">{t("switchHistory")}</span>
                  <CircleArrowRight strokeWidth={1.8} size={18} />
                </div>
              </DialogClose>
              <Button variant="ghost" size="icon" onClick={() => onDelete(i)}>
                <Trash />
              </Button>
            </div>
          ))}
        </div>
        <DialogFooter>
          <DialogClose render={<Button variant="outline" />}>
            {t("cancel")}
          </DialogClose>
          {history.length > 0 && (
            <DialogClose
              render={<Button variant="destructive" />}
              onClick={onConfirm}
            >
              {t("clearHistory")}
            </DialogClose>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export default History
