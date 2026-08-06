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
} from "@/components/ui/dialog"
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

import { Input } from "@/components/ui/input.tsx"
import { Button } from "@/components/ui/button.tsx"
import { useStore } from "@nanostores/react"
import { $address, $domainList, updateAddress } from "@/lib/store/store.ts"
import { toast } from "@/components/ui/toast"
import { MessageCircleWarning } from "lucide-react"
import { type language, useTranslations } from "@/i18n/ui"

function EditAddress({
  children,
  lang,
}: {
  children: React.ReactElement
  lang: string
}) {
  const [address, setAddress] = useState("")
  const domainList = useStore($domainList)
  const domains = domainList.map((value) => ({ label: value, value }))

  const t = useTranslations(lang as language)

  function onDomainChange(value: string | null) {
    if (!value) {
      return
    }
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
    toast.add({ title: t("changeNew") + " " + address, type: "success" })
  }

  return (
    <Dialog onOpenChange={onOpenChange}>
      <DialogTrigger render={children} />
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{t("edit")}</DialogTitle>
          <DialogDescription className="flex items-center justify-center gap-1 sm:justify-start">
            <MessageCircleWarning size={20} strokeWidth={1.8} />
            {t("editWarn")}
          </DialogDescription>
        </DialogHeader>
        <div className="flex items-center justify-center sm:justify-start">
          <Input
            className="w-32 text-right"
            value={address?.split("@")[0]}
            onChange={(e) => onInputChange(e.currentTarget.value)}
          />
          <span className="bg-secondary mx-1 rounded-sm p-1">@</span>
          <Select
            items={domains}
            value={address?.split("@")[1]}
            onValueChange={onDomainChange}
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                {domains.map(({ label, value }) => (
                  <SelectItem key={value} value={value}>
                    {label}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
        <DialogFooter>
          <DialogClose render={<Button variant="outline" />}>
            {t("cancel")}
          </DialogClose>
          <DialogClose render={<Button />} onClick={onConfirm}>
            {t("confirm")}
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export default EditAddress
