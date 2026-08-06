import React, { useState } from "react"
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog.tsx"
import { $history, updateAddress } from "@/lib/store/store.ts"
import { toast } from "sonner"
import { CircleArrowRight, Frown, Mail, Trash } from "lucide-react"
import * as AlertDialogPrimitive from "@radix-ui/react-alert-dialog"
import { Button } from "@/components/ui/button.tsx"
import { clsx } from "clsx"
import { type language, useTranslations } from "@/i18n/ui"
import { fmtString } from "@/lib/utils.ts"

function History({
  children,
  lang,
}: {
  children: React.ReactNode
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
    toast.success(t("clearHistoryTip"))
  }

  function onSwitch(value: string) {
    updateAddress(value)
    toast.success(t("changeNew") + " " + value)
  }

  function onDelete(i: number) {
    const next = history.filter((_, index) => index !== i)
    $history.set(next)
    setHistory(next)
  }

  return (
    <AlertDialog onOpenChange={onOpenChange}>
      <AlertDialogTrigger asChild>{children}</AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{t("history")}</AlertDialogTitle>
          <AlertDialogDescription>
            {fmtString(t("historyTotal"), history.length)}
          </AlertDialogDescription>
        </AlertDialogHeader>
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
              <AlertDialogPrimitive.Cancel asChild onClick={() => onSwitch(v)}>
                <button
                  type="button"
                  className="group bg-sidebar text-muted-foreground hover:text-foreground hover:bg-secondary flex flex-1 items-center rounded-sm border px-3 py-2 text-left transition-colors hover:cursor-pointer"
                >
                  <Mail size={16} className="mr-2" />
                  {v}
                  <div className="flex-1" />
                  <div className="hidden items-center gap-2 group-hover:flex">
                    <span className="text-sm">{t("switchHistory")}</span>
                    <CircleArrowRight strokeWidth={1.8} size={18} />
                  </div>
                </button>
              </AlertDialogPrimitive.Cancel>
              <Button variant="ghost" size="icon" onClick={() => onDelete(i)}>
                <Trash />
              </Button>
            </div>
          ))}
        </div>
        <AlertDialogFooter>
          <AlertDialogCancel>{t("cancel")}</AlertDialogCancel>
          {history.length > 0 && (
            <AlertDialogPrimitive.Action asChild>
              <Button variant="destructive" onClick={onConfirm}>
                {t("clearHistory")}
              </Button>
            </AlertDialogPrimitive.Action>
          )}
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

export default History
