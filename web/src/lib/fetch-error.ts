import { toast } from "@/components/ui/toast"
import { ABORT_SAFE } from "@/lib/constant.ts"

export function fetchError(e: any) {
  if (e === ABORT_SAFE || e.name === "AbortError") {
    return
  }
  toast.add({ title: e.message ?? String(e), type: "error" })
}
