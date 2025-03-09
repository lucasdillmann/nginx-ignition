export function buildLoginUrl() {
    const { pathname, search } = location
    const returnTo = !search ? pathname : `${pathname}?${search}`
    return `/login?returnTo=${encodeURI(returnTo)}`
}
