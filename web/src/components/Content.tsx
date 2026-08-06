import Actions from "@/components/Actions.tsx"
import Detail from "@/components/Detail.tsx"
import Mounted from "@/components/Mounted.tsx"
import { Button } from "@/components/ui/button.tsx"
import { Skeleton } from "@/components/ui/skeleton.tsx"
import { type language, useTranslations } from "@/i18n/ui.ts"
import { ABORT_SAFE } from "@/lib/constant.ts"
import { $address, initStore } from "@/lib/store/store.ts"
import type { Envelope } from "@/lib/types.ts"
import {
  apiFetch,
  fetchError,
  fmtDate,
  fmtFrom,
  fmtString,
  unwrapApi,
} from "@/lib/utils.ts"
import { useStore } from "@nanostores/react"
import { clsx } from "clsx"
import {
  ClipboardCopy,
  ExternalLink,
  Frown,
  Loader,
  RotateCw,
} from "lucide-react"
import { useEffect, useRef, useState } from "react"
import { toast } from "sonner"

const PAGE_SIZE = 30

function fetchPage(address: string, signal: AbortSignal, beforeId?: number) {
  const params = new URLSearchParams({
    to: address,
    limit: String(PAGE_SIZE),
  })
  if (beforeId !== undefined) {
    params.set("before_id", String(beforeId))
  }
  return apiFetch<Envelope[]>(`/api/fetch?${params}`, { signal })
}

function Content({ lang }: { lang: string }) {
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [hasMore, setHasMore] = useState(false)
  const [envelopes, setEnvelopes] = useState<Envelope[]>([])
  const address = useStore($address)
  const controller = useRef<AbortController | null>(null)
  const t = useTranslations(lang as language)

  useEffect(() => {
    apiFetch<string[]>("/api/domain")
      .then((domainList) => initStore(domainList))
      .catch(fetchError)
  }, [])

  useEffect(() => {
    if (!address) {
      return
    }

    const currentController = new AbortController()
    controller.current = currentController
    let latestId = 0

    async function poll() {
      while (!currentController.signal.aborted) {
        try {
          const params = new URLSearchParams({
            to: address,
            id: String(latestId),
          })
          const res = await fetch(`/api/fetch/latest?${params}`, {
            signal: currentController.signal,
          })
          if (res.status === 204) {
            continue
          }

          const envelope = await unwrapApi<Envelope>(res)
          if (currentController.signal.aborted) {
            return
          }
          if (envelope.id <= latestId) {
            continue
          }

          latestId = envelope.id
          envelope.animate = true
          setEnvelopes((current) => [envelope, ...current])
          toast.success(fmtString(t("receiveNew"), envelope.from))
        } catch (error) {
          if (currentController.signal.aborted) {
            return
          }
          fetchError(error)
          await new Promise((resolve) => setTimeout(resolve, 1000))
        }
      }
    }

    async function start() {
      try {
        const list = await fetchPage(address, currentController.signal)
        if (currentController.signal.aborted) {
          return
        }

        const page = list.slice(0, PAGE_SIZE)
        latestId = page[0]?.id ?? 0
        setEnvelopes(page)
        setHasMore(list.length > PAGE_SIZE)
        void poll()
      } catch (error) {
        if (!currentController.signal.aborted) {
          fetchError(error)
        }
      } finally {
        if (!currentController.signal.aborted) {
          setLoading(false)
        }
      }
    }

    setLoading(true)
    setLoadingMore(false)
    setEnvelopes([])
    setHasMore(false)
    void start()

    return () => {
      currentController.abort(ABORT_SAFE)
      if (controller.current === currentController) {
        controller.current = null
      }
    }
  }, [address, lang])

  async function loadMore() {
    if (!address || !hasMore || loadingMore || envelopes.length === 0) {
      return
    }

    const currentController = controller.current
    if (!currentController) {
      return
    }

    const beforeId = envelopes[envelopes.length - 1].id
    setLoadingMore(true)
    try {
      const list = await fetchPage(address, currentController.signal, beforeId)
      if (currentController.signal.aborted) {
        return
      }

      const page = list.slice(0, PAGE_SIZE)
      setEnvelopes((current) => [...current, ...page])
      setHasMore(list.length > PAGE_SIZE)
    } catch (error) {
      if (!currentController.signal.aborted) {
        fetchError(error)
      }
    } finally {
      if (!currentController.signal.aborted) {
        setLoadingMore(false)
      }
    }
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
          <button
            type="button"
            aria-label={t("copyAddress")}
            onClick={copyToClipboard}
            className="hover:bg-sidebar flex items-center self-stretch border-0 bg-transparent transition-colors hover:cursor-pointer hover:border-r"
          >
            <ClipboardCopy className="mx-2" size={20} strokeWidth={1.8} />
          </button>
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
            <button
              type="button"
              className={clsx(
                "hover:bg-secondary group text-muted-foreground block w-full bg-transparent px-4 py-2 text-left transition-colors duration-300 hover:cursor-pointer",
                envelope.animate && "animate-in slide-in-from-right"
              )}
            >
              <div className="flex items-center space-y-1">
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
            </button>
          </Detail>
        ))}
        {hasMore && envelopes.length > 0 && (
          <div className="flex justify-center py-2">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={loadMore}
              disabled={loadingMore}
            >
              {loadingMore && <RotateCw className="animate-spin" />}
              {loadingMore ? t("loadingMore") : t("loadMore")}
            </Button>
          </div>
        )}
      </div>
      <div className="flex-1" />
    </div>
  )
}

export default Content
