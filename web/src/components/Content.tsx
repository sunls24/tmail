import Actions from "@/components/Actions.tsx"
import Detail from "@/components/Detail.tsx"
import Mounted from "@/components/Mounted.tsx"
import { Skeleton } from "@/components/ui/skeleton.tsx"
import { type language, useTranslations } from "@/i18n/ui.ts"
import { ABORT_SAFE } from "@/lib/constant.ts"
import { $address, initStore } from "@/lib/store/store.ts"
import type { Envelope } from "@/lib/types.ts"
import { fetchError, fmtDate, fmtFrom, fmtString } from "@/lib/utils.ts"
import { useStore } from "@nanostores/react"
import { clsx } from "clsx"
import {
  ClipboardCopy,
  ExternalLink,
  Frown,
  Loader,
  RotateCw,
} from "lucide-react"
import { useEffect, useMemo, useRef, useState } from "react"
import { toast } from "sonner"

function Content({ lang }: { lang: string }) {
  const [latestId, setLatestId] = useState(-1)
  const [loading, setLoading] = useState(true)
  const [envelopes, setEnvelopes] = useState<Envelope[]>([])
  const controller = useRef<AbortController>(null)

  const address = useStore($address)

  const t = useMemo(() => useTranslations(lang as language), [])

  useEffect(() => {
    fetch("/api/domain")
      .then((res) => res.json())
      .then((res) => initStore(res))
      .catch(fetchError)

    return () => controller.current?.abort(ABORT_SAFE)
  }, [])

  useEffect(() => {
    if (latestId < 0) {
      return
    }

    fetchLatest().catch(fetchError)
  }, [latestId])

  useEffect(() => {
    if (!address) {
      return
    }
    controller.current?.abort(ABORT_SAFE)
    controller.current = new AbortController()

    setLoading(true)
    setEnvelopes([])
    setLatestId(-1)
    fetchAll()
      .catch(fetchError)
      .finally(() => setLoading(false))
  }, [address])

  async function fetchAll() {
    const res = await fetch("/api/fetch?to=" + address, {
      signal: controller.current!.signal,
    })
    if (!res.ok) {
      toast.error((await res.json()).message)
      return
    }
    const list = await res.json()
    setEnvelopes(list)
    setLatestId(list.length > 0 ? list[0].id : 0)
  }

  async function fetchLatest() {
    const res = await fetch(`/api/fetch/latest?to=${address}&id=${latestId}`, {
      signal: controller.current!.signal,
    })
    if (!res.ok) {
      setTimeout(() => fetchLatest().catch(fetchError), 1000)
      toast.error((await res.json()).message)
      return
    }
    if (res.status === 204) {
      setTimeout(() => fetchLatest().catch(fetchError))
      return
    }
    const e = await res.json()
    e.animate = true
    setEnvelopes([e, ...envelopes])
    setLatestId(e.id)
    toast.success(fmtString(t("receiveNew"), e.from))
  }

  function copyToClipboard() {
    navigator.clipboard
      .writeText(address)
      .then(() => toast.success(t("copy") + " " + address))
      .catch((e) => toast.error(e.message ?? e))
  }

  return (
    <div className="flex w-full flex-col pb-4">
      <div className="block sm:hidden">
        <Actions lang={lang} />
      </div>
      <div className="relative border-x">
        <div className="animate-fill absolute h-1 bg-green-400" />
        <div className="flex flex-wrap items-center">
          <div className="bg-sidebar flex h-12 items-center border-r px-4">
            <Mounted fallback={<Skeleton className="h-6 w-48" />}>
              <span className="font-mono font-semibold">{address}</span>
            </Mounted>
          </div>
          <div
            onClick={copyToClipboard}
            className="hover:bg-sidebar flex items-center self-stretch transition-colors hover:cursor-pointer hover:border-r"
          >
            <ClipboardCopy className="mx-2" size={20} strokeWidth={1.8} />
          </div>
          <div className="flex-1" />
          <div className="text-muted-foreground hidden font-medium sm:inline">
            {t("realTime")}
          </div>
          <Loader size={20} strokeWidth={1.8} className="mx-2 animate-spin" />
        </div>
      </div>
      <div className="min-h-0 divide-y overflow-y-auto rounded-b-sm border">
        {envelopes.length === 0 && (
          <div className="text-muted-foreground flex items-center justify-center gap-1 py-5.5">
            {loading ? (
              <>
                <RotateCw className="animate-spin" size={20} />
                <span>{t("listLoading")}</span>
              </>
            ) : (
              <>
                <Frown size={20} />
                <span>{t("listEmpty")}</span>
              </>
            )}
          </div>
        )}
        {envelopes.map((envelope) => (
          <Detail lang={lang} key={envelope.id} envelope={envelope}>
            <div
              className={clsx(
                "hover:bg-secondary group text-muted-foreground space-y-1 px-4 py-2 transition-colors duration-300 hover:cursor-pointer",
                envelope.animate && "animate-in slide-in-from-right"
              )}
            >
              <div className="flex items-center">
                <span className="text-foreground">{envelope.subject}</span>
                <ExternalLink
                  size={16}
                  className="invisible mx-2 hidden group-hover:visible sm:block"
                />
                <div className="flex-1" />
                {envelope.to != address && <span>{envelope.to}</span>}
              </div>
              <div className="flex justify-between text-sm">
                <div className="truncate">{fmtFrom(envelope.from)}</div>
                <div className="shrink-0">{fmtDate(envelope.created_at)}</div>
              </div>
            </div>
          </Detail>
        ))}
      </div>
      <div className="flex-1" />
    </div>
  )
}

export default Content
