import { persistentAtom } from "@nanostores/persistent"
import { randomAddress } from "@/lib/utils.ts"
import { atom } from "nanostores"

export const $domainList = atom<string[]>([])
export const $address = persistentAtom<string>("address")
export const $history = persistentAtom<string[]>("history", [], {
  encode: JSON.stringify,
  decode: JSON.parse,
})

export function initStore(domainList: string[]) {
  domainList && $domainList.set(domainList)
  if ($address.get()) {
    return
  }
  const cfg = localStorage.getItem("config")
  if (!cfg) {
    updateAddress(randomAddress(domainList[0]))
  } else {
    const state = JSON.parse(cfg).state
    $history.set(state.history)
    updateAddress(state.mail + state.domain)
  }
}

export function updateAddress(address: string) {
  const current = $address.get()
  if (current === address) {
    return
  }
  $address.set(address)
  if (!current) {
    return
  }
  const history = $history.get().filter((item) => item !== address)
  $history.set([current, ...history])
}
