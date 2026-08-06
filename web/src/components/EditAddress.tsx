import React, { useState } from "react"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

import { Input } from "@/components/ui/input.tsx"
import { useStore } from "@nanostores/react"
import { $address, $domainList, updateAddress } from "@/lib/store/store.ts"
import { toast } from "sonner"
import { MessageCircleWarning } from "lucide-react"
import { type language, useTranslations } from "@/i18n/ui"

function EditAddress({
  children,
  lang,
}: {
  children: React.ReactNode
  lang: string
}) {
  const [address, setAddress] = useState("")
  const domainList = useStore($domainList)

  const t = useTranslations(lang as language)

  function onDomainChange(value: string) {
    setAddress(`${address!.split("@")[0]}@${value}`)
  }

  function onInputChange(value: string) {
    value = value.replace(/[^a-zA-Z0-9-_.]/g, "")
    if (value.length > 12) {
      value = value.slice(0, 12)
    }
    setAddress(`${value}@${address!.split("@")[1]}`)
  }

  function onOpenChange(open: boolean) {
    if (open) {
      setAddress($address.get())
    }
  }

  function onConfirm() {
    if (address === $address.get()) {
      return
    }
    updateAddress(address)
    toast.success(t("changeNew") + " " + address)
  }

  return (
    <AlertDialog onOpenChange={onOpenChange}>
      <AlertDialogTrigger asChild>{children}</AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{t("edit")}</AlertDialogTitle>
          <AlertDialogDescription className="flex items-center justify-center gap-1 sm:justify-start">
            <MessageCircleWarning size={20} strokeWidth={1.8} />
            {t("editWarn")}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <div className="flex items-center justify-center sm:justify-start">
          <Input
            className="w-32 text-right"
            value={address?.split("@")[0]}
            onChange={(e) => onInputChange(e.currentTarget.value)}
          />
          <span className="bg-secondary mx-1 rounded-sm p-1">@</span>
          <Select value={address?.split("@")[1]} onValueChange={onDomainChange}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {domainList.map((v) => (
                <SelectItem key={v} value={v}>
                  {v}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <AlertDialogFooter>
          <AlertDialogCancel>{t("cancel")}</AlertDialogCancel>
          <AlertDialogAction onClick={onConfirm}>
            {t("confirm")}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

export default EditAddress
